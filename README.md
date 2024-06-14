<h1 align="center">askai</h1>

<p align="center">
AI is with you.<br>
<small>Powered by <a href="https://openai.com/blog/openai-api">OpenAI API</a></small>
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
- [Getting Started](#getting-started)
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

## Getting Started

### Setup

First, run `askai` with the `--configure` flag set.

```console
$ askai --configure
```

Upon doing this, an interactive setup session will be initiated.  
Input the following items:

- OpenAI API Key: API key created on [the OpenAI platform](https://platform.openai.com).
- Model: The chat completion model to use. (Available: `gpt-4`, `gpt-4-0613`, `gpt-4-32k`, `gpt-4-32k-0613`, `gpt-3.5-turbo`, `gpt-3.5-turbo-0613`, `gpt-3.5-turbo-16k`, `gpt-3.5-turbo-16k-0613`)

Once the setup is completed, a configuration file named `.askai` will be created.  
Now, you are all set.

> [!NOTE]
> If you run `askai --configure` with the `--global` flag set, the configuration file will be created as `$HOME/.askai`.  
> This file is referred to when executing `askai`, if there is no `.askai` in the current directory.

### Ask a question to AI

By passing a question as an argument to `askai`, you can obtain an answer from AI.

```console
$ askai 'Please write a program in Go that outputs "Hello world”.'
```

![](./assets/oneshot.gif)

### Ask a question to AI interactively

By running `askai` with the `-i` or `--interactive` flag set, you can ask questions to AI in an interactive manner.

```console
$ askai -i
```

![](./assets/demo.gif)

### Pass a question from stdin

You can also pass questions from stdin.

```sh
# pipe
$ echo 'QUESTION' | askai

# file
$ askai < question.txt
```

You can pass both input from stdin and arguments at the same time.  
By utilizing this, for example, you can also summarize the contents of a specific file.

```console
$ cat README.md | askai 'Please summarize this content'
```

![](./assets/summarize.gif)

### Ask a question to AI from the Emacs buffer

Please add the following Emacs Lisp code to ~/.emacs.d/init.el.

```Emacs Lisp
(defun askai-on-region (start end)
  (interactive "r")
  (let* ((output-buf (get-buffer-create "*Askai Output*"))
         (err-buf "*Askai Error*")
         (cmd "askai")
	 (cur-buf (current-buffer)))
    (shell-command-on-region start end cmd output-buf nil err-buf)
    (with-current-buffer output-buf
      (let ((str (buffer-substring-no-properties (point-min) (point-max))))
	(with-current-buffer cur-buf
	  (insert "\n\n")
	  (insert str))))))
(global-set-key (kbd "M-c") 'askai-on-region)
```

Please restart Emacs and open the `*scratch*` buffer and please input the following sentence.

Then, please select region the following text in Emacs and press the `ESC-c` key.

```
Please write a program in Go that outputs "Hello world”.
```

The following string will be inserted after the cursor position.

````
You can use the package fmt to print "Hello world" in Go language like this:

```GO
package main

import "fmt"

func main() {
    fmt.Println("Hello world")
}
```

You have to run this program in the environment where Go is installed. You can also use online Go compilers available on the internet if you don't have Go installed in your system.
````


## LICENSE

[MIT](./LICENSE)
