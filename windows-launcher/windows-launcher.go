package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/jessevdk/go-flags"
	"golang.org/x/sys/windows"
)

func runProgram(exe string, args string, elevated bool) error {
	cwd, _ := os.Getwd()

	var verb string
	if elevated {
		verb = "runas"
	} else {
		verb = "open"
	}

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 // SW_NORMAL

	return windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
}

func errorBox(label string, errm error) {
	message := fmt.Sprintf(label, errm.Error())
	messagePtr, err := syscall.UTF16PtrFromString(message)
	if err != nil {
		panic(err)
	}

	_, err = windows.MessageBox(0, messagePtr, nil, 0)
	if err != nil {
		panic(err)
	}

	os.Exit(1)
}

func infoBox(caption string, message string) {
	captionPtr, err := syscall.UTF16PtrFromString(caption)
	if err != nil {
		panic(err)
	}
	messagePtr, err := syscall.UTF16PtrFromString(message)
	if err != nil {
		panic(err)
	}

	_, err = windows.MessageBox(0, messagePtr, captionPtr, 0)
	if err != nil {
		panic(err)
	}
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
		errorBox("Error trying to get executable path: %s", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		errorBox("Error trying to evaluate symlinks for executable: %s", err)
	}

	dir := filepath.Dir(exe)
	return dir
}

func runWebview(exePath string) {
	var webviewExe string
	if exePath == "" {
		dir := executableDir()
		webviewExe = filepath.Join(dir, "webview.exe")
	} else {
		if strings.HasPrefix(exePath, "./") || strings.HasPrefix(exePath, ".\\") || strings.HasPrefix(exePath, "../") || strings.HasPrefix(exePath, ".\\") {
			webviewExe = filepath.Join(executableDir(), exePath)
		} else {
			webviewExe = exePath
		}
	}

	if fileExists(webviewExe) {
		cmd := exec.Command(webviewExe, os.Args[1:]...)
		cmd.Run()
	} else {
		errorBox("Unable to start webview: %s", fmt.Errorf("File not found: %s", webviewExe))
	}
}

type options struct {
	// Stuff from webview. We have them all here so they show up in the help text
	Dir    string `short:"d" long:"dir" description:"path to serve" default:"./static" json:"dir"`
	URL    string `long:"url" description:"instead of serving files, navigate to this url" default:"" json:"url"`
	Title  string `long:"title" description:"title of the webview window" default:"webview" json:"title"`
	Width  int    `long:"width" description:"width of the webview window" default:"800" json:"width"`
	Height int    `long:"height" description:"height of the webview window" default:"600" json:"height"`

	// Stuff only the launcher uses
	Exe string `long:"webviewExe" decsription:"path to the webview exe to run" json:"webviewExe"`
}

func loadJSONConfig(opts *options) {
	maybeJSONFile := filepath.Join(executableDir(), "launcher.json")
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
	// We copy the flags from webview.go so they're shown in help
	var opts options

	parser := flags.NewParser(&opts, flags.Default|flags.IgnoreUnknown)
	_, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	loadJSONConfig(&opts)

	// Get a list of currently-exempt apps
	cmd := exec.Command("checknetisolation", "LoopbackExempt", "-s")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		// If checknetisolation isn't present, maybe we're in wine or something?
		// Just launch the app normally and hope it works.
		runWebview(opts.Exe)
		return
	}
	outputstr := string(output)

	loopbackIsExempt := strings.Contains(outputstr, "microsoft.win32webviewhost_cw5n1h2txyewy")
	fmt.Println("loopbackIsExempt", loopbackIsExempt)
	if loopbackIsExempt {
		// We're already exempt, so we can launch the app normally.
		runWebview(opts.Exe)
	} else {
		infoBox("Permission required", fmt.Sprintf("To launch %s, we need to add a Loopback Exemption to your Network Isolation settings, which will allow the program to access localhost. Administrator permission will be requested in order to do this.", opts.Title))
		runProgram("C:\\Windows\\System32\\CheckNetIsolation.exe", "LoopbackExempt -a -n=Microsoft.Win32WebViewHost_cw5n1h2txyewy", true)

		// Now we should be able to run the webview...
		runWebview(opts.Exe)
	}
}
