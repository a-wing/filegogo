# filegogo

A file transfer tool that can be used in the browser webrtc p2p

## Build && Install

```sh
npm install
make
sudo make install
sudo systemctl start filegogo
```

## Run Development

```sh
cp conf/config.json .

npm install

# server
make run

# frontend
npm run dev
```

## Config

```json
{
  "wsUrl": "ws://localhost:8033/topic/",
  "iceServers": [
    {
      "urls": "stun:stun.services.mozilla.com",
      "username": "louis@mozilla.com",
      "credential": "webrtcdemo"
    }, {
      "urls": ["stun:stun.example.com", "stun:stun-1.example.com"]
    }
  ]
}
```

[Reference iceServer config](https://developer.mozilla.org/en-US/docs/Web/API/RTCIceServer)

iceServer Use Other

For example [coturn](https://github.com/coturn/coturn) && [gortcd](https://github.com/gortc/gortcd)

### coturn

```sh
apt install coturn

# Change config
vim /etc/turnserver.conf
```

