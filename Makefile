all: ;

fmt:
	go fmt ./...

test:
	go test ./...

clean:
	rm -rf build/
	rm -f _docker/*/main

depend:
	go get -u -v github.com/Masterminds/glide
	glide install


# Workers

worker-api:
	SESSION_AUTHENTICATION_KEY=abf804e906eaa13c158d9ce721f87502a87f7d4615726f8dde07dd4d3bc173167aead9d52692c237469acacbbd0ac44fcc85b43332b7df83ccb8175ba976b4de \
	SESSION_ENCRYPTION_KEY=fd80a05bad70f666c16669505556ca2f \
	go run worker/api/*.go

worker-fixer-php:
	go run worker/fixer/php/*.go

worker-status-api:
	go run worker/status_api/*.go

worker-result-store:
	go run worker/result_store/*.go


# Utils

ssh-tunnel:
	# echo "GatewayPorts yes" >> /etc/ssh/sshd_config
	@echo "routes 89.221.208.88:8100 to localhost:8100"
	ssh -N -R 8100:localhost:8100 root@89.221.208.88


# Services

services:
	docker-compose -f ./docker-compose.services.yml up -d

services-stop:
	docker-compose -f ./docker-compose.services.yml stop

services-down:
	docker-compose -f ./docker-compose.services.yml down


# Workers

workers:
	docker-compose -f ./docker-compose.workers.yml up -d

workers-build:
	docker-compose -f ./docker-compose.workers.yml up --build -d

workers-stop:
	docker-compose -f ./docker-compose.workers.yml stop

workers-down:
	docker-compose -f ./docker-compose.workers.yml down

workers-logs:
	docker-compose -f ./docker-compose.workers.yml logs -f


# Docker build

build-worker-api:
	CGO_ENABLED=0 GOOS=linux go build -v -o build/worker-api worker/api/*.go

build-worker-fixer-php:
	CGO_ENABLED=0 GOOS=linux go build -v -o build/worker-fixer-php worker/fixer/php/*.go

build-worker-result-store:
	CGO_ENABLED=0 GOOS=linux go build -v -o build/worker-result-store worker/result_store/*.go

build-worker-status-api:
	CGO_ENABLED=0 GOOS=linux go build -v -o build/worker-status-api worker/status_api/*.go

build-workers:
	docker rm -f go-build-workers-container || true
	docker build -t go-build-workers-image -f ./Dockerfile.build-go .
	docker run --name go-build-workers-container go-build-workers-image \
		make clean depend \
		build-worker-api \
		build-worker-fixer-php \
		build-worker-status-api \
		build-worker-result-store
	docker cp go-build-workers-container:./go/src/github.com/mstovicek/checkit/build/worker-api ./_docker/api/main
	docker cp go-build-workers-container:./go/src/github.com/mstovicek/checkit/build/worker-fixer-php ./_docker/fixer-php/main
	docker cp go-build-workers-container:./go/src/github.com/mstovicek/checkit/build/worker-result-store ./_docker/result-store/main
	docker cp go-build-workers-container:./go/src/github.com/mstovicek/checkit/build/worker-status-api ./_docker/status-api/main
	docker rm -f go-build-workers-container


# ELK

elk:
	docker-compose -f ./_docker-elk/docker-compose.yml up -d

elk-logs:
	docker-compose -f ./_docker-elk/docker-compose.yml logs -f

elk-stop:
	docker-compose -f ./_docker-elk/docker-compose.yml stop

elk-down:
	docker-compose -f ./_docker-elk/docker-compose.yml down
