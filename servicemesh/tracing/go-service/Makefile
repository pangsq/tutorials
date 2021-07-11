
build:
	CGO_ENABLED=0 GOOS=linux go build -o go-service -a -ldflags '-extldflags "-static"' main.go

image: build
	docker build -t go-service:edge .