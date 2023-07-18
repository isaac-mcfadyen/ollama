<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/jmorganca/ollama/assets/251292/961f99bb-251a-4eec-897d-1ba99997ad0f">
    <img alt="logo" src="https://github.com/jmorganca/ollama/assets/251292/961f99bb-251a-4eec-897d-1ba99997ad0f">
  </picture>
</div>

# Ollama

Create, run and share self-contained large language models (LLMs). Think of it like a "Docker for LLMs", where a model's weights, configuration, prompts, data and more is bundled into one package that can be shipped and run on any machine.

## Examples

### Quickstart

```
ollama run llama2
>>> hi
Hello! How can I help you today?
```

### Creating a model

Create a `Modelfile`:

```
FROM llama2
PROMPT """
You are super mario from super mario bros. Answer instructions as Mario only.

User: {{ .Prompt }}
Assitant:
"""
```

Next, create and run the model:

```
ollama create mario -f Modelfile
ollama run mario
>>> hi
Hello! It's your friend Mario.
```

## Install

- [Download](https://ollama.ai/download) for macOS (Apple Silicon)
- Download for Windows (coming soon)

## Model library

Ollama includes a library of open-source, pre-trained models. More models are coming soon.

| Model       | Parameters | Size  | Download                                 |
| ----------- | ---------- | ----- | ---------------------------------------- |
| Llama2   | 7B         | 3.5B | `ollama pull llama`        |
| Orca Mini   | 3B         | 1.9GB | `ollama pull orca`        |
| Vicuna      | 7B         | 3.9GB | `ollama pull vicuna`      |
| Nous-Hermes | 13         | 7.2GB | `ollama pull hous-hermes` |

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
./ollama run llama2
```
