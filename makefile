FOLDER = "example"
# set env value NAME
NAME = "main"
# set env value VERSION
VERSION = "1.0.0"

build:
	go build -o bin/$(NAME) -ldflags "-X main.Version=$(VERSION)" $(FOLDER)/$(NAME).go

run:
	make build
	./bin/$(NAME)