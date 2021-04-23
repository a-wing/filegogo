
export default class Webrtc {
  constructor(iceServers, channel) {
    this.channel = channel
    this.iceServers = iceServers
    this.dataChannel = {}
    this.onConnected = () => {}
  }

  onMessage(data) {
    const msg = JSON.parse(data)
    if (msg.sdp) {
      console.log('Recv:', msg.type)
      if (msg.type === 'offer') {
        this.answer(msg)
      } else {
        this.onAnswer(msg)
      }
    } else if (msg.ice) {
      this.onIncomingICE(msg.ice)
    }
  }

  init() {
    this.pc = new RTCPeerConnection({
      iceServers: this.iceServers
    })

    this.pc.addEventListener('iceconnectionstatechange', () => {
      console.log('iceconnectionstatechange', this.pc.iceConnectionState)
    })
    this.pc.addEventListener('icecandidate', ev => {
      if (ev.candidate) {
        this.channel.send(JSON.stringify({ ice: ev.candidate }))
      }
    })
  }

  onIncomingICE(ice) {
    const candidate = new RTCIceCandidate(ice)
    console.log(ice)
    this.pc.addIceCandidate(candidate).then(r => {
      console.log(r)
    }).catch(ev => {
      console.log(ev)
    })
  }

  offer() {
    this.init()
    const pc = this.pc

    // Set dataChannel
    this.dataChannel = pc.createDataChannel('channel')
    this.dataChannel.onopen = ev => this.onConnected(this.dataChannel)

    pc.createOffer().then(offer => {
      console.log('on Create offer')
      pc.setLocalDescription(offer)
      this.channel.send(JSON.stringify(offer))
    })
  }

  onAnswer(sdp) {
    this.pc.setRemoteDescription(sdp)
  }

  answer(sdp) {
    this.init()
    const pc = this.pc

    // Set dataChannel
    pc.ondatachannel = event => {
      this.dataChannel = event.channel
      event.channel.onopen = ev => this.onConnected(event.channel)
    }

    this.pc.setRemoteDescription(sdp)

    this.pc.createAnswer().then(answer => {
      pc.setLocalDescription(answer)
      this.channel.send(JSON.stringify(answer))
    })
  }
}
