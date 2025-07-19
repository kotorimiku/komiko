FROM alpine:latest AS builder

WORKDIR /app

RUN apk add --no-cache gcc go

COPY . .

RUN go mod download

RUN go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o komiko .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/komiko .

EXPOSE 6060

VOLUME ./data

CMD ["./komiko"]
