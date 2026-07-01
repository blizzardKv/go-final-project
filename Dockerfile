FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/scheduler ./scheduler
COPY web ./web
EXPOSE 7540
ENV TODO_PORT=7540
CMD ["./scheduler"]
