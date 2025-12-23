# Audio parser

A small go programm to extract interesting information from music files

## Features

- [ ] Parse mp3's
- [ ] Extract the sound curve


## Usage

### Running the CLI

From the project root:

```bash
cd src/cmd/cli
go run main.go <path-to-audio-file>
```

Example:

```bash
cd src/cmd/cli
go run main.go ../../examples/jazz.mp3
```

### Building the WASM Package

From the project root:

```bash
cd src/cmd/wasm
GOOS=js GOARCH=wasm go build -o wasm.wasm main.go
```

The compiled WASM binary will be created as `wasm.wasm` in the `src/cmd/wasm` directory.


## Ideas