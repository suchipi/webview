# webview

A cross-platform program that launches a webview pointed at either a URL or files on disk (in which case it spawns a local http server for them on an automatically-allocated port).

On macOS, it uses Cocoa/WebKit, on Windows 10, it uses Edge, and on Linux/FreeBSD, it uses gtk-webkit2 (so it depends on GTK3 and GtkWebkit2 on Linux; Ubuntu users can `sudo apt-get install libwebkit2gtk-4.0-dev`).

## Installation

Get the binaries from the `bin` folder in this repo.

## Usage

```
webview [options]

  --dir string
        path to serve (default ".")
  --url string
        instead of serving files, load this url
  --title string
        title of the webview window (default "webview")
  --width int
        width of the webview window (default 800)
  --height int
        height of the webview window (default 600)
```

## Notes

Running `window.close()` from JavaScript will close the webview window and exit the program.

Uses [webview/webview](https://github.com/webview/webview).
