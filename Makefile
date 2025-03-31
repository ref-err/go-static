BUILD_DIR = ./build
BINARY_NAME= go-static
FLAGS = CGO_ENABLED=0
RUN_FLAGS = --port=8080 --root=$(HOME)

help:
	@echo "Usage:"
	@echo "make build - compile binary"
	@echo "make run   - run compiled binary"
	@echo "make clean - remove build artifacts"

build:
	@mkdir -p $(BUILD_DIR)
	@$(FLAGS) go build -o $(BUILD_DIR)/$(BINARY_NAME)

run: build
	@#echo "Command-line arguments: $(RUN_FLAGS)"
	@$(BUILD_DIR)/$(BINARY_NAME) $(RUN_FLAGS)

clean:
	@rm -rf $(BUILD_DIR)

.PHONY: help build run clean