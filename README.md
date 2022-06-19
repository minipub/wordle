# Wordle!

A Wordle Game.

## Features

1. `wordle-puzzle` support classic inactive shell mode

2. `wordle` support C-S mode

3. `wordle-solver` support solving puzzle automatically (Todo)

## Installation

### Binaries

Download the binary from the [releases](https://github.com/minipub/wordle/releases) page.

### From Source

##### Requirements

* Go1.18+

If you already have the Go1.18 SDK installed, you can use the go tool to install wordle-puzzle:

`go install github.com/minipub/wordle/cmd/wordle-puzzle@latest`

## Usage

### C-S Mode

Run server

```
wordle-puzzle
```

Run client

You can run any command tool to connect local tcp port 8080, eg: telnet

```
telnet localhost 8080
```

If it work well, you will see the following response

```
root@orz ~ $ socat tcp:localhost:8080 -

 _       __               ____        __
| |     / /___  _________/ / /__     / /
| | /| / / __ \/ ___/ __  / / _ \   / /
| |/ |/ / /_/ / /  / /_/ / /  __/  /_/
|__/|__/\____/_/   \__,_/_/\___/  (_)


Please input a five-letter word and Press <Enter> to confirm.

input:

```

Help yourself to input and enjoy the fun!
