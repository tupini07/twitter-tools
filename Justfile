# shows this message and lists commands:
help:
    @just --list

# applies code checking and linting
lint: 
    staticcheck ./...
    go vet ./...

# builds the project
build: lint
    go build

# builds and then runs the executable
run: 
    go run .

# adds missing dependencies to go.mod
tidy:
    go mod tidy

# updates dependencies
update-deps: 
    go get -u all
    just tidy