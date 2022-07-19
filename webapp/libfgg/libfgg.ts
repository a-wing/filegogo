import log from 'loglevel'

import Pool from './pool/pool'
import { IConn } from './transport/conn'
import { IFile } from "./pool/file/file"
import { Meta, Hash } from "./pool/data"

let uniqueID: number = 0

function getUniqueID(): string {
  return (uniqueID++).toString()
}

const loopWait        = 100
const maxPendingCount = 100

const methodMeta = "meta"
const methodData = "data"
const methodHash = "hash"

type Rpc = {
  [_: string]: (_: any) => Promise<any>
}

type Pending = {
  [_: string]: (_: any) => void
}

export default class Fgg {
  private pool: Pool = new Pool()
  private conn: IConn[] = []

  private finish: boolean = false

  private rpc: Rpc = {
    [methodMeta]: async (_: any): Promise<any> => {
      const meta = await this.pool.sendMeta()
      this.onPreTran(meta)
      return meta
    },
    [methodData]: (data: any): any => { return data },
    [methodHash]: async (_: any): Promise<any> => {
      const hash = this.pool.sendHash()
      this.onPostTran(hash)
      return hash
    },
  }

  private pending: Pending = {}

  private pendingCount: number = 0

  onPreTran: (_: Meta) => void = (_: Meta) => {}
  onPostTran: (_: Hash) => void = (_: Hash) => {}

  addConn(conn: IConn): void {
    conn.setOnRecv((head: ArrayBuffer, body: ArrayBuffer): void => {
      this.recv((new TextDecoder("utf-8").decode(head)), body)
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
      this.finish = true
    }
  }

  run(): Promise<void> {
    return new Promise((resolve) => {
      const timer = setInterval(async () => {
        if (maxPendingCount > this.pendingCount) {
          this.getData()
        }

        if (this.finish) {
          const ok = await this.clientHash()
          if (!ok) {
            log.error("checkSum error")
          }
          log.warn("run finish")

          clearInterval(timer)
          resolve()
        }

      }, loopWait)
    })
  }

  // RPC: Send
  send(head: string, body: ArrayBuffer): void {
    log.debug("SEND", head, body.byteLength)
    this.conn.length > 0 && this.conn[this.conn.length - 1].send((new TextEncoder()).encode(head).buffer, body)
  }

  // RPC: Recv
  async recv(head: string, body: ArrayBuffer): Promise<void> {
    log.debug("RECV", head, body.byteLength)
    const rpc = JSON.parse(head)

    if ("method" in rpc) {
      let res = null
      let err = null
      let body = new ArrayBuffer(0)
      if (rpc.method in this.rpc) {
        try {
          if (rpc.method === methodData) {
            body = await this.pool.sendData(rpc.params)
          }
          res = await this.rpc[rpc.method](rpc.params)
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
        this.send(JSON.stringify(res
          ? {
            jsonrpc: "2.0",
            result: res,
            id: rpc.id,
          }
          : {
            jsonrpc: "2.0",
            error: err,
            id: rpc.id,
          }), body)
      } else {
        // notification
      }

    } else if ("result" in rpc || "error" in rpc) {
      if (rpc.result) {
        if (body.byteLength != 0) {
          this.pendingCount--
          // TODO:
          this.pool.recvData(rpc.result, body)
        } else {
          this.pending[rpc.id](rpc.result)
        }
      } else {
        // TODO: error
        this.pending[rpc.id](rpc.error)
        log.error(rpc.error)
      }
    } else {
      log.warn("Unknown message:", rpc)
    }
  }

  async clientMeta(): Promise<void> {
    const meta = await this.call(methodMeta, null)
    this.onMeta(meta)
  }

  private onMeta(meta: Meta): void {
    this.pool.recvMeta(meta)
    this.onPreTran(meta)
  }

  private call(method: string, params: any): Promise<any> {
    const rpc = {
      jsonrpc: "2.0",
      method: method,
      params: params,
      id: getUniqueID(),
    }

    const head = JSON.stringify(rpc)
    this.send(head, new ArrayBuffer(0))

    return new Promise(resolve => {
      this.pending[rpc.id] = resolve
    })
  }

  private asyncCall(method: string, params: any): void {
    const rpc = {
      jsonrpc: "2.0",
      method: method,
      params: params,
      id: getUniqueID(),
    }

    const head = JSON.stringify(rpc)
    this.send(head, new ArrayBuffer(0))
  }

  private getData() {
    this.pendingCount++
    const c = this.pool.next()
    if (!c) return
    this.asyncCall(methodData, c)
  }

  private async clientHash(): Promise<boolean> {
    const hash = await this.call(methodHash, null)
    this.onPostTran(hash)
    return this.pool.recvHash(hash)
  }
}
