build:
	@echo "Building project..."
	go build -v -buildvcs=false -o=build/manga
	@echo "Project built"

start:
	./build/manga
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