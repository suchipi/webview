cd webview
go build -ldflags="-H windowsgui" -o ..\bin\windows-amd64\webview.exe .
cd ..

cd windows-launcher
go build -ldflags="-H windowsgui" -o ..\bin\windows-amd64\launcher.exe .
cd ..
