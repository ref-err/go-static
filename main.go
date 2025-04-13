package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	port := flag.Int("port", 8080, "File server port")             // port from cmd-line args
	root := flag.String("root", ".", "File server root directory") // root dir from cmd-line args
	versionFlag := flag.Bool("version", false, "Print version")

	flag.Parse()

	if *versionFlag {
		fmt.Println("go-static version", version)
		os.Exit(0)
	}

	server.Addr = ":" + fmt.Sprint(*port)   // applying port to server
	fs := http.FileSystem(http.Dir(*root))  // applying root dir to http.FileSystem
	server.Handler = handler.FileServer(fs) // overriding handler

	log.Println("Running go-static version", version)
	log.Println("Server listening at localhost:" + fmt.Sprint(*port))
	log.Println("Root directory: " + *root)
	log.Fatal(server.ListenAndServe()) // run server
}
