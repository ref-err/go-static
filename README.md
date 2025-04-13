# go-static

A simple file server that runs locally on your PC, written in Go.

## Running
To run `go-static`, you can download binary for your system in Releases page. But you can also install latest build by running:
```bash
go install github.com/ref-err/go-static@latest
```

## Usage
### Command-line arguments:
- `--port=PORT`: File server port _(default: 8080)_
- `--root=path/to/dir`: File server root directory _(default: ".")_

### Example
```bash
go-static --port=3000 --root="$HOME/public" 
```

This command runs a file server at port 3000 and a root directory `public` in the user's home.  

## Building
To build this program, you can run:
```bash
go build -o "build/go-static"
```
To build AND run:
```bash
go run . <go-static args> # (for example: --port=3000 --root="$HOME")
```
To build for other OS and other architecture:
```bash
GOOS=windows GOARCH=amd64 go build -o "build/go-static.exe"
```

## License
This program is licensed under the [MIT License](https://opensource.org/license/MIT).