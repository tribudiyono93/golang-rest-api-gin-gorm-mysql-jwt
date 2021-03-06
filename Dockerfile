FROM golang:1.16.0-buster
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go get -u github.com/swaggo/swag/cmd/swag
COPY . .
RUN swag init -g server.go
RUN go build -o main
EXPOSE 8080
CMD ["./main"]