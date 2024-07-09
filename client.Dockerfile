FROM golang:1.22-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY client .
RUN go build -o tcpclient .
CMD ["./tcpclient"]