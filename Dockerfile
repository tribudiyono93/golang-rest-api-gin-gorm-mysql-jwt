FROM golang:1.16.0-buster
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main
EXPOSE 8080
CMD ["./main"]