.PHONY: all

run:
	go run cmd/tracker/main.go

test:
	go test ./tests

image-build:
	docker build --tag go-task-tracker .

container-run:
	docker run -d -p 7540:7540 -v host-volume:/db go-task-tracker