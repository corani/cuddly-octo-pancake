# Testing GitHub Models

Test repo to use [GitHub Models API](https://docs.github.com/en/rest/models) from a GitHub Action using a simple Go client.

This was inspired by a [PR](https://github.com/tmc/langchaingo/pull/1258) for langchaingo, from where I borrowed some of the client code.

>[!warning]
> Do not use this, it's just a demo 😊

## Usage

Run the application with the following options:

```sh
go run . [flags]
```

### Flags

- `-models`   Print available models and exit.
- `-message`  Specify the user message for chat completion (default: "What is the capital of France?").
- `-help`     Show help message and usage.

### Examples


Print available models:

```console
$ go run . -models

Available models:
Model: openai/gpt-4.1, Tags: [default]
Model: openai/gpt-3.5, Tags: [fast]
...
```

Perform a chat completion with a custom message:

```console
$ go run . -message "Tell me a joke."

system: You are a helpful assistant.
user: Tell me a joke.
assistant: Why did the scarecrow win an award? Because he was outstanding in his field!
```

Show help:

```console
$ go run . -help

Usage of cuddly-octo-pancake:
  -help
        Show help message and usage.
  -message string
        User message for chat completion (default "What is the capital of France?")
  -models
        Print available models and exit
```
