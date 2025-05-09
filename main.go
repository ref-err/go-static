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
	"github.com/ref-err/go-static/config"
	"github.com/ref-err/go-static/handler"
	"github.com/skip2/go-qrcode"
)

var version = "dev"

// creating http server
var server = &http.Server{
	Addr:           ":8080",
	Handler:        nil,
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
	var startTime = time.Now()

	port := flag.Int("port", 8080, "Server HTTP port")                                    // port from cmd-line args
	root := flag.String("root", ".", "Serve files from this directory")                   // root dir from cmd-line args
	versionFlag := flag.Bool("version", false, "Prints version")                          // self-explanatory
	openFlag := flag.Bool("open", false, "Opens file server in default browser")          // open browser
	qrFlag := flag.Bool("qr", false, "Show QR-code that redirects to localhost:YOURPORT") // show qr
	configPath := flag.String("config-file", "", "Path to the config file")               // path to config file

	flag.Parse()

	if *configPath != "" {
		config, err := config.LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("Error loading config file: %v", err)
		}

		root = &config.Root
		port = &config.Port
		openFlag = &config.Open
		qrFlag = &config.Qr
	}

	handler.SetRoot(root)

	if *versionFlag {
		fmt.Printf("go-static %s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	mux := http.NewServeMux()

	handler.RegisterEndpoints(mux, version, startTime) // registering endpoints
	server.Addr = ":" + fmt.Sprint(*port)              // applying port to server
	fs := http.FileSystem(http.Dir(*root))             // applying root dir to http.FileSystem
	mux.Handle("/", handler.FileServer(fs))            // using mux to handle everything
	server.Handler = mux                               // overriding handler

	log.Printf("Running go-static %s on %s %s\n", version, runtime.GOOS, runtime.GOARCH)
	log.Println("Server listening at http://localhost:" + fmt.Sprint(*port))
	log.Println("Root directory:", *root)

	if *qrFlag {
		generateQR(port)
	}

	if *openFlag {
		err := browser.OpenURL("http://localhost:" + fmt.Sprint(*port))
		if err != nil {
			log.Fatalln("Failed to open URL: ", err)
		}
	}

	log.Fatal(server.ListenAndServe()) // run server
}

func generateQR(port *int) {
	log.Println("Your QR-code:")
	qr, err := qrcode.New("http://localhost:"+fmt.Sprint(*port), qrcode.Low)
	if err != nil {
		panic(err)
	}
	qr.DisableBorder = true

	fmt.Println(qr.ToString(true))
}
