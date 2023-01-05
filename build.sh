#!/usr/bin/bash

function jumpto
{
    label=$1
    cmd=$(sed -n "/$label:/{:a;n;p;ba};" $0 | grep -v ':$')
    eval "$cmd"
    exit
}

arg=$1

start=${1:-"start"}

jumpto $start


start:
if [[ "$#" -ne 2 ]]; then
    jumpto help
fi
if [[ "$2" == "dev" ]]; then
    jumpto dev
fi
jumpto end

help:
	echo "You can run next commands:"
	echo "  - dev"
	echo "    Makes docker image of the scraping service sutiable for the dev -"
	echo "    it has fast build."
	echo "  - prod"
	echo "  - docker"
	echo "  - compose"
	echo "  - docs"
	echo "  - tests"
	echo "  - go"
    echo "    Create executable of the main go file"
    echo "    (/cmd/scraing_service/main.go)."
	echo "  - run"
jumpto end

dev:
	echo "Deleting previous version of the golang app"
    rm build/main
	echo "Building golang app"
	set GOARCH=amd64
	set GOOS=linux
	go build -o ./build/main ./cmd/scraping_service/main.go
    if [[ "$?" -ne 0 ]]; then
       echo "ERROR with building your golang file!"
       jumpto end
    fi
	set GOARCH=amd64
	set GOOS=windows

	echo "Setting up docker"
	docker build -t sleepygiantpandabear/scraping_service:dev . -f Dockerfile
	echo "Setting up docker-compose"
	docker-compose -f ./docker-compose_dev.yml kill
	docker-compose -f ./docker-compose_dev.yml build 
	docker-compose -f ./docker-compose_dev.yml up -d
jumpto end

end:
