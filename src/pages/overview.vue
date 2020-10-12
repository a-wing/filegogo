<template>
  <div class="container">
    <b-field label="address">
      <b-input v-model="address"></b-input>
    </b-field>
    <section>
      <b-field>
        <b-upload v-model="dropFiles"
                  multiple
                  drag-drop
                  expanded>
          <section class="section">
            <div class="content has-text-centered">
              <p>
              <b-icon
                icon="upload"
                size="is-large">
              </b-icon>
              </p>
              <p>Drop your files here or click to upload</p>
            </div>
          </section>
        </b-upload>
      </b-field>

      <div class="tags">
        <span v-for="(file, index) in dropFiles"
              :key="index"
              class="tag is-primary" >
              {{file.name}}
              <button class="delete is-small"
                      type="button"
                      @click="deleteDropFile(index)">
              </button>
        </span>
      </div>
    </section>

    <div class="buttons">
      <b-button type="is-warning is-light" @click="getPeerList">getPeerList</b-button>
      <b-button type="is-warning is-light" @click="onFileComplete">onFileComplete</b-button>
    </div>
  </div>
</template>

<script>
import streamSaver from 'streamsaver'

export default {
  data: () => ({
    address: 'ws://localhost:8033/ws/1234',
    pc: {},
    cable: {},
    dataChannel: {},
    signChannel: {},
    dropFiles: [],
    fileStream: {},
    blobs: [],
    buffer: 10,
    pointer: 0,
    step: 1024 * 256,
    isComplete: false
  }),
  created() {
    this.connect()
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
        this.signChannel.send('req')
      }
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
            this.putPeerList()
            this.offer()
          } else if (msg.res != null) {
            // Client
            this.showConfirm(msg.res)
          } else if (msg.close != null) {
            this.isComplete = true
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
            this.blobs.push(ev.data)
            if (this.blobs.length >= this.buffer) {
              this.write(this.blobs)
              this.blobs = []
            }

            this.signChannel.send('req')
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
        if (this.blobs.length != 0) {
          this.write(this.blobs)
          this.blobs = []
        } else {
          this.onFileComplete()
        }
      }
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
    deleteDropFile(index) {
      this.dropFiles.splice(index, 1)
    },
    fileList() {
      let list = []
      if (this.dropFiles.length !== 0) {
        this.dropFiles.forEach(file => {
          list.push({
            name: file.name,
            size: file.size,
            type: file.type
          })
        })
      }
      return list
    },
    getPeerList() {
      this.cable.send(JSON.stringify({
        req: this.fileList()
      }))
      this.fileStream = streamSaver.createWriteStream('filename.txt').getWriter()
    },
    putPeerList() {
      this.cable.send(JSON.stringify({
        res: this.fileList()
      }))
    },
    showConfirm(data) {
      // TODO
    },
    onFileComplete() {
      this.fileStream.close()
    },
    sendBlob() {
      let p = this.pointer

      if (p >= this.dropFiles[0].size) {
        this.cable.send(JSON.stringify({ close: true }))
      }

      this.dropFiles[0].slice(p, p + this.step).arrayBuffer().then(buffer => {
        this.dataChannel.send(buffer)
      })
      this.pointer = p + this.step
    }
  }
}
</script>
