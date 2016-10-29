# pick-up-api
Pick Up App

## Running Pickup Locally
```
# Run from the root of this directory
docker build --rm -t pickup-pg config/db
docker build --rm -t pickup-go-server -f config/api/Dockerfile .
docker run -p 5432:5432 --name pickup-db -e POSTGRES_PASSWORD=pickEmUp -d pickup-pg
docker run -p 3000:3000 -v $(pwd):/go/src/github.com/pick-up-api --name pickup-server --link pickup-db:pickup-postgres -d pickup-go-server
```

## Local Server
Local API server should be available at <docker-ip>:3000

Confirm its available at <docker-ip>:3000/hello/world

## Running Test Suite
```
docker exec pickup-server go get github.com/stretchr/testify/assert
docker exec pickup-server go test  ./tests/utils/validation/strings_test.go
```

### Helpful articles
* [Golang Docker](https://blog.golang.org/docker)
* [Dockerizing Go w/ Local Filesync](https://medium.com/developers-writing/docker-powered-development-environment-for-your-go-app-6185d043ea35#.r58sq9cr2)
* [Architecting Go Web Apps](https://larry-price.com/blog/2015/06/25/architecture-for-a-golang-web-app)
