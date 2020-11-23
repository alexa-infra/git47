
GO_SRC := $(shell find . -name "*.go")
FA_DIR := node_modules/@fortawesome/fontawesome-free/webfonts
FA_FONTS_SOURCE := $(wildcard $(FA_DIR)/*.ttf) $(wildcard $(FA_DIR)/*.woff) $(wildcard $(FA_DIR)/*.woff2) $(wildcard $(FA_DIR)/*.eot)
FA_FONTS_TARGET := $(patsubst $(FA_DIR)/%, static/webfonts/%, $(FA_FONTS_SOURCE))

git47: $(GO_SRC)
	go build backend/main/main.go

check:
	gofmt -l -w **/*.go
	golint ./...

clean:
	rm git47
	rm -r static

static/css/styles.css: css/styles.css css/content.css
	npx postcss css/styles.css -o $@

static/css/fontawesome.css: css/fontawesome.css
	npx postcss $< -o $@

static/favicon.ico: css/favicon.ico
	@test -d static || mkdir static
	cp $< $@

static/webfonts/%: $(FA_DIR)/%
	@test -d static/webfonts || mkdir -p static/webfonts
	cp $< $@

static: static/css/styles.css static/css/fontawesome.css static/favicon.ico $(FA_FONTS_TARGET)

all: git47 static

dev:
	modd
