.PHONY: run test vet fmt check-format build check

run:
	go run ./cmd/forum

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w .

check-format:
	test -z "$$(gofmt -l .)"

build:
	mkdir -p bin
	go build -o bin/forum ./cmd/forum

.PHONY: seed

seed:
	go run ./cmd/seed

# Docker
DOCKER_IMAGE := literary-lions
DOCKER_TAG := 1.0.0
DOCKER_CONTAINER := literary-lions

docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	docker run --rm -it -p 8080:8080 --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-start:
	docker run -d --name $(DOCKER_CONTAINER) -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-stop:
	docker stop $(DOCKER_CONTAINER)

docker-logs:
	docker logs -f $(DOCKER_CONTAINER)

docker-status:
	docker ps -a

docker-clean:
	docker rm -f $(DOCKER_CONTAINER) 2>/dev/null || true

docker-image-clean:
	docker image rm $(DOCKER_IMAGE):$(DOCKER_TAG)
