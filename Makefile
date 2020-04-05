
git47: *.go
	go build

check:
	go fmt *.go
	golint *.go

clean:
	rm git47
	rm -r static

static/output.css: styles.css content.css
	npx postcss styles.css -o static/output.css

all: git47 static/output.css
