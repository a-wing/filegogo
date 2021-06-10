import log from 'loglevel'

export default class Webrtc {
	//conn datachannel.ReadWriteCloser

	//OnOpen    func()
	//OnClose   func()
	//OnError   func(error)
	//OnMessage func([]byte, bool)

  //config: RTCConfiguration
	//pc:     RTCPeerConnection | undefined
	pc:     RTCPeerConnection
	//onSignSend: (msg: string)=>void
	onSignSend: (msg: any)=>void
  //dataChannel: RTCDataChannel | undefined
  dataChannel: RTCDataChannel

  constructor(config: RTCConfiguration) {
    //this.config = config
    //this.pc = undefined
    this.pc = new RTCPeerConnection(config)
    //this.dataChannel = undefined
    //this.dataChannel = this.pc.createDataChannel("data")
    this.dataChannel = this.pc.createDataChannel("data", {
      negotiated: true,
      id: 1234,
    })


    this.onSignSend = () => {}
    //this.onmessage = () => {}
  }


  //signRecv(raw: string) {
  signRecv(raw: any) {
    log.debug(raw)
    //const p = JSON.parse(raw)
    const p = raw
    switch (p.type) {
      case "sdp":
        this.recvSdp(p.data)
        //c.RecvSdp(data)
        break
      case "ice":
        this.pc!.addIceCandidate(p.data).then(r => {
          log.info(r)
        }).catch(ev => {
          log.error(ev)
        })
        break
      default:
        log.warn(raw)
    }
  }

  recvSdp(sdp: any) {
    switch (sdp.type) {
      case "offer":
        this.getPeerConnection()
        this.recvOffer(sdp)
        this.sendAnswer()
        break
      case "answer":
        this.recvAnswer(sdp)
        break
      case "pranswer":
        // TODO:
        log.error("webrtc.SDPTypePranswer")
        break
      case "rollback":
        // TODO:
        log.error("webrtc.SDPTypeRollback")
        break
      default:
        log.error(sdp)
    }
  }

  start() {
    this.getPeerConnection()
    this.sendOffer()
  }

  sendOffer() {
    this.pc!.createOffer()
      .then(offer => {
        this.pc!.setLocalDescription(offer)
        this.onSignSend({
          type: 'offer',
          data: offer,
        })
        //this.onSignSend(JSON.stringify({
        //  type: 'offer',
        //  data: offer,
        //}))
      })
      .catch(e => {
        log.error(e)
      })
  }
  recvOffer(offer: RTCSessionDescription) {
    this.pc!.setRemoteDescription(offer)
  }
  sendAnswer() {
    this.pc!.createAnswer()
      .then(answer => {
        this.pc!.setLocalDescription(answer)
        this.onSignSend(JSON.stringify({
          type: 'answer',
          data: answer,
        }))
      })
      .catch(e => {
        log.error(e)
      })
  }
  recvAnswer(answer: RTCSessionDescription) {
    this.pc!.setRemoteDescription(answer)
  }

  getPeerConnection() {
    //this.pc = new RTCPeerConnection(this.config)
    //const pc = new RTCPeerConnection(this.config)
    //this.pc = pc
    //const pc = this.pc

    this.pc.addEventListener('iceconnectionstatechange', () => {
      console.log('iceconnectionstatechange', this.pc!.iceConnectionState)
    })
    this.pc.addEventListener('icecandidate', ev => {
      if (ev.candidate) {
        //this.channel.send(JSON.stringify({ ice: ev.candidate }))
        this.onSignSend(JSON.stringify({
          type: "ice" ,
          data: ev.candidate,
        }))
      }
    })

    //const dataChannel = pc.createDataChannel("data", {
    //  negotiated: true,
    //  id: 1234,
    //})
    //this.dataChannel = dataChannel
    const dataChannel = this.dataChannel

    dataChannel.onopen = () => {
      log.warn("DataChanne Open:", dataChannel.id, dataChannel.label)
    }

  }

  //send(data: string | ArrayBuffer | Blob) {
  send(data: any) {
    this.dataChannel!.send(data)
  }

}

