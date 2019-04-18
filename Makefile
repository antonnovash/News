build:
	go build -o bin/news-reader main.go

run:
	bin/news-reader

dockerbuild:
	docker build -t sportsnews .

dockerrun:
	docker run -p 3306:3306 --name sportsnews -e MYSQL_ROOT_PASSWORD=1234567 -d mysql:latest
