run:
	go run cmd/online-buddy/main.go

build:
	go build cmd/online-buddy/main.go

docker_build:
	docker build --progress=plain --no-cache -t online-buddy .

docker_run:
	docker run online-buddy

docker_compose_all:
	docker-compose up --build

docker_compose:
	docker-compose up --build online-buddy

redis-cli:
	redis-cli

redis-cli-cluster:
	redis-cli -c -p 6079

redis-reset:
	redis-cli -c -p 6079 FLUSHALL
	redis-cli -c -p 6079 CLUSTER RESET
	redis-cli -c -p 6179 FLUSHALL
	redis-cli -c -p 6179 CLUSTER RESET
	redis-cli -c -p 6279 FLUSHALL
	redis-cli -c -p 6279 CLUSTER RESET
	redis-cli -c -p 6080 CLUSTER RESET
	redis-cli -c -p 6180 CLUSTER RESET
	redis-cli -c -p 6280 CLUSTER RESET

redis-list-channels:
	redis-cli -c -p 6079 pubsub channels
	redis-cli -c -p 6179 pubsub channels
	redis-cli -c -p 6279 pubsub channels