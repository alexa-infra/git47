
ROOT_DIR := $(shell pwd)

git47: *.go
	go build

check:
	gofmt -l -w **/*.go
	golint ./...

clean:
	rm git47
	rm -r static

static/output.css: css/styles.css css/content.css
	npx postcss css/styles.css -o static/output.css

static/favicon.ico: css/favicon.ico
	cp css/favicon.ico static/

all: git47 static/output.css static/favicon.ico

dev:
	modd
