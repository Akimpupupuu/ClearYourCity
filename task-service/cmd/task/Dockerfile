FROM golang:1.25-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/task/exe /app/cmd/task

FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/cmd/task/exe /app
CMD ["/app/exe"]