package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/pkg/browser"
	"github.com/ref-err/go-static/handler"
)

var version = "dev"

// creating http server
var server = &http.Server{
	Addr:           ":8080",
	Handler:        http.FileServer(http.Dir(".")),
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}

func init() {
	// making sure that server shuts down properly
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Stopping...")
		server.Shutdown(context.TODO())
	}()
}

func main() {
	port := flag.Int("port", 8080, "Server HTTP port")                          // port from cmd-line args
	root := flag.String("root", ".", "Serve files from this directory")         // root dir from cmd-line args
	versionFlag := flag.Bool("version", false, "Prints version")                // self-explanatory
	openFlag := flag.Bool("open", false, "Open file server in default browser") // same ^

	flag.Parse()

	if *versionFlag {
		fmt.Printf("go-static %s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	server.Addr = ":" + fmt.Sprint(*port)   // applying port to server
	fs := http.FileSystem(http.Dir(*root))  // applying root dir to http.FileSystem
	server.Handler = handler.FileServer(fs) // overriding handler

	log.Printf("Running go-static %s on %s %s\n", version, runtime.GOOS, runtime.GOARCH)
	log.Println("Server listening at http://localhost:" + fmt.Sprint(*port))
	log.Println("Root directory:", *root)

	if *openFlag {
		err := browser.OpenURL("http://localhost:" + fmt.Sprint(*port))
		if err != nil {
			log.Fatalln("Failed to open URL: ", err)
		}
	}

	log.Fatal(server.ListenAndServe()) // run server
}
