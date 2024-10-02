OUTPUT = game

.PHONY: build
build: 
	@echo "Building..."
	@mkdir -p bin
	@go build -o ./bin/$(OUTPUT) ./cmd/game
	@echo "Build complete"

.PHONY: run
run:
	@./bin/$(OUTPUT)

.PHONY: dev
dev: build
	@./bin/$(OUTPUT)

.PHONY: clean
clean:
	@rm -f ./bin/$(OUTPUT)
