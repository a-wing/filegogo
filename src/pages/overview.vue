<template>
  <div class="container">
    <b-field label="address">
      <b-input v-model="address"></b-input>
    </b-field>
    <section>
      <div v-if="recv.name">
        {{ recv.name }}
        {{ recv.size }}
        {{ recv.type }}
        <b-button type="is-warning is-light" @click="confirmGet">Confirm Recv</b-button>
      </div>
      <div v-else>
        <b-field class="file is-primary" :class="{'has-name': !!file}">
          <b-upload v-model="file" class="file-label" @input=onSelect >
              <span class="file-cta">
                  <b-icon class="file-icon" icon="upload"></b-icon>
                  <span class="file-label">Click to upload</span>
              </span>
              <span class="file-name" v-if="file">
                  {{ file.name }}
              </span>
          </b-upload>
        </b-field>
      </div>
   </section>

  </div>
</template>

<script>
import streamSaver from 'streamsaver'
import SparkMD5 from 'spark-md5'

export default {
  data: () => ({
    address: 'ws://localhost:8033/ws/1234',
    pc: {},
    cable: {},
    file: null,
    send: {},
    recv: {},
    dataChannel: {},
    signChannel: {},
    fileStream: {},
    spark: {},
    checksum: '',
    pointer: 0,
    step: 1024 * 256,
    isComplete: false
  }),
  created() {
    this.connect()

    this.spark = new SparkMD5.ArrayBuffer()
  },
  computed: {
    isServer() {
      return this.$route.params.id ? false : true
    }
  },
  methods: {
    onPWSConnect() {
      if (!this.isServer) {
        this.getPeerList()
      }
    },
    onP2PConnect() {
      if (!this.isServer) {
      }
    },
    onSelect(file) {
      this.putPeerList()
    },
    connect() {
      const cable = new WebSocket(this.address)
      this.cable = cable

      cable.onopen = event => {
        console.log('ws open')
        this.onPWSConnect()
      }

      cable.onclose = event => {
        console.log('ws close')
      }

      cable.onmessage = event => {
        try {
          const msg = JSON.parse(event.data)
          if (msg.sdp != null) {
            console.log('Recv:', msg.type)
            if (msg.type === 'offer') {
              this.answer(msg)
            } else {
              this.onAnswer(msg)
            }
          } else if (msg.ice != null) {
            this.onIncomingICE(msg.ice)
          } else if (msg.req != null) {
            // Server
            if (this.file !== null) {
              this.putPeerList()
            }
            this.offer()
          } else if (msg.res != null) {
            // Client
            this.recv = msg.res[0]
          } else if (msg.checksum != null) {
            this.isComplete = true
            this.checksum = msg.checksum
            console.log(this.checksum)
            this.next()
          } else {
            console.log(msg)
          }
        } catch (e) {
          console.log(e)
        }
      }
    },
    init() {
      const configuration = {
        iceServers: [
          { urls: 'stun:stun.l.google.com:19302' }
        ]
      }

      const pc = new RTCPeerConnection(configuration)
      this.pc = pc

      pc.addEventListener('iceconnectionstatechange', () => {
        console.log('iceconnectionstatechange', pc.iceConnectionState)
      })
      pc.addEventListener('icecandidate', ev => {
        if (ev.candidate === null) {
          console.log(pc)
        } else {
          this.cable.send(JSON.stringify({ ice: ev.candidate }))
        }
      })
    },
    offer() {
      this.init()
      const pc = this.pc

      this.dataChannel = pc.createDataChannel('dataChannel', { reliable: true })
      this.signChannel = pc.createDataChannel('signChannel', { reliable: true })

      this.dataChannel.onopen = () => {
        console.log('data channel open')
      }
      this.dataChannel.onclose = () => {
        console.log('data channel close')
      }

      this.signChannel.onmessage = ev => {
        if (ev.target.label === 'signChannel') {
          this.sendBlob()
        }
      }

      pc.createOffer().then(offer => {
        console.log('on Create offer')
        pc.setLocalDescription(offer)
        this.cable.send(JSON.stringify(offer))
      })

    },
    onAnswer(sdp) {
      this.pc.setRemoteDescription(sdp)
    },
    answer(sdp) {
      this.init()
      const pc = this.pc

      pc.ondatachannel = event => {
        console.log(event)

        if (event.channel.label === 'signChannel') {
          this.signChannel = event.channel
          this.signChannel.onopen = () => {
            console.log('data channel open')
            this.onP2PConnect()
          }
        } else {
          this.dataChannel = event.channel

          this.dataChannel.onmessage = ev => {

            // Md5
            this.spark.append(ev.data)

            this.write([ev.data])
          }
        }
      }

      this.pc.setRemoteDescription(sdp)

      this.pc.createAnswer().then(answer => {
        pc.setLocalDescription(answer)
        this.cable.send(JSON.stringify(answer))
      })
    },
    write(buf) {
      const blob = new Blob(buf)
      const readableStream = blob.stream()
      console.log(blob)

      const reader = readableStream.getReader()
      const pump = () => reader.read()
        .then(res => res.done
          ? this.next()
          : this.fileStream.write(res.value).then(pump))

      pump()
    },
    next() {
      if (this.isComplete) {
        if (this.spark.end() === this.checksum) {
          console.log('Md5 check success')
        }
        this.onFileComplete()
      } else {
        this.signChannel.send('req')
      }
    },
    confirmGet() {
      this.fileStream = streamSaver.createWriteStream(this.recv.name).getWriter()
      this.signChannel.send('req')
    },
    onIncomingICE(ice) {
      const candidate = new RTCIceCandidate(ice)
      console.log(ice)
      this.pc.addIceCandidate(candidate).then(r => {
        console.log(r)
      }).catch(ev => {
        console.log(ev)
      })
    },
    reqData() {
      this.cable.send(JSON.stringify({ event: 'req' }))
    },
    fileList() {
      let list = []
      if (this.file !== null) {
        list.push({
          name: this.file.name,
          size: this.file.size,
          type: this.file.type
        })
      }
      return list
    },
    getPeerList() {
      this.cable.send(JSON.stringify({
        req: this.fileList()
      }))
    },
    putPeerList() {
      this.cable.send(JSON.stringify({
        res: this.fileList()
      }))
    },
    onFileComplete() {
      this.fileStream.close()
    },
    sendBlob() {
      let p = this.pointer

      if (p >= this.file.size) {
        this.checksum = this.spark.end()
        this.cable.send(JSON.stringify({ checksum: this.checksum }))
        console.log(this.checksum)
      }

      this.file.slice(p, p + this.step).arrayBuffer().then(buffer => {

        // Md5
        this.spark.append(buffer)

        this.dataChannel.send(buffer)
      })
      this.pointer = p + this.step
    }
  }
}
</script>
