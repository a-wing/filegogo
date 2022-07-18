import log from 'loglevel'

import { assert, describe, it } from 'vitest'

import VirtualFile from './pool/file/virtual'

import IOVirtual from './transport/virtual'
import { DataChunk, Meta, Hash } from "./pool/data"

import Fgg from './libfgg'
import { IConn } from "./transport/conn"
import { encode, decode } from './transport/protocol'

class Conn implements IConn {
  conn: any
  constructor(conn: any) {
    this.conn = conn
  }
  send(head: ArrayBuffer, body: ArrayBuffer): Promise<void> {
    this.conn.send(encode(head, body))
    return new Promise<void>((fn) => fn())
  }
  setOnRecv(fn: (head: ArrayBuffer, body: ArrayBuffer) => void): void {
    this.conn.onmessage = (data: ArrayBuffer) => {
      const [head, body] = decode(data)
      fn(head, body)
    }
  }
}

describe('io virtual test', async () => {
  const nameSend = "send_xxx"
  const size = 1024*1024

  const fileSend = new VirtualFile(nameSend, size)
  const fileRecv = new VirtualFile("recv_xxx", size)

  const [a, b] = IOVirtual(2)

  const sender = new Fgg()
  const recver = new Fgg()

  sender.setSend(fileSend)
  recver.setRecv(fileRecv)

  sender.addConn(new Conn(a))
  recver.addConn(new Conn(b))

  log.setLevel(log.levels.TRACE)

  recver.onPreTran = (meta: Meta): void => {
    console.log(meta)
  }

  await recver.clientMeta()

  it('a data', () => {
    assert.equal(true, true)
  })

})
