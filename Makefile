all:
	docker network create app
	docker run --network app --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres -p 5432:5432 -d postgres
	docker build -t to_postgres -f scripts/to_postgres/Dockerfile scripts/to_postgres
	docker run --network app to_postgres
	docker build -t app -f app/Dockerfile app 
	docker run -it --network app -p 8080:8080 -d --name app app

clean:
	docker stop app postgres
	docker container remove app postgres
	docker network remove app