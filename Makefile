BLOG_SRC := articles/*.md style.css main.go
OUTPUT_DIR := public
DEPLOY_DIR := kadse@jedicke.uberspace.de:web/nobloat.org
IMAGES_DIR := $(OUTPUT_DIR)/images

.PHONY: build deploy dev clean convert

build: $(BLOG_SRC)
	mkdir -p $(OUTPUT_DIR)
	go run ./...

deploy: build
	rsync -av --delete $(OUTPUT_DIR)/ $(DEPLOY_DIR)/


dev:
	go run -tags watch ./... -watch

clean:
	rm -rf $(OUTPUT_DIR)
