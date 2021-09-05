# filegogo

A file transfer tool that can be used in the browser webrtc p2p

[![Go Report Card](https://goreportcard.com/badge/github.com/a-wing/filegogo)](https://goreportcard.com/report/github.com/a-wing/filegogo)
[![GitHub release](https://img.shields.io/github/tag/a-wing/filegogo.svg?label=release)](https://github.com/a-wing/filegogo/releases)
[![license](https://img.shields.io/github/license/a-wing/filegogo.svg?maxAge=2592000)](https://github.com/a-wing/filegogo/blob/master/LICENSE)

[send.22333.fun](https://send.22333.fun) | [send.cn.22333.fun](https://send.cn.22333.fun)

[![Demo.gif](https://i.postimg.cc/wTyzyHMc/Peek-2020-10-24-11-29.gif)](https://postimg.cc/8jS992hj)

## Depend

- golang >= 1.16

## Build && Install

```sh
npm install
make
sudo make install
sudo systemctl start filegogo
```

## Run Development

```sh
cp conf/filegogo.toml .
# server
make run

cd webapp

npm install

# frontend
npm run start
```

## Config

[Reference iceServer config](https://developer.mozilla.org/en-US/docs/Web/API/RTCIceServer)

iceServer Use Other

For example [coturn](https://github.com/coturn/coturn) Or [gortcd](https://github.com/gortc/gortcd)

### coturn

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
