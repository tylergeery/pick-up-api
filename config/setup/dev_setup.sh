#!/bin/bash
# Run from the root of this directory

# defaults
clean=0

# handle cli flags
while :; do
    case $1 in
        -c|--clean) # clean up existing containers
            clean=1
            break
            ;;
        *)
            echo 'flag $1 not recognized'
            break
            ;;
    esac

    shift
done

#clean up existing containers
if [[ "$clean" = 1 ]]; then
    echo 'Cleaning up existing docker containers'
    eval "docker kill pickup-server"
    eval "docker kill pickup-db"
    eval "docker rm pickup-server"
    eval "docker rm pickup-db"
fi

# create new images and start up local containers
eval "docker build --rm -t pickup-pg config/db"
eval "docker build --rm -t pickup-go-server -f config/api/Dockerfile ."
eval "docker run -p 5432:5432 --name pickup-db -e POSTGRES_PASSWORD=pickEmUp -d pickup-pg"
eval "docker run -p 3000:3000 -v $(pwd):/go/src/github.com/pick-up-api --name pickup-server --link pickup-db:pickup-postgres -d pickup-go-server"
