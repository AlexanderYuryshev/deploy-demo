FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app .
EXPOSE 8080
CMD [ "./bot" ]