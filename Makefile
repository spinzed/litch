output: *.go
	go build
	cp ./litch /bin
	cp ./litch /usr/bin
	rm litch

uninstall:
	rm -f /bin/litch
	rm -f /usr/bin/litch
