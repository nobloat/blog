BLOG_SRC := content/*.md style.css main.go
OUTPUT_DIR := public

.PHONY: build deploy dev

build: $(BLOG_SRC)
	mkdir -p $(OUTPUT_DIR)
	go run main.go

deploy: build
	#rsync -av --delete $(OUTPUT_DIR)/ $(DEPLOY_DIR)/

dev:
	go run main.go -watch

clean:
	rm -rf $(OUTPUT_DIR)
