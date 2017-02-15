osx: 
	env GOOS=darwin GOARCH=amd64 go build -o souparchive_osx_64
linux:
	env GOOS=linux GOARCH=amd64 go build -o souparchive_linux_64
win:
	env GOOS=windows GOARCH=amd64 go build -o souparchive_win_64.exe
clean:
	rm -f souparchive_osx_64 souparchive_linux_64 souparchive_win_64.exe

release: osx linux win

.PHONY: clean
