
GO_SRC := $(shell find . -name "*.go")
FA_DIR := node_modules/@fortawesome/fontawesome-free/webfonts
FA_FONTS_SOURCE := $(wildcard $(FA_DIR)/*.ttf) $(wildcard $(FA_DIR)/*.woff) $(wildcard $(FA_DIR)/*.woff2) $(wildcard $(FA_DIR)/*.eot)
FA_FONTS_TARGET := $(patsubst $(FA_DIR)/%, internal/web/static/webfonts/%, $(FA_FONTS_SOURCE))

git47: $(GO_SRC)
	go build -o git47 main.go

check:
	gofmt -l -w ./internal/**/*.go main.go
	golint ./...

clean:
	rm git47
	rm -r internal/web/static

internal/web/static/css/styles.css: css/styles.css css/content.css
	npx postcss css/styles.css -o $@

internal/web/static/css/fontawesome.css: css/fontawesome.css
	npx postcss $< -o $@

internal/web/static/favicon.ico: css/favicon.ico
	@test -d internal/web/static || mkdir internal/web/static
	cp $< $@

internal/web/static/webfonts/%: $(FA_DIR)/%
	@test -d internal/web/static/webfonts || mkdir -p internal/web/static/webfonts
	cp $< $@

static: internal/web/static/css/styles.css internal/web/static/css/fontawesome.css internal/web/static/favicon.ico $(FA_FONTS_TARGET)

all: static git47

dev:
	modd
