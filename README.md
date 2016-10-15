# pick-up-api
Pick Up App


## Running Pickup Locally
```
# Run from the root of this directory
docker build --rm -t pickup-pg config/db
docker build --rm -t pickup-go-server -f config/api/Dockerfile .
docker run -p 5432:5432 --name pickup-db -e POSTGRES_PASSWORD=pickEmUp -d pickup-pg
docker run -p 3000:3000 -v $(pwd)/api:/go/src/pickup-server --name pickup-server --link pickup-db:pickup-postgres -d pickup-go-server
```
