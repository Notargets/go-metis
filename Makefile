.PHONY: deps test clean install-metis

# Install METIS and dependencies
install-metis:
	# Clone METIS and its dependency GKlib from the official repos
	git clone https://github.com/KarypisLab/GKlib.git
	git clone https://github.com/KarypisLab/METIS.git
	# Build GKlib first (METIS depends on it)
	cd GKlib && make config prefix=/usr/local && sudo make install
	# Then build METIS
	cd METIS && make config shared=1 prefix=/usr/local gklib_path=/usr/local && sudo make install
	sudo ldconfig
	# Add to your environment (add to ~/.bashrc or ~/.zshrc)
	@echo "Adding environment variables to ~/.bashrc"
	@echo 'export LD_LIBRARY_PATH=/usr/local/lib:$$LD_LIBRARY_PATH' >> $(HOME)/.bashrc
	@echo 'export CGO_CFLAGS="-I/usr/local/include"' >> $(HOME)/.bashrc
	@echo 'export CGO_LDFLAGS="-L/usr/local/lib"' >> $(HOME)/.bashrc
	@echo "Please run 'source ~/.bashrc' or restart your shell"

# Install Go dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the library
build:
	go build -v ./...

# Clean build artifacts
clean:
	rm -rf GKlib METIS
	rm -f coverage.out coverage.html
	go clean

# Format code
fmt:
	go fmt ./...
	gofmt -s -w .

# Run linters
lint:
	golangci-lint run

# Install development tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Full setup: install METIS and Go dependencies
setup: install-metis deps

# Check if METIS is installed
check-metis:
	@echo "Checking METIS installation..."
	@ldconfig -p | grep metis || (echo "METIS not found in library path" && exit 1)
	@ls -la /usr/local/include/metis.h || (echo "metis.h not found" && exit 1)
	@echo "METIS appears to be installed correctly"
