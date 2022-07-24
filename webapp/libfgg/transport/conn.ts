interface IConn {
  send(head: ArrayBuffer, body: ArrayBuffer): Promise<void>
  setOnRecv(fn: (head: ArrayBuffer, body: ArrayBuffer) => void): void
}

export type {
  IConn,
}

import { encode, decode } from './protocol'

export default class Conn implements IConn {
  conn: WebSocket | RTCDataChannel
  constructor(conn: WebSocket | RTCDataChannel) {
    this.conn = conn
    this.conn.binaryType = "arraybuffer"
  }
  async send(head: ArrayBuffer, body: ArrayBuffer): Promise<void> {
    this.conn.send(encode(head, body))
  }
  setOnRecv(fn: (head: ArrayBuffer, body: ArrayBuffer) => void): void {
    this.conn.onmessage = (ev: MessageEvent<ArrayBuffer>) => {
      const [head, body] = decode(ev.data)
      fn(head, body)
    }
  }
}
