<template>
  <div class="main">
    <div class="card">
      <section>
        <div v-if="isReceiver">
          <div class="detail">
            <b-taglist attached>
              <b-tag type="is-info is-light" size="is-large">Filename</b-tag>
              <b-tag type="is-link is-light" size="is-large">{{ file.name }}</b-tag>
            </b-taglist>
            <b-taglist attached>
              <b-tag type="is-warning is-light" size="is-large">Size</b-tag>
              <b-tag type="is-danger is-light" size="is-large">{{ humanFileSize(file.size) }}</b-tag>
            </b-taglist>
            <b-taglist attached>
              <b-tag type="is-success is-light" size="is-large">Type</b-tag>
              <b-tag type="is-light" size="is-large">{{ file.type }}</b-tag>
            </b-taglist>
          </div>
        </div>
        <div v-else>
          <canvas ref="qrcode"></canvas>

          <div class="address" v-if="address !== ''">
            <b-tag type="is-link is-light" class="address-text">{{ address }}</b-tag>
            <b-button size="is-small" type="is-danger" rounded outlined @click="copy2clipboard">Copy</b-button>
          </div>

        </div>

        <b-progress type="is-link" size="is-small" :value=progress format="percent"></b-progress>

        <div v-if="isReceiver">
          <b-button type="is-warning is-light is-fullwidth" icon-left="download" @click="confirmGet">{{ file.name || "File Error" }}</b-button>
        </div>
        <div v-else>
          <b-upload v-model="file" @input=onSelect expanded>
            <a class="button is-success is-fullwidth">
              <b-icon icon="upload"></b-icon>
              <span>{{ file.name || "Click to upload"}}</span>
            </a>
          </b-upload>
        </div>

      </section>

    </div>
  </div>
</template>

<script>
import { WritableStream } from "web-streams-polyfill/ponyfill";
import streamSaver from 'streamsaver'
streamSaver.WritableStream = WritableStream

import SparkMD5 from 'spark-md5'
import QRCode from 'qrcode'
import wretch from 'wretch'
import copy from 'copy-to-clipboard'
import humanFileSize from 'filesize'

export default {
  data: () => ({
    iceServers: [
      { urls: 'stun:stun.l.google.com:19302' }
    ],
    address: '',
    pc: {},
    cable: {},
    file: {},
    dataChannel: {},
    signChannel: {},
    fileStream: {},
    spark: new SparkMD5.ArrayBuffer(),
    checksum: '',
    pointer: 0,
    step: 1024 * 256,
    isReceiver: false,
    isComplete: false
  }),
  created() {
    let wsUrl = 'ws://localhost:8033/topic/'
    wretch("./config.json").get().json()
    .then(res => {
      this.iceServers = res.iceServers
      console.log(this.iceServers)
      this.onConfigServer(wsUrl)
    })
    .catch(err => {
      console.log(err)
      this.onConfigServer(wsUrl)
    })
  },
  computed: {
    progress() {
      return ( this.pointer / this.file.size ) * 100
    },
    isServer() {
      return this.$route.params.id ? false : true
    }
  },
  methods: {
    onConfigServer(wsUrl) {
      this.connect(this.$route.params.id
        ? wsUrl + this.$route.params.id
        : wsUrl
      )
    },
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
    humanFileSize(size) {
      return humanFileSize(size)
    },
    copy2clipboard() {
      copy(this.address)
    },
    connect(address) {
      console.log(address)
      const cable = new WebSocket(address)
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
          } else if (msg.topic != null) {
            // Get topic name
            this.onTopic(msg.topic)
          } else if (msg.ice != null) {
            this.onIncomingICE(msg.ice)
          } else if (msg.req != null) {
            // Server
            if (this.file.name) {
              this.putPeerList()
            }
            this.offer()
          } else if (msg.res != null) {
            // Client
            this.file = msg.res[0]
            this.isReceiver = true
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
    onTopic(topic) {
      let address = document.location.href + 't/' + topic
      this.address = address
      QRCode.toCanvas(this.$refs.qrcode, address, {
        width: 400
      }, error => {
        if (error) console.error(error)
        console.log('Create QRCode:', address);
      })
    },
    init() {
      const configuration = {
        iceServers: this.iceServers
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

            // computed progress
            this.pointer = this.pointer + this.step

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
      this.fileStream = streamSaver.createWriteStream(this.file.name).getWriter()
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

<style>
.main {
  margin: 100px 0px;
  display: flex;
  flex-direction: row;
  justify-content: center;
}

.card {
  flex-direction: column;
}

.detail {
  margin: 50px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.address {
  margin: 20px;
  display: flex;
  justify-content: center;
}

.address-text {
  margin: 3px;
}

</style>
