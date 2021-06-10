import log from 'loglevel'

import WebSocket from './websocket'
import Transfer from './transfer'

export default class LibFgg {
  //ws: WebSocket
  ws: any
  conn: any

  tran: any
  onShare: (addr: string) => void
	onPreTran:  (meta: any) => void
	onPostTran: (meta: any) => void
	//OnPreTran:  (meta: Transfer.MetaFile) => void
	//OnPostTran: (meta: Transfer.MetaHash) => void

  constructor() {
    this.onShare = () => {}
    this.onPreTran = () => {}
    this.onPostTran = () => {}

    this.tran = new Transfer()
  }

  useWebsocket(addr: string) {
    log.debug("websocket connect: ", addr)
    this.ws = new WebSocket(addr)

    // @ts-ignore
    this.ws.onmessage = (ev) => {
      this.recv(ev)
    }
    this.ws.connect((addr: string) => {
      this.onShare(addr)
      this.conn = this.ws

      this.send(JSON.stringify({
        method: "reqlist",
      }))
    })
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

      this.onPreTran(this.tran.getMetaFile())
    }
  }

  sendData() {
    this.tran.read((buffer: any) => {
      this.send(buffer)
    }, () => {
      log.warn("transfer complete")
    })
  }

  recv(ev: MessageEvent) {
    const data = ev.data
    if (data instanceof ArrayBuffer) {
      this.tran.write(data, () => {
        this.send(JSON.stringify({
          method: "reqdata",
        }))
      }, () => {
        this.send(JSON.stringify({
          method: "reqsum",
        }))
      })
    } else {
      log.trace(data)
      const rpc = JSON.parse(data)
      switch (rpc.method) {
        case "webrtc":
          //t.rtc.SignRecv(*rpc.Params)
          break
        case "reqlist":
          this.reslist()
          break
        case "getfile":
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
          this.tran.verifyHash(rpc.params)
          break
        case "filelist":
          log.warn(this)
          this.tran.setMetaFile(rpc.params)

          this.onPreTran(rpc.params)
          this.send(JSON.stringify({
            method: "getfile",
          }))
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
}
