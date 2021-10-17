<h1 align="center">
  <img src="./webapp/public/logo512.png" alt="Filegogo" width="200">
  <br>Filegogo<br>
</h1>

<h4 align="center">A file transfer tool that can be used in the browser webrtc p2p</h4>

[![Build Status](https://github.com/a-wing/filegogo/workflows/ci/badge.svg)](https://github.com/a-wing/filegogo/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/a-wing/filegogo)](https://goreportcard.com/report/github.com/a-wing/filegogo)
[![GitHub release](https://img.shields.io/github/tag/a-wing/filegogo.svg?label=release)](https://github.com/a-wing/filegogo/releases)
[![license](https://img.shields.io/github/license/a-wing/filegogo.svg?maxAge=2592000)](https://github.com/a-wing/filegogo/blob/master/LICENSE)

[send.22333.fun](https://send.22333.fun)

[![Demo.gif](https://i.postimg.cc/wTyzyHMc/Peek-2020-10-24-11-29.gif)](https://postimg.cc/8jS992hj)

## Depend

- golang >= 1.16

## Build && Install

```sh
make
```

## Run Development

```sh
cp conf/filegogo.toml .
cp conf/filegogo-server.toml .

# run server
make run

# run webapp
cd webapp

npm install

# frontend
npm run start
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
