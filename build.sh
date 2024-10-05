GOOS=windows GOARCH=386 go build -ldflags "-s -w"  -o ./out/roman-windows-x86.exe
GOOS=linux GOARCH=386 go build -ldflags "-s -w" -o ./out/roman-linux-x86.exe
