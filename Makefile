.DEFAULT_GOAL := hot

.PHONY:fmt vet build run hot
fmt:
	go fmt ./...
vet: fmt
	go vet ./...
build: vet
	go build
run:
	templ generate
	@go run main.go fake_contacts.go

hot: vet
	echo "Running hot reload mode"
	templ generate --watch --proxy="http://localhost:3000" --cmd="go run main.go fake_contacts.go"


.PHONY:clean
clean:
	go clean