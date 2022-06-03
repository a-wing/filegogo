import log from 'loglevel'

import WebRTC from './webrtc'
import Transfer from './transfer'

export default class LibFgg {
  ws: WebSocket | null
  rtc: WebRTC | null
  conn: WebSocket | WebRTC | null

  sender: boolean = false

  tran: Transfer

	onPreTran:  (meta: any) => void
	onPostTran: (meta: any) => void

  onRecvFile: () => void
	//OnPreTran:  (meta: Transfer.MetaFile) => void
	//OnPostTran: (meta: Transfer.MetaHash) => void

  constructor() {
    this.onPreTran = () => {}
    this.onPostTran = () => {}
    this.onRecvFile = () => {}

    this.tran = new Transfer()
    this.ws = null
    this.rtc = null
    this.conn = null
  }

  useWebsocket(addr: string) {
    log.debug("websocket connect: ", addr)
    const ws = new WebSocket(addr)

    // This browser default
    // Firefox is Blob
    // Chrome, Safari is ArrayBuffer
    ws.binaryType = "arraybuffer"
    ws.onclose = () => { log.debug("websocket disconnected") }
    ws.onerror = () => { log.debug("websocket error") }

    ws.onmessage = (ev: MessageEvent) => {
      this.recv(ev)
    }
    ws.onopen = () => {
      this.ws = ws
      this.conn = ws
      this.send(JSON.stringify({
        method: "reqlist",
      }))
    }
  }

  close() {
    this.conn?.close()
    this.conn = null
  }

  useWebRTC(config: RTCConfiguration, callback: ()=>void) {
    log.warn("using WebRTC")
    const rtc = new WebRTC(config)
    rtc.onSignSend = (data: any) => {
      this.send(JSON.stringify({
        method: "webrtc",
        params: data,
      }))
    }

    rtc.dataChannel.onmessage = (data: any) => {
      log.debug(data)
      this.recv(data)
    }

    rtc.dataChannel.onopen = () => {
      this.conn = rtc
      log.warn("data channel is open")
      callback()
    }
    this.rtc = rtc
  }

  runWebRTC() {
    this.rtc?.start()
  }

  sendFile(file: File) {
    this.sender = true
    this.tran.send(file)
    this.reslist()
  }

  reslist() {
    if (this.tran.file) {
      this.send(JSON.stringify({
        method: "filelist",
        params: this.tran.getMetaFile()
      }))
    }
  }

  sendData() {
    this.tran.read((buffer: any) => {
      this.send(buffer)
    }, () => {
      log.warn("transfer complete")
    })
  }

  async recv(ev: MessageEvent) {
    const data = ev.data
    if (data instanceof ArrayBuffer) {
      await this.tran.write(data)
      if (this.tran.complete) {
        this.send(JSON.stringify({
          method: "reqsum",
        }))
      } else {
        this.send(JSON.stringify({
          method: "reqdata",
        }))
      }
    } else {
      log.trace(data)
      const rpc = JSON.parse(data)
      switch (rpc.method) {
        case "webrtc":
          log.warn("method 'webrtc'")
          this.rtc?.signRecv(rpc.params)
          break
        case "reqlist":
          this.sender && this.reslist()
          break
        case "getfile":
          this.onPreTran(this.tran.getMetaFile())
          this.sendData()
          break
        case "reqdata":
          this.sendData()
          break
        case "reqsum":
          this.send(JSON.stringify({
            method: "ressum",
            params: this.tran.getMetaHash(),
          }))
          break
        case "ressum":
          if (this.tran.verifyHash(rpc.params)) {
            log.info("md5 verify success")
          } else {
            log.error("md5 verify failure")
          }
          break
        case "filelist":
          log.warn(this)
          this.tran.setMetaFile(rpc.params)
          this.onPreTran(rpc.params)

          this.onRecvFile()

          break
        default:
          log.warn(rpc)
          break
      }
    }
  }

  send(data: string) {
    this.conn?.send(data)
  }

  getfile() {
    log.warn("getfile")
    this.send(JSON.stringify({
      method: "getfile",
    }))
  }
}
