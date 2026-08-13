BUILD_DIR := build

.PHONY: all configure build clean run format

all: build

configure:
	cmake -S . -B $(BUILD_DIR)

build: configure
	cmake --build $(BUILD_DIR)

run: build
	$(BUILD_DIR)/clox

clean:
	rm -rf $(BUILD_DIR)

format:
	clang-format -i src/main.c
