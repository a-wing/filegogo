import log from 'loglevel'

import WebRTC from './webrtc'
import Transfer from './transfer'

export default class LibFgg {
  //ws: WebSocket
  ws: any
  rtc: any
  conn: any

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
  }

  useWebsocket(addr: string) {
    log.debug("websocket connect: ", addr)
    this.ws = new WebSocket(addr)

    // This browser default
    // Firefox is Blob
    // Chrome, Safari is ArrayBuffer
    this.ws.binaryType = "arraybuffer"
    this.ws.onopen = () => {
      this.conn = this.ws
      this.send(JSON.stringify({
        method: "reqlist",
      }))
    }
    this.ws.onclose = () => { log.debug("websocket disconnected") }
    this.ws.onerror = () => { log.debug("websocket error") }

    this.ws.onmessage = (ev: MessageEvent) => {
      this.recv(ev)
    }
  }

  useWebRTC(config: RTCConfiguration, callback: ()=>void) {
    this.rtc = new WebRTC(config)
    this.rtc.onSignSend = (data: any) => {
      this.send(JSON.stringify({
        method: "webrtc",
        params: data,
      }))
    }

    this.rtc.dataChannel.onmessage = (data: any) => {
      log.debug(data)
      this.recv(data)
    }

    this.rtc.dataChannel.onopen = () => {
      this.conn = this.rtc
      log.warn("data channel is open")
      callback()
    }

  }

  runWebRTC() {
    this.rtc.start()
  }

  sendFile(file: File) {
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
          this.rtc.signRecv(rpc.params)
          break
        case "reqlist":
          this.reslist()
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
          if (rpc.share && rpc.token) {
            log.warn(this)
            this.ws.updateServer(rpc.share)
            this.ws.token = rpc.token
            //ws.onmessage = this.onmessage
            //callback(this.server)
          }
          break
      }
    }
  }

  send(data: string) {
    this.conn.send(data)
  }

  getfile() {
    this.send(JSON.stringify({
      method: "getfile",
    }))
  }
}
