#FROM golang:1.15-alpine AS builder
#
#ARG COMMAND_HANDLER=serverd
#
## Set necessary environmet variables needed for our image
#ENV GO111MODULE=on \
#    CGO_ENABLED=0 \
#    GOOS=linux \
#    GOARCH=amd64
#
## Move to working directory /build
#WORKDIR /build
#
## Copy and download dependency using go mod
#COPY go.mod go.sum ./
#RUN go mod download
#
## Copy the code into the container
#COPY . .
#
## Build the application
#RUN go build -o dist/main cmd/${COMMAND_HANDLER}/main.go


# Build a small image
FROM alpine

ARG COMMAND_HANDLER=serverd

ENV API_HTTP_PORT 8080
EXPOSE $API_HTTP_PORT

ENV MONITORING_PORT 8086
EXPOSE $MONITORING_PORT

COPY ./dist/${COMMAND_HANDLER} ./main
#COPY --from=builder /build/dist/main /
#COPY migrations /migrations


# Command to run
ENTRYPOINT ["./main"]