# go-static

A simple file server that runs locally on your PC, written in Go.

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
This program uses `make` build system. But you can also use `go build` if you don't have `make` installed.  
Full list of `make` tasks:
```bash
make help
```

## License
This program is licensed under the [MIT License](https://opensource.org/license/MIT).