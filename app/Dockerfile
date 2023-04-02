FROM golang:1.19.4-alpine3.17 AS base

WORKDIR "/build"

ARG VERSION

COPY ./src .
#-ldflags "-X main.version=${VERSION}"
RUN go build -o ./app ./main.go

FROM alpine:3.17
RUN apk update && apk add tzdata
WORKDIR "/app"
COPY --from=base /build/app ./app
EXPOSE 3001
ENTRYPOINT chmod +x ./app && ./app serve