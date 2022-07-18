interface IConn {
  send(head: ArrayBuffer, body: ArrayBuffer): Promise<void>
  setOnRecv(fn: (head: ArrayBuffer, body: ArrayBuffer) => void): void
}

export type {
  IConn,
}
