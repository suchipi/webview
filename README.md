# static-webview

A cross-platform program that serves files from disk and shows them in a native webview.

On macOS, it uses Cocoa/WebKit, on Windows 10, it uses Edge, and on Linux/FreeBSD, it uses gtk-webkit2 (so it depends on GTK3 and GtkWebkit2 on Linux).

## Usage

```
static-webview [options]

  -dir string
        path to serve (default "./static")
  -height int
        height of the webview window (default 600)
  -title string
        title of the webview window (default "webview")
  -width int
        width of the webview window (default 800)
```

## Notes

Running `window.close()` from JavaScript will close the webview window and exit the program.

Uses [webview/webview](https://github.com/webview/webview).
