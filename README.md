<h1 align="center">
  <img src="./webapp/public/logo512.png" alt="Filegogo" width="200">
  <br>Filegogo<br>
</h1>

<p align="center">
  <a href="https://send.22333.fun">send.22333.fun</a>
</p>

<h4 align="center">A file transfer tool that can be used in the browser webrtc p2p</h4>

<p align="center">
  <a href="https://github.com/a-wing/filegogo/actions">
    <img src="https://github.com/a-wing/filegogo/workflows/ci/badge.svg" alt="Github Actions">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/a-wing/filegogo">
  <a href="https://goreportcard.com/report/github.com/a-wing/filegogo">
    <img src="https://goreportcard.com/badge/github.com/a-wing/filegogo" alt="Go Report Card">
  </a>
  <a href="https://github.com/a-wing/filegogo/releases">
    <img src="https://img.shields.io/github/release/a-wing/filegogo/all.svg" alt="GitHub Release">
  </a>
  <a href="https://github.com/a-wing/filegogo/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/a-wing/filegogo.svg?maxAge=2592000" alt="License">
  </a>
</p>

[![Demo.gif](https://i.postimg.cc/wTyzyHMc/Peek-2020-10-24-11-29.gif)](https://postimg.cc/8jS992hj)

## Depend

- golang >= 1.16

## Build && Install

```sh
make
```

## Run Development

### Server

```bash
go run ./main.go server
```

### Webapp

```bash
cd webapp

npm install

# frontend
npm run start
```

### Client

> run cli client. For example:

```bash
# send command
go run ./main.go send -s http://localhost:8080/6666 <file>

# recv command
go run ./main.go recv -s http://localhost:8080/6666 <file>
```

## Config

[Reference iceServer config](https://developer.mozilla.org/en-US/docs/Web/API/RTCIceServer)

### Built-in turn server

```toml
# Enable Built-in turn server
[turn]

# if no set, use random user
user = "filegogo:filegogo"

realm = "filegogo"
listen = "0.0.0.0:3478"

# Public ip
# if aws, aliyun
publicIP = "0.0.0.0"
relayMinPort = 49160
relayMaxPort = 49200
```

### iceServer Use Other

For example: [coturn](https://github.com/coturn/coturn)

```sh
apt install coturn
```

```ini
# /etc/turnserver.conf

listening-ip={YOUR_IP_ADDRESS}
relay-ip={YOUR_IP_ADDRESS}

# Public ip
# if aws, aliyun
external-ip={YOUR_IP_ADDRESS}

fingerprint
lt-cred-mech
user=filegogo:filegogo
realm=filegogo

```
