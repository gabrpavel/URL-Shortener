FROM golang:1.22.3

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o url-shortener ./cmd/url-shortener

ENTRYPOINT ["./url-shortener"]