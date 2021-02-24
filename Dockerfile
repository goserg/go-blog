#build stage
FROM golang:alpine AS builder
RUN mkdir /app
ADD . /app/
WORKDIR /app/cmd
RUN go build -o main
EXPOSE 8080
CMD ["/app/cmd/main"]
