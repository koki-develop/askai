<h1 align="center">askai</h1>

<p align="center">
AI is with you.
</p>

<p align="center">
<a href="https://github.com/koki-develop/askai/releases/latest"><img src="https://img.shields.io/github/v/release/koki-develop/askai" alt="GitHub release (latest by date)"></a>
<a href="https://github.com/koki-develop/askai/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/koki-develop/askai/ci.yml?logo=github" alt="GitHub Workflow Status"></a>
<a href="https://codeclimate.com/github/koki-develop/askai/maintainability"><img src="https://img.shields.io/codeclimate/maintainability/koki-develop/askai?style=flat&amp;logo=codeclimate" alt="Maintainability"></a>
<a href="https://goreportcard.com/report/github.com/koki-develop/askai"><img src="https://goreportcard.com/badge/github.com/koki-develop/askai" alt="Go Report Card"></a>
<a href="./LICENSE"><img src="https://img.shields.io/github/license/koki-develop/askai" alt="LICENSE"></a>
</p>

<p align="center">
<img src="./assets/demo.gif" >
</p>

## Contents

- [Installation](#installation)
- [Usage](#usage)
- [LICENSE](#license)

## Installation

### Homebrew Tap

```console
$ brew install koki-develop/tap/askai
```

### `go install`

```console
$ go install github.com/koki-develop/askai@latest
```

### Releases

Download the binary from the [releases page](https://github.com/koki-develop/askai/releases/latest).

## Usage

```console
$ askai --help
AI is with you.

Usage:
  askai [flags] [question]

Flags:
  -k, --api-key string   the OpenAI API key
      --configure        configure askai
  -g, --global           configure askai globally (only for --configure)
  -h, --help             help for askai
  -i, --interactive      interactive mode
  -m, --model string     the chat completion model to use (default "gpt-3.5-turbo")
  -v, --version          version for askai
```

## LICENSE

[MIT](./LICENSE)
