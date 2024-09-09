# webview

A cross-platform program that launches a webview pointed at either a URL or files on disk (in which case it spawns a local http server for them on an automatically-allocated port).

A tiny wrapper around https://github.com/webview/webview. On macOS, it uses Cocoa/WebKit, on Windows 10, it uses Edge, and on Linux/FreeBSD, it uses gtk-webkit2 (so it depends on GTK3 and GtkWebkit2 on Linux; Ubuntu users can `sudo apt-get install libwebkit2gtk-4.0-dev`).

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

Instead of passing CLI flags, you can also create a `webview.json` file and put it next to the webview binary. If present, it will be used in combination with the command-line arguments.

Example `webview.json`:
```
{
  "url": "https://google.com",
  "title": "Some Site",
  "width": 640,
  "height": 480
}
```

In the case that there is a conflict between the CLI flags and the JSON, the value specified by the JSON will be used instead of the value specified by the flag.

## Notes

Running `window.close()` from JavaScript will close the webview window and exit the program.

Uses [webview/webview](https://github.com/webview/webview).
