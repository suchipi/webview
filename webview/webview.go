package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"

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

type options struct {
	Dir    string `short:"d" long:"dir" description:"path to serve" default:"." json:"dir"`
	URL    string `long:"url" description:"instead of serving files, navigate to this url" default:"" json:"url"`
	Title  string `long:"title" description:"title of the webview window" default:"webview" json:"title"`
	Width  int    `long:"width" description:"width of the webview window" default:"800" json:"width"`
	Height int    `long:"height" description:"height of the webview window" default:"600" json:"height"`
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func executableDir() string {
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		panic(err)
	}

	dir := filepath.Dir(exe)
	return dir
}

func loadJSONConfig(opts *options) {
	maybeJSONFile := filepath.Join(executableDir(), "webview.json")
	if fileExists(maybeJSONFile) {
		jsonFile, err := os.Open(maybeJSONFile)
		if err != nil {
			// Couldn't open the file; just bail.
			return
		}
		defer jsonFile.Close()

		bytes, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			return
		}

		json.Unmarshal(bytes, opts)
	}
}

func main() {
	port, listener := getport()

	var opts options

	parser := flags.NewParser(&opts, flags.Default|flags.IgnoreUnknown)
	_, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	loadJSONConfig(&opts)

	if opts.URL == "" {
		go fileServer(listener, opts.Dir)
		url := fmt.Sprintf("http://localhost:%d/", port)
		runWebview(url, opts.Title, opts.Width, opts.Height)
	} else {
		runWebview(opts.URL, opts.Title, opts.Width, opts.Height)
	}
}
