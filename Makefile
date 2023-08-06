run:
	go run main.go draw.go sprites.go

install:
	# download needed images from old pokemon repo
	@mkdir -p img
	go run install.go
