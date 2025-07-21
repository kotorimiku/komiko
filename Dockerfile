FROM alpine:latest AS builder

ARG VERSION=dev

WORKDIR /app

RUN apk add --no-cache gcc go nodejs npm

COPY . .

RUN cd web && npm install -g pnpm
RUN cd web && pnpm install
RUN cd web && pnpm run build

RUN go mod download

RUN go build \
    -a -installsuffix cgo \
    -ldflags="-w -s -X komiko/version.Version=${VERSION}" \
    -o komiko .

FROM alpine:latest

ENV VERSION=${VERSION}

WORKDIR /app

COPY --from=builder /app/komiko .

EXPOSE 6060

VOLUME ./data

CMD ["./komiko"]
