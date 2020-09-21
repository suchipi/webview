package main

import (
	"fmt"
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

func runWebview(exePath string) {
	var webviewExe string
	if exePath == "" {
		exe, err := os.Executable()
		if err != nil {
			errorBox("Error trying to get executable path: %s", err)
		}
		exe, err = filepath.EvalSymlinks(exe)
		if err != nil {
			errorBox("Error trying to evaluate symlinks for executable: %s", err)
		}

		dir := filepath.Dir(exe)
		webviewExe = filepath.Join(dir, "webview.exe")
	} else {
		webviewExe = exePath
	}

	if fileExists(webviewExe) {
		cmd := exec.Command(webviewExe, os.Args[1:]...)
		cmd.Run()
	} else {
		errorBox("Unable to start webview: %s", fmt.Errorf("File not found: %s", webviewExe))
	}
}

func main() {
	var opts struct {
		Title string `long:"title" description:"title of the webview window" default:"webview"`
		Exe   string `long:"webviewExe" decsription:"path to the webview exe to run"`
	}

	parser := flags.NewParser(&opts, flags.Default|flags.IgnoreUnknown)
	_, err := parser.Parse()
	if err != nil {
		panic(err)
	}

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
