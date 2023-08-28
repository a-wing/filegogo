FROM node:20-alpine as builder-node

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

COPY . .

RUN npm run build

FROM golang:1.20-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY --from=builder-node /app/server/build /app/server/build

RUN make build

# Bin
FROM alpine AS bin

COPY --from=builder /app/conf/filegogo.toml /etc/filegogo.toml
COPY --from=builder /app/filegogo /usr/bin/filegogo

EXPOSE 8080/tcp

CMD ["server"]

ENTRYPOINT ["/usr/bin/filegogo"]
