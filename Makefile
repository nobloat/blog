BLOG_SRC := articles/*.md style.css main.go
OUTPUT_DIR := public
DEPLOY_DIR := kadse@jedicke.uberspace.de:web/nobloat.org

.PHONY: build deploy dev

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
