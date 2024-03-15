run:
	templ generate
	@go run main.go fake_contacts.go

hot:
	echo "Running hot reload mode"
	templ generate --watch --proxy="http://localhost:3000" --cmd="go run main.go fake_contacts.go"