# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o silayang .

# Runtime stage
FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Makassar
ENV SILAYANG_ADDR=:8080
ENV SILAYANG_DATA=/app/data/silayang.json
ENV SILAYANG_UPLOADS=/app/uploads

RUN mkdir -p /app/data /app/uploads

COPY --from=builder /app/silayang /app/silayang
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static

EXPOSE 8080

CMD ["/app/silayang"]
