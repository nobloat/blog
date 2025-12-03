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

test:


dev:
	go run ./... -watch

clean:
	rm -rf $(OUTPUT_DIR)

convert:
	@if [ -z "$(file)" ]; then \
		echo "Usage: make convert file=path/to/image.jpg [out=name.png]"; \
		exit 1; \
	fi
	mkdir -p $(IMAGES_DIR)
	go run ./... image "$(file)" "$(or $(out),$(IMAGES_DIR)/$(notdir $(basename $(file)))_gray.png)"
