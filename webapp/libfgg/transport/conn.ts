interface IConn {
  send(_: ArrayBuffer, __: ArrayBuffer): Promise<void>
  setOnRecv(_: (_: ArrayBuffer, __: ArrayBuffer) => void): void
}

export type {
  IConn,
}
