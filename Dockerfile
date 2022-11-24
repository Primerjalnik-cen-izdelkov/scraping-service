FROM golang:1.16-alpine AS builder

WORKDIR /app

# we should copy only go.mod, install dependicies and then copy the whole app

# COPY . .
COPY go.mod .
COPY go.sum .

RUN go mod edit -go='1.16' -replace='golang.org/x/sys'='golang.org/x/sys@v0.0.0-20220811171246-fbc7d0a398ab'
RUN go mod download

COPY docs/ ./docs/
COPY scrapy_grocery_stores/ ./scrapy_grocery_stores/
COPY pkg/ ./pkg/
COPY internal/ ./internal/
COPY cmd/ ./cmd/
#COPY scrapy_grocery_stores/ ./scrapy_grocery_stores/
RUN go mod vendor

#RUN go build ./cmd/scraping_service/main.go
#CMD ["./main"]
CMD ["go", "run", "./cmd/scraping_service/main.go"]