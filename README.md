# glox - A lox interpreter written in Go

## Build and run immediately

```bash
# Run the REPL
go run src/cmd/glox.go

# Run a .lox source code file
go run src/cmd/glox.go <path_to_script>
```

## Build

```bash
# Create a binary in bin/
go build -o bin/ src/cmd/glox.go

# Run the binary
./bin/glox
```

## Build and install

```bash
# Build binary and store in GOPATH
go install src/cmd/glox.go

# Run from GOPATH (ensure GOPATH is set in $PATH)
glox
```
