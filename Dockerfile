FROM golang:1.15-alpine

# Git is required for go mod to download dependencies.
RUN apk --no-cache add git zip

ARG DIR=/project
WORKDIR $DIR

# Pre-cache dependencies.
ADD go.mod $DIR/
RUN go mod download

ENV CGO_ENABLED=0
ADD . $DIR
