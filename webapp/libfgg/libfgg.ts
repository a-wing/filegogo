import log from 'loglevel'

import Pool from './pool/pool'
import Conn from './transport/conn'
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

const methodWebrtcUp  = "webrtc-up"
const methodWebrtcIce = "webrtc-ice"
const methodWebrtcSdp = "webrtc-sdp"

type Rpc = {
  [_: string]: (_: any) => Promise<any>
}

type Pending = {
  [_: string]: {
    resolve: (_: any) => void,
    reject: (_: any) => void,
  }
}

export default class Fgg {
  private pool: Pool = new Pool()
  private conn: IConn[] = []

  private finish: boolean = false

  private rpc: Rpc = {
    [methodMeta]: async (meta: any): Promise<any> => {
      if (meta) {
        this.onMeta(meta)
      } else {
        meta = await this.pool.sendMeta()
        this.onPreTran(meta)
      }
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

  onSendFile: (_: Meta) => void = () => {}
  onRecvFile: (_: Meta) => void = () => {}

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
    const fn = async () => {
      const meta = await this.pool.sendMeta()
      this.onSendFile(meta)
      this.notify(methodMeta, meta)
    }
    fn()
  }

  setRecv(file: IFile): void {
    this.pool.setRecv(file)

    this.pool.OnFinish = () => {
      this.finish = true
    }
  }

  async run(): Promise<void> {
    if (this.rpc[methodWebrtcUp]) {
      await this.rpc[methodWebrtcUp](null)
    }

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
          this.pending[rpc.id]?.resolve(rpc.result)
        }
      } else {
        this.pending[rpc.id]?.reject(new Error(rpc.error))
        log.debug(rpc.error)
      }
    } else {
      log.warn("Unknown message:", rpc)
    }
  }

  async clientMeta(): Promise<void> {
    try {
      const meta = await this.call(methodMeta, null)
      if (meta) {
        this.onMeta(meta)
      }
    } catch (e) {
      // Ignore this error
      log.debug(e)
    }
  }

  private onMeta(meta: Meta): void {
    this.onRecvFile(meta)
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

    return new Promise((resolve, reject) => {
      this.pending[rpc.id] = {
        resolve: resolve,
        reject: reject,
      }
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

  private notify(method: string, params: any): void {
    const rpc = {
      jsonrpc: "2.0",
      method: method,
      params: params,
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
    try {
      const hash = await this.call(methodHash, null)
      this.onPostTran(hash)
      return this.pool.recvHash(hash)
    } catch (err) {
      console.log(err)
    }
    return false
  }

  async useWebsocket(addr: string): Promise<void> {
    log.debug("websocket connect: ", addr)
    const ws = new WebSocket(addr)
    this.addConn(new Conn(ws))
    return new Promise((resolve) => {
      ws.onopen = () => {
        resolve()
      }
    })
  }

  async useWebRTC(cfg: RTCConfiguration): Promise<void> {
    const peerConnection = new RTCPeerConnection(cfg)

    // === Up ===
    this.rpc[methodWebrtcUp] = async (): Promise<any> => {
      const offer = await peerConnection.createOffer()
      await peerConnection.setLocalDescription(offer)
      const sdp = await this.call(methodWebrtcSdp, offer)
      await peerConnection.setRemoteDescription(sdp)
    }
    // === Up ===

    // === SDP ===
    this.rpc[methodWebrtcSdp] = async (sdp): Promise<any> => {
      switch (sdp.type) {
        case "offer":
          await peerConnection.setRemoteDescription(sdp)
          const answer = await peerConnection.createAnswer()
          await peerConnection.setLocalDescription(answer)
          return answer
        case "answer":
          // TODO
          break
        case "pranswer":
          // TODO:
          log.error("webrtc.SDPTypePranswer", sdp)
          break
        case "rollback":
          // TODO:
          log.error("webrtc.SDPTypeRollback", sdp)
          break
        default:
          log.error(sdp)
          break
      }
    }
    // === SDP ===

    // === DataChannel ===
    // Create a datachannel with label 'data'
    const datachannel = peerConnection.createDataChannel("data", {
      negotiated: true,

      // Chrome needs id < 1024
      // But: ID An 16-bit numeric ID for the channel; permitted values are 0-65534
      // https://developer.mozilla.org/en-US/docs/Web/API/RTCPeerConnection/createDataChannel#rtcdatachannelinit_dictionary
      //id: 1023
      id: 0
    })

    let conn: Conn | null = null
    datachannel.onopen = () => {
      conn = new Conn(datachannel)
      this.addConn(conn)
    }

    // This State 'disconnected' no call onClose
    datachannel.onclose = () => {
      log.info(`Data channel '${datachannel.label}'-'${datachannel.id}' closed.`)
      if (conn) {
        this.delConn(conn)
        conn = null
      }
    }

    datachannel.onerror = (err) => {
      log.error(err)
    }
    // === DataChannel ===

    // === ICECandidate ===
    this.rpc[methodWebrtcIce] = async (candidate): Promise<any> => {
      log.trace(`RECV ICE: ${candidate}`)
      peerConnection.addIceCandidate(candidate)
    }

    peerConnection.onicecandidate = (ice) => {
      if (ice.candidate) {
        //const data = JSON.stringify(ice.candidate)
        log.trace(`SEND ICE: ${ice.candidate}`)
        this.notify(methodWebrtcIce, ice.candidate)
      } else {
        log.debug("ICE Server already end")
      }
    }

    peerConnection.oniceconnectionstatechange = () => {
      const connectionState = peerConnection.connectionState
      log.info(`ICE Connection State has changed: ${connectionState}`)
      if (connectionState == "disconnected") {
        if (conn) {
          this.delConn(conn)
          conn = null
        }
      }
    }
    // === ICECandidate ===
  }

  setOnProgress(fn: (c: number) => void) {
    this.pool.OnProgress = fn
  }
}
