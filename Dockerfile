# syntax=docker/dockerfile:1

##
## Build Web
##
FROM node:lts AS build-web
WORKDIR /
# Copy our preact app
COPY /web ./
# Install NPM packages including preact-cli
RUN npm install
# Create a production build
RUN npm run build

##
## Build API
##
FROM golang:latest AS build-api
WORKDIR /
# Copy go module configuration
COPY go.mod ./
COPY go.sum ./
# Copy over our API
COPY /api ./
# Download go modules
RUN go mod download
# Build our package 
RUN go build -o /toiler-web

##
## Deploy
##
FROM gcr.io/distroless/base
WORKDIR /
# Copy our binary
COPY --from=build-api /toiler-web ./toiler-web
# Copy our built preact app
COPY --from=build-web /build ./public
# Copy our db migrations
COPY /database/migrations ./database/migrations
# Expose port 8080
# EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/toiler-web"]