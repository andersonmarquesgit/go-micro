FROM golang:alpine AS build

WORKDIR /app

COPY . .

RUN go build -o frontEndApp ./cmd/web

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/frontEndApp /app/frontEndApp

CMD ["/app/frontEndApp"]

EXPOSE 8081