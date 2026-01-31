# gosh 
[![CI](https://github.com/archstrap/gosh/actions/workflows/ci.yml/badge.svg)](https://github.com/archstrap/gosh/actions/workflows/ci.yml)
[![Auto Tag](https://github.com/archstrap/gosh/actions/workflows/auto-tag.yml/badge.svg)](https://github.com/archstrap/gosh/actions/workflows/auto-tag.yml)


**gosh** is a minimal, interactive Unix-style shell written in Go. It provides a read–eval–print loop (REPL) with built-in commands, piping, I/O redirection, and raw terminal handling.

---

## Tags

`go` · `shell` · `repl` · `terminal` · `cli` · `unix` · `parser` · `builtins`

---

## Features

### Core

- **Interactive REPL** — Raw terminal mode with prompt (configurable via `PS` env var).
- **Built-in commands** — `cd`, `pwd`, `echo`, `exit`, `type`, `history`.
- **External programs** — Run any executable from `PATH`.
- **Piping** — Chain commands with `|` (e.g. `ls | head -5`).
- **I/O redirection** — `<` and `>` for stdin/stdout (including `2>` for stderr).

### UX

- **History** — Up/Down arrows, persisted via `HISTFILE`.
- **Tab completion** — Builtins and executables; double-tab lists options.
- **Ctrl+C** — Interrupt current line.
- **Ctrl+D** — Exit (after saving history).
- **`.shellrc`** — Optional config file loaded at startup.

### Implementation

- **Parser** — Tokenizer/lexer with quoted strings (`'` and `"`), escapes, and redirects.
- **No forking for builtins** — Builtins run in-process; external commands via `exec`.
- **Structured commands** — Parsed into commands with args and redirects, then executed in order with pipes.

---

## Install

### Prerequisites

- **Go 1.25+** — [Install Go](https://go.dev/doc/install) if needed.
- Unix-like environment (Linux, macOS, WSL).

### From source

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/gosh.git
   cd gosh
   ```

2. **Build the binary**

   ```bash
   make build
   ```

   Or without Make:

   ```bash
   go build -o gosh app/*.go
   ```

   This produces a `gosh` binary in the current directory.

3. **(Optional) Install to your PATH**

   - **User install** (recommended): copy into a directory that’s on your `PATH`, e.g. `~/bin`:

     ```bash
     mkdir -p ~/bin
     cp gosh ~/bin/
     ```

     Ensure `~/bin` is in your `PATH` (e.g. add `export PATH="$HOME/bin:$PATH"` to your shell config).

   - **System-wide**: install to `/usr/local/bin` (may require `sudo`):

     ```bash
     sudo cp gosh /usr/local/bin/
     ```

4. **Run gosh**

   ```bash
   ./gosh
   ```

   Or, if you installed it to your PATH:

   ```bash
   gosh
   ```

### From a release (GitHub)

Push a tag (e.g. `v1.0.0`) to trigger a [GitHub Release](https://github.com/yourusername/gosh/releases) with a Linux binary attached. Download `gosh-linux-amd64` from the latest release.

---

## CI / Releases

- **Runs on:** push to `main`, pull requests targeting `main`, and tag pushes `v*`.
- **Artifacts:** Each run builds the binary; download it from the [Actions](https://github.com/yourusername/gosh/actions) run summary (Artifacts).
- **Releases:** Pushing a tag (e.g. `git push origin v1.0.0`) creates a GitHub Release and attaches the built binary.

---

## Development

| Command            | Description                    |
|--------------------|--------------------------------|
| `make build`       | Build `gosh` binary            |
| `make run`         | Build and run `gosh`           |
| `make test`        | Run tests                      |
| `make test-coverage` | Tests + HTML coverage report |
| `make lint`        | Run golangci-lint (if installed) |
| `make help`        | List all targets               |

---

## Project structure

```
.
├── app/
│   ├── main.go      # Entry point, REPL loop, raw terminal
│   ├── parser.go    # Tokenizer & command parser
│   ├── command.go   # Builtins (cd, pwd, echo, type, exit, history)
│   ├── execute.go   # Command execution, piping, redirects
│   ├── trie.go      # Tab completion (Trie)
│   ├── history.go   # History storage and navigation
│   ├── file.go      # File/executable lookup
│   ├── setup.go     # .shellrc loading
│   ├── color.go     # Color helpers
│   └── util.go      # Shared utilities
├── .shellrc         # Optional shell config
├── Makefile
├── go.mod
└── README.md
```

---

## Resources

- [os/exec patterns (Go)](https://www.dolthub.com/blog/2022-11-28-go-os-exec-patterns/)
- [Beej's Guide to Unix IPC](https://beej.us/guide/bgipc/)

---

## License

See repository for license information.
