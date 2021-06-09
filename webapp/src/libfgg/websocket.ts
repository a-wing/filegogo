import log from 'loglevel'

const PrefixShare = "/share/"
//const PrefixShort = "/s/"

export default class WebSocketConn {
  server: string
  token: string
  onmessage: ((ev: MessageEvent) => void)

  ws: any

  constructor(addr: string) {
    this.server = addr
    this.token = ""
    this.onmessage = () => {}
  }

  // callback(addr string)
  connect(callback: (addr: string) => void) {
    const ws = new WebSocket(this.authServer())

    this.ws = ws

    // TODO:
    // This browser default
    // Firefox is Blob
    // Chrome, Safari is ArrayBuffer
    ws.binaryType = "arraybuffer"

    ws.onopen = () => { log.debug("websocket connected") }
    ws.onclose = () => { log.debug("websocket disconnected") }
    ws.onerror = () => { log.debug("websocket error") }
    //ws.onmessage = (ev) => {
    //  this.onmessage(ev)
    //}
    ws.onmessage = ({ data }) => {
      try {
        const msg = JSON.parse(data)
        if (msg.share && msg.token) {
          this.updateServer(msg.share)
          this.token = msg.token
          ws.onmessage = this.onmessage
          callback(this.server)
        }
      } catch (e) {
        console.log(e)
      }
    }

  }

  // () => string
  authServer() {
    if (this.token === "") {
      return this.server
    } else {
		  return this.server + "?token=" + this.token
    }
  }

  // share string
  updateServer(share: string) {
    const u = new URL(this.server)
    u.pathname = PrefixShare + share
    this.server = u.toString()
  }

  send(data: any) {
    this.ws.send(data)
  }
}
