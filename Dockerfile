FROM golang:alpine as builder

RUN apk add --no-cache make git
WORKDIR /scraper-src
COPY . /scraper-src
RUN go mod download && \
    go build -o scraper . && \
    mv ./scraper /scraper

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /scraper /
ENTRYPOINT ["/scraper"]
