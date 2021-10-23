FROM node:lts-alpine as builder-node

WORKDIR /app

COPY webapp .

RUN npm install && npm run build

FROM golang:1.16-buster AS builder

WORKDIR /src

COPY . .

COPY --from=builder-node /app/build /src/server/build

RUN make build

# Bin
FROM scratch AS bin

COPY --from=builder /src/filegogo /usr/bin/filegogo

EXPOSE 8080/tcp

ENTRYPOINT ["/usr/bin/filegogo"]
