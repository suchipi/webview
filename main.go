package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/webview/webview"
)

type portAndListener struct {
	port     int
	listener net.Listener
}

func getport() portAndListener {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port

	return portAndListener{port, listener}
}

func fileServer(info portAndListener, path string) {
	fs := http.FileServer(http.Dir(path))
	http.Handle("/", fs)

	panic(http.Serve(info.listener, nil))
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
	info := getport()
	path := flag.String("dir", "./static", "path to serve")
	title := flag.String("title", "webview", "title of the webview window")
	width := flag.Int("width", 800, "width of the webview window")
	height := flag.Int("height", 600, "height of the webview window")

	flag.Parse()

	go fileServer(info, *path)
	url := fmt.Sprintf("http://localhost:%d/", info.port)
	runWebview(url, *title, *width, *height)
}
