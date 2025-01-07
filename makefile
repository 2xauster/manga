start:
	./build/manga

build:
	@echo "Building project..."
	go build -v -buildvcs=false -o=build/manga
	@echo "Project built"

test:
	go test

dev:
	go run . --isDev 

reset:
	@echo "Resetting build..."
	rm -r build
	@echo "build directory removed"

tidy:
	go mod tidy