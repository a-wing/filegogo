import log from 'loglevel'

import Pool from './pool/pool'
import { IConn } from './transport/conn'
import { IFile } from "./pool/file/file"
import { DataChunk, Meta, Hash } from "./pool/data"

let uniqueID: number = 0

function getUniqueID(): string {
  uniqueID =+ 1
  return uniqueID.toString()
}

const loopWait        = 10
const maxPendingCount = 100

const methodMeta = "meta"
const methodData = "data"
const methodHash = "hash"

type Rpc = {
  [_: string]: (_: string) => string
}

export default class Fgg {
  private pool: Pool = new Pool()
  private conn: IConn[] = []

  private rpc: Rpc = {
    [methodMeta]: (data: any): any => { return data },
    [methodData]: (data: any): any => { return data },
    [methodHash]: (data: any): any => { return data },
  }

  private pendingCoun: number = 0

  onPreTran: (_: Meta) => void = (_: Meta) => {}
  onPostTran: (_: Hash) => void = (_: Hash) => {}

  addConn(conn: IConn): void {
    conn.setOnRecv((head: ArrayBuffer, body: ArrayBuffer) => void {
      //log.debug(head, body)
      //TODO
    })
    this.conn.push(conn)
  }

  delConn(conn: IConn): void {
    this.conn = this.conn.filter(c => c !== conn)
  }

  setSend(file: IFile): void {
    this.pool.setSend(file)
  }

  setRecv(file: IFile): void {
    this.pool.setRecv(file)

    this.pool.OnFinish = () => {
      // TODO: OnFinish
    }
  }

  // RPC: Send
  send(head: ArrayBuffer, body: ArrayBuffer): void {
    log.trace((new TextDecoder("utf-8").decode(head)), body.byteLength)
    this.conn.length > 0 || this.conn[this.conn.length - 1].send(head, body)
  }

  // RPC: Recv
  recv(head: ArrayBuffer, body: ArrayBuffer): void {
    log.trace((new TextDecoder("utf-8").decode(head)), body.byteLength)
    const rpc = JSON.parse((new TextDecoder("utf-8").decode(head)))

    if ("method" in rpc) {
      let res = null
      let err = null
      if (rpc.method in this.rpc) {
        try {
          res = this.rpc[rpc.method](rpc.params)
        } catch (error) {
          err = error
        }
      } else {
        err = {
          code: -32601,
          message: "Method not found"
        }
      }

      // request
      if ("id" in rpc) {
        this.send((new TextEncoder()).encode(JSON.stringify(res
          ? {
            jsonrpc: "2.0",
            result: res,
            id: rpc.id,
          }
          : {
            jsonrpc: "2.0",
            error: err,
            id: rpc.id,
          })).buffer, new ArrayBuffer(0))
      } else {
        // notification
      }

    } else if ("result" in rpc || "error" in rpc) {

    } else {
      //TODO
    }
  }
}
