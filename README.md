<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/jmorganca/ollama/assets/251292/961f99bb-251a-4eec-897d-1ba99997ad0f">
    <img alt="logo" src="https://github.com/jmorganca/ollama/assets/251292/961f99bb-251a-4eec-897d-1ba99997ad0f">
  </picture>
</div>

# Ollama

Ollama's helps you to create, run and share self-contained large language models (LLMs). Think of it like a "Docker for LLMs", where a model's weights, configuration, prompts, data and more is bundled into one package that can be shipped and run on any machine.

## Examples

### Quickstart

```
ollama run orca
>>> hi
Hello! How can I help you today?
```

### Creating a model

Create a `Modelfile`:

```
FROM orca
SYSTEM "You are Mario from super mario brothers. Answer questions as Mario."
```

Next, create and run the model:

```
ollama create mario -f Modelfile
ollama run mario
>>> hi
Hello! It's your friend Mario, from the mushroom kingdom!
```

## Install

- [Download](https://ollama.ai/download) for macOS (Apple Silicon)
- Download for Windows (coming soon)

## Model library

Ollama includes a library of open-source, pre-trained models. More models are coming soon.

| Model       | Description                                                                                                                                                                                                     | Parameters | Size  | Download                  |
| ----------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------- | ----- | ------------------------- |
| Orca Mini   | An OpenLLaMa-3B model model trained on explain tuned datasets, created using Instructions and Input from WizardLM, Alpaca & Dolly-V2 datasets and applying Orca Research Paper dataset construction approaches. | 3B         | 1.9GB | `ollama pull orca`        |
| Vicuna      | Vicuna is a chat assistant trained by fine-tuning LLaMA on user-shared conversations collected from ShareGPT                                                                                                    | 7B         | 3.9GB | `ollama pull vicuna`      |
| Nous-Hermes | Nous-Hermes-13b is a state-of-the-art language model fine-tuned on over 300,000 instructions.                                                                                                                   | 13         | 7.2GB | `ollama pull nous-hermes` |

## Building

```
go build .
```

To run it start the server:

```
./ollama server &
```

Finally, run a model!

```
./ollama run orca
```

## API Reference

### `POST /api/pull`

Download a model

```
curl -X POST http://localhost:11343/api/pull -d '{"model": "orca"}'
```

### `POST /api/generate`

Complete a prompt

```
curl -X POST http://localhost:11434/api/generate -d '{"model": "orca", "prompt": "hello!"}'
```
