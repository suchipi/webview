package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/webview/webview"
)

func getport() (int, net.Listener) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port

	return port, listener
}

func fileServer(listener net.Listener, path string) {
	fs := http.FileServer(http.Dir(path))
	http.Handle("/", fs)

	panic(http.Serve(listener, nil))
}

func close() {
	os.Exit(0)
}

func runWebview(url string, title string, width int, height int) {
	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle(title)
	w.SetSize(width, height, webview.HintNone)
	w.Navigate(url)
	w.Bind("close", close)
	w.Run()
}

func main() {
	port, listener := getport()

	var opts struct {
		Dir    string `short:"d" long:"dir" description:"path to serve" default:"./static"`
		URL    string `long:"url" description:"instead of serving files, navigate to this url" default:""`
		Title  string `long:"title" description:"title of the webview window" default:"webview"`
		Width  int    `short:"w" long:"width" description:"width of the webview window" default:"800"`
		Height int    `short:"h" long:"height" description:"height of the webview window" default:"600"`
	}

	parser := flags.NewParser(&opts, flags.Default|flags.IgnoreUnknown)
	_, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	if opts.URL == "" {
		go fileServer(listener, opts.Dir)
		url := fmt.Sprintf("http://localhost:%d/", port)
		runWebview(url, opts.Title, opts.Width, opts.Height)
	} else {
		runWebview(opts.URL, opts.Title, opts.Width, opts.Height)
	}
}
