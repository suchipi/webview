On windows, the webview will not be able to access localhost unless the following command has been run as administrator at least once before:

```
CheckNetIsolation.exe LoopbackExempt -a -n=Microsoft.Win32WebViewHost_cw5n1h2txyewy
```

If you're using the webview to show local files, you'll need that, because the program runs an http server on localhost in order to show them.

Running `launcher.exe` will check if this has been run before, request admin permissions and run it if it hasn't, and then run `webview.exe`.

Any flags you pass to `launcher.exe` will get forwarded to `webview.exe`.

So, if you're distributing an application, your shortcut should point to launcher.exe, not webview.exe.
