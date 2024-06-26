# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
MAIN_FILE = cmd/app/main.go
BINARY_NAME = GoSt
BINARY_UNIX = $(BINARY_NAME)_unix
BINARY_WIN = $(BINARY_NAME)_windows
TEST_PROJECT = goat
TEST_PROJECT_DIR = "E:/GoLandProjects/gost/$(TEST_PROJECT)/*"

# Frontend parameters
FRONTEND_DIR = ui/front
NPMCMD = npm
NPMINSTALL = $(NPMCMD) install
NPMRUNBUILD = $(NPMCMD) run build

# All targets
all: test build create new init run release frontend clean

# Test target
test:
	$(GOTEST) -v ./...

# Create target
create:
	$(GORUN) $(MAIN_FILE) create $(TEST_PROJECT)

# New target
new:
	$(GORUN) $(MAIN_FILE) new $(TEST_PROJECT)

# Init target
init:
	$(GORUN) $(MAIN_FILE) init $(TEST_PROJECT)

# Run target
run:
	$(GORUN) $(MAIN_FILE) create $(TEST_PROJECT)

# Build target
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Release target
release: clean
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
	zip $(BINARY_UNIX).zip $(BINARY_UNIX)

# Frontend target
frontend:
	cd $(FRONTEND_DIR) && $(NPMINSTALL) && $(NPMRUNBUILD)

# Clean target
clean:
	rm -f $(TEST_PROJECT_DIR)
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_UNIX).zip
	rm -f $(BINARY_WIN)
	rm -f $(BINARY_WIN).zip

.PHONY: all test build create new init run release frontend clean