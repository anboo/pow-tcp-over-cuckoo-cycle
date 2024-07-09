FROM golang:1.22-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY server .
RUN go build -o tcpserver .
EXPOSE 12345
CMD ["./tcpserver"]