.PHONY: build
build:
	docker build -t goji:latest .

.PHONY: up
up:
	docker run -p 8080:8080 goji:latest
