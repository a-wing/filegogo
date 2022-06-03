import log from 'loglevel'

export default class Webrtc {
	//OnOpen    func()
	//OnClose   func()
	//OnError   func(error)
	//OnMessage func([]byte, bool)

	pc:     RTCPeerConnection
	//onSignSend: (msg: string)=>void
	onSignSend: (msg: any)=>void
  //dataChannel: RTCDataChannel | undefined
  dataChannel: RTCDataChannel

  constructor(config: RTCConfiguration) {
    this.pc = new RTCPeerConnection(config)
    this.dataChannel = this.pc.createDataChannel("data", {
      negotiated: true,

      // Chrome needs id < 1024
      // But: ID An 16-bit numeric ID for the channel; permitted values are 0-65534
      // https://developer.mozilla.org/en-US/docs/Web/API/RTCPeerConnection/createDataChannel#rtcdatachannelinit_dictionary
      //id: 1023
      id: 0
    })


    // This browser default
    // Firefox is Blob
    // Chrome, Safari is ArrayBuffer
    this.dataChannel.binaryType = "arraybuffer"

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
          type: 'sdp',
          data: offer,
        })
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
        this.onSignSend({
          type: 'sdp',
          data: answer,
        })
      })
      .catch(e => {
        log.error(e)
      })
  }
  recvAnswer(answer: RTCSessionDescription) {
    this.pc!.setRemoteDescription(answer)
  }

  getPeerConnection() {
    this.pc.addEventListener('iceconnectionstatechange', () => {
      console.log('iceconnectionstatechange', this.pc!.iceConnectionState)
    })
    this.pc.addEventListener('icecandidate', ev => {
      if (ev.candidate) {
        this.onSignSend({
          type: "ice" ,
          data: ev.candidate,
        })
      }
    })
  }

  //send(data: string | ArrayBuffer | Blob) {
  send(data: any) {
    log.info(data)
    this.dataChannel!.send(data)
  }

  close() {
    //TODO:
  }

}
