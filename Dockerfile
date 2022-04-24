# Compile stage
FROM golang:1.18 AS build-env

ADD . /debug_container
WORKDIR /debug_container

RUN go build -o /server

# Final stage
FROM debian:buster

EXPOSE 8000

WORKDIR /
COPY --from=build-env /server /

CMD ["/server"]