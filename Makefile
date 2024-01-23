run:
	@go run main.go

hot:
	echo "Running hot reload mode"
	templ generate --watch --proxy="http://localhost:3000" --cmd="go run main.go"