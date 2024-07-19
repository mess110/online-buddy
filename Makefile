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