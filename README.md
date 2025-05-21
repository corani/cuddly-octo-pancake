# Testing GitHub Models

Test repo to use [GitHub Models API](https://docs.github.com/en/rest/models) from a GitHub Action using a simple Go client.

This was inspired by a [PR](https://github.com/tmc/langchaingo/pull/1258) for langchaingo, from where I borrowed some of the client code.

>[!warning]

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

```sh
go run . -models
```

Perform a chat completion with a custom message:

```sh
go run . -message "Tell me a joke."
```

Show help:

```sh
go run . -help
```
