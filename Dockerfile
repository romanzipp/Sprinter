# syntax=docker/dockerfile:1

## Build go application
FROM golang:1.19-buster AS build

WORKDIR /app

COPY . ./
RUN go mod download
RUN go build -o /sprinter

## Build frontend
FROM node AS build-frontend

WORKDIR /app

COPY . ./
RUN yarn
RUN yarn build

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /sprinter /sprinter
COPY --from=build-frontend /dist/app.css /dist/app.css
COPY --from=build-frontend /static /static

EXPOSE 8080

USER nonroot:nonroot

CMD ["/sprinter"]
