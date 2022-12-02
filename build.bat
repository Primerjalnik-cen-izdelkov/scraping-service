@echo off

set arg=%1

if [%arg%]==[] goto :help
if %arg%==dev goto :dev
if %arg%==go goto :go

goto :help

:help
	echo You can run next commands:
	echo   - dev
	echo     Makes docker image of the scraping service sutiable for the dev -
	echo     it has fast build.
	echo   - prod
	echo   - docker
	echo   - compose
	echo   - docs
	echo   - tests
	echo   - go
    echo     Create executable of the main go file
    echo     (/cmd/scraing_service/main.go).
	echo   - run
goto :end

:dev
	echo Deleting previous version of the golang app
    del build\main
	echo Building golang app
	set GOARCH=amd64
	set GOOS=linux
	go build -o ./build/main ./cmd/scraping_service/main.go
	set GOARCH=amd64
	set GOOS=windows
    IF %ERRORLEVEL% NEQ 0 ( 
       echo ERROR with building your golang file!
       goto :end
    )
	echo Setting up docker
	docker build -t sleepygiantpandabear/scraping_service:latest . -f Dockerfile
	echo Setting up docker-compose
	docker compose -f ./docker-compose_dev.yml kill
	docker compose -f ./docker-compose_dev.yml build 
	docker compose -f ./docker-compose_dev.yml up -d
goto :end

:go
	echo Building golang app
    go build cmd\scraping_service\main.go
goto :end

:end
