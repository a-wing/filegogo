import log from 'loglevel'

export default class WebSocketConn {
  private ws: WebSocket | null = null
  onmessage: (ev: MessageEvent) => void = (_: MessageEvent) => {}

  async useWebsocket(addr: string): Promise<void> {
    // Close always if exist
    if (this.ws) {
      await this.close()
    }

    const ws = new WebSocket(addr)
    this.ws = ws

    // This browser default
    // Firefox is Blob
    // Chrome, Safari is ArrayBuffer
    ws.binaryType = "arraybuffer"
    ws.onclose = () => { log.debug("websocket disconnected") }
    ws.onerror = () => { log.debug("websocket error") }

    ws.onmessage = this.onmessage
    const onopen = async (): Promise<void> => {
      return new Promise<void>(resolve => {
        ws.onopen = () => {
          console.log("opened", ws.readyState)
          resolve()
        }
      })
    }
    return onopen()
  }

  async close(): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      if (this.ws) {
        this.ws.onclose = () => {
          this.ws = null
          resolve()
        }
        this.ws.close()
      } else {
        reject()
      }
    })
  }

  send(data: string | Blob | ArrayBuffer | ArrayBufferView) {
    if (this.ws) {
      this.ws.send(data)
    } else {
      log.error(data)
    }
  }
}
