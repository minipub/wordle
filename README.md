# Wordle!

A Command Line Wordle Game.

[What's wordle?](https://en.wikipedia.org/wiki/Wordle)

## How To Play

Players have six attempts to guess a five-letter word

After every attempt, each letter is marked as either Green, Yellow or Red

Green: letter is correct, and in the correct position

Yellow: it's in the answer, but not in the right position

Red: it's not in the answer at all

**Reference from wiki**

## Have A Try

## Features

1. `wordle puzzle` support classic inactive shell mode

2. `wordle puzzle` support C-S mode

3. `wordle solver` support solving puzzle automatically

## Installation

### Binaries

Download the binary from the [releases](https://github.com/minipub/wordle/releases) page.

### From Source

##### Requirements

* Go1.18+

If you already have the Go1.18 SDK installed, you can use the go tool to install

Install wordle:

```
go install github.com/minipub/wordle@latest
```

## Usage

### Requirements

* iTerm2 or any other terminal that support 256 Colors

### Interactive Mode

```
wordle puzzle
```

### C-S Mode

Run server

```
wordle puzzle --mode 2
```

Run client

You can run any command line tool to connect local tcp port 8080, eg: telnet

```
telnet localhost 8080
```

If it works well, you will see the following response

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

Run auto solver

Sometimes you lost your patience or just get bored with this game, you can call the auto solver to help (cheat)

```
wordle solver 2>/dev/null
```

Also, you can redirect Stderr to a file to look the whole analysis

Help yourself to input and enjoy the fun!
