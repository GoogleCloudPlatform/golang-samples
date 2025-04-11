.PHONY: bin cross upload

VERSION=`git rev-parse HEAD`
DATE=`date +%Y-%m-%d_%H:%M`
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildDate=${DATE}"

bin:
	rm -f ${GOPATH}/bin/gimmeproj
	go install ${LDFLAGS}

cross:
	GOOS=linux   GOARCH=amd64 go build ${LDFLAGS} -o bin/linux_amd64/gimmeproj
	GOOS=darwin  GOARCH=amd64 go build ${LDFLAGS} -o bin/darwin_amd64/gimmeproj
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/windows_amd64/gimmeproj.exe

upload: cross
	gsutil -m cp -a public-read -r bin/* gs://gimme-proj/
