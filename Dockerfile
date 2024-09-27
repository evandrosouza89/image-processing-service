# Build stage
FROM golang:1.22.4-alpine3.20 AS build

RUN apk add --no-cache \
    gcc \
    musl-dev \
    vips-dev \
    upx

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o image-processor-service .
RUN upx --best --lzma image-processor-service

# Final stage
FROM scratch

COPY --from=build /usr/lib /usr/lib
COPY --from=build /lib /lib
COPY --from=build /app/image-processor-service /app/image-processor-service
COPY --from=build /app/static /app/static

WORKDIR /app

EXPOSE 8080

CMD ["./image-processor-service"]