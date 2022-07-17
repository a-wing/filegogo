import { IConn } from "./conn"

import { encode, decode } from './protocol'

export default class Conn implements IConn {
  conn: WebSocket
  constructor(conn: WebSocket) {
    this.conn = conn
    this.conn.binaryType = "arraybuffer"
  }
  send(head: ArrayBuffer, body: ArrayBuffer): Promise<void> {
    this.conn.send(encode(head, body))
    return new Promise<void>((fn) => fn())
  }
  setOnRecv(fn: (_: ArrayBuffer, __: ArrayBuffer) => void): void {
    this.conn.onmessage = (ev: MessageEvent<ArrayBuffer>) => {
      const [head, body] = decode(ev.data)
      fn(head, body)
    }
  }
}
