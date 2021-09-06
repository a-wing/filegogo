FROM node:lts-alpine as builder-node

WORKDIR /app

COPY webapp .

RUN npm install && npm run build

FROM golang:1.16-buster AS builder

#ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /src

COPY . .

COPY --from=builder-node /app/build /src/server/build

RUN make build

# Bin
FROM scratch AS bin

COPY --from=builder /src/filegogo /usr/bin/filegogo
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

ENV TZ=Asia/Shanghai

EXPOSE 8033/tcp

ENTRYPOINT ["/usr/bin/filegogo"]
