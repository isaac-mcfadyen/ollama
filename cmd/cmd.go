package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/jmorganca/ollama/api"
	"github.com/jmorganca/ollama/server"
)

func cacheDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(home, ".ollama")
}

func create(cmd *cobra.Command, args []string) error {
	filename, _ := cmd.Flags().GetString("file")
	client := api.NewClient()

	request := api.CreateRequest{Name: args[0], Path: filename}
	fn := func(resp api.CreateProgress) error {
		fmt.Println(resp.Status)
		return nil
	}

	if err := client.Create(context.Background(), &request, fn); err != nil {
		return err
	}

	return nil
}

func RunRun(cmd *cobra.Command, args []string) error {
	_, err := os.Stat(args[0])
	switch {
	case errors.Is(err, os.ErrNotExist):
		if err := pull(args[0]); err != nil {
			var apiStatusError api.StatusError
			if !errors.As(err, &apiStatusError) {
				return err
			}

			if apiStatusError.StatusCode != http.StatusBadGateway {
				return err
			}
		}
	case err != nil:
		return err
	}

	return RunGenerate(cmd, args)
}

func push(cmd *cobra.Command, args []string) error {
	client := api.NewClient()

	request := api.PushRequest{Name: args[0]}
	fn := func(resp api.PushProgress) error {
		fmt.Println(resp.Status)
		return nil
	}

	if err := client.Push(context.Background(), &request, fn); err != nil {
		return err
	}
	return nil
}

func RunPull(cmd *cobra.Command, args []string) error {
	return pull(args[0])
}

func pull(model string) error {
	client := api.NewClient()

	var bar *progressbar.ProgressBar

	currentLayer := ""
	request := api.PullRequest{Name: model}
	fn := func(resp api.PullProgress) error {
		if resp.Digest != currentLayer && resp.Digest != "" {
			if currentLayer != "" {
				fmt.Println()
			}
			currentLayer = resp.Digest
			layerStr := resp.Digest[7:23] + "..."
			bar = progressbar.DefaultBytes(
				int64(resp.Total),
				"pulling "+layerStr,
			)
		} else if resp.Digest == currentLayer && resp.Digest != "" {
			bar.Set(resp.Completed)
		} else {
			currentLayer = ""
			fmt.Println(resp.Status)
		}
		return nil
	}

	if err := client.Pull(context.Background(), &request, fn); err != nil {
		return err
	}
	return nil
}

func RunGenerate(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		// join all args into a single prompt
		return generate(cmd, args[0], strings.Join(args[1:], " "))
	}

	if term.IsTerminal(int(os.Stdin.Fd())) {
		return generateInteractive(cmd, args[0])
	}

	return generateBatch(cmd, args[0])
}

var generateContextKey struct{}

func generate(cmd *cobra.Command, model, prompt string) error {
	if len(strings.TrimSpace(prompt)) > 0 {
		client := api.NewClient()

		spinner := progressbar.NewOptions(-1,
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionThrottle(60*time.Millisecond),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionSetElapsedTime(false),
			progressbar.OptionClearOnFinish(),
		)

		go func() {
			for range time.Tick(60 * time.Millisecond) {
				if spinner.IsFinished() {
					break
				}

				spinner.Add(1)
			}
		}()

		var latest api.GenerateResponse

		generateContext, ok := cmd.Context().Value(generateContextKey).([]int)
		if !ok {
			generateContext = []int{}
		}

		request := api.GenerateRequest{Model: model, Prompt: prompt, Context: generateContext}
		fn := func(resp api.GenerateResponse) error {
			if !spinner.IsFinished() {
				spinner.Finish()
			}

			latest = resp

			fmt.Print(resp.Response)

			cmd.SetContext(context.WithValue(cmd.Context(), generateContextKey, resp.Context))
			return nil
		}

		if err := client.Generate(context.Background(), &request, fn); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println()

		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			return err
		}

		if verbose {
			latest.Summary()
		}
	}

	return nil
}

func generateInteractive(cmd *cobra.Command, model string) error {
	fmt.Print(">>> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if err := generate(cmd, model, scanner.Text()); err != nil {
			return err
		}

		fmt.Print(">>> ")
	}

	return nil
}

func generateBatch(cmd *cobra.Command, model string) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		prompt := scanner.Text()
		fmt.Printf(">>> %s\n", prompt)
		if err := generate(cmd, model, prompt); err != nil {
			return err
		}
	}

	return nil
}

func RunServer(_ *cobra.Command, _ []string) error {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("OLLAMA_PORT")
	if port == "" {
		port = "11434"
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return err
	}

	return server.Serve(ln)
}

func NewCLI() *cobra.Command {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	rootCmd := &cobra.Command{
		Use:          "ollama",
		Short:        "Large language model runner",
		SilenceUsage: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRunE: func(_ *cobra.Command, args []string) error {
			// create the models directory and it's parent
			return os.MkdirAll(filepath.Join(cacheDir(), "models"), 0o700)
		},
	}

	cobra.EnableCommandSorting = false

	createCmd := &cobra.Command{
		Use:   "create MODEL",
		Short: "Create a model from a Modelfile",
		Args:  cobra.MinimumNArgs(1),
		RunE:  create,
	}

	createCmd.Flags().StringP("file", "f", "Modelfile", "Name of the Modelfile (default \"Modelfile\")")

	runCmd := &cobra.Command{
		Use:   "run MODEL [PROMPT]",
		Short: "Run a model",
		Args:  cobra.MinimumNArgs(1),
		RunE:  RunRun,
	}

	runCmd.Flags().Bool("verbose", false, "Show timings for response")

	serveCmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"start"},
		Short:   "Start ollama",
		RunE:    RunServer,
	}

	pullCmd := &cobra.Command{
		Use:   "pull MODEL",
		Short: "Pull a model from a registry",
		Args:  cobra.MinimumNArgs(1),
		RunE:  RunPull,
	}

	pushCmd := &cobra.Command{
		Use:   "push MODEL",
		Short: "Push a model to a registry",
		Args:  cobra.MinimumNArgs(1),
		RunE:  push,
	}

	rootCmd.AddCommand(
		serveCmd,
		createCmd,
		runCmd,
		pullCmd,
		pushCmd,
	)

	return rootCmd
}
