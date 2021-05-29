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

        <div v-if="progress > 0 && progress < 100">
          <b-progress type="is-link" size="is-small" :value=progress format="percent"></b-progress>
        </div>

        <div v-else-if="progress >= 100">
          <b-button type="is-success is-light is-fullwidth" disabled>{{ "Completed" }}</b-button>
        </div>
        <div v-else>

        <div v-if="isReceiver">
          <b-button type="is-warning is-light is-fullwidth" :disabled="!isConnect" icon-left="download" @click="confirmGet">{{ file.name || "File Error" }}</b-button>
        </div>
        <div v-else-if="!pwsConnect">
          <b-button type="is-danger is-light is-fullwidth" disabled icon-left="upload">{{ "Not Connect" }}</b-button>
        </div>
        <div v-else>
          <b-upload v-model="file" @input=onSelect expanded>
            <a class="button is-success is-fullwidth">
              <b-icon icon="upload"></b-icon>
              <span>{{ file.name || "Click to upload"}}</span>
            </a>
          </b-upload>
        </div>

        </div>

      </section>

    </div>
  </div>
</template>

<script>
import { Sender, Recver } from '../lib/transfer'

import Webrtc from '../connect/webrtc'

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
    transfer: {},
    progress: 0,
    cable: {},
    webrtc: {},
    file: {},
    dataChannel: {},
    pwsConnect: false,
    p2pConnect: false,
    isReceiver: false,
    isComplete: false
  }),
  created() {
    let wsUrl = 'ws://localhost:8033/topic/'
    wretch('./config.json').get().json()
      .then(res => {
        this.iceServers = res.iceServers

        wsUrl = res.wsUrl || wsUrl
        this.onConfigServer(wsUrl)
      })
      .catch(err => {
        console.log(err)
        this.onConfigServer(wsUrl)
      })
  },
  computed: {
    isConnect() {
      return this.pwsConnect && this.p2pConnect
    },
    isServer() {
      return !this.$route.params.id
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
      this.webrtc = new Webrtc(this.iceServers, this.cable)
      this.webrtc.onConnected = dataChannel => {
        console.log('P2P onConnected')
        this.dataChannel = dataChannel
        this.onP2PConnect()
      }
    },
    onP2PConnect() {
      this.p2pConnect = true
      if (!this.isServer) {
      }
      if (this.file.name) {
        this.preSend()
      }
    },
    onSelect(file) {
      this.putPeerList()
      if (this.p2pConnect) {
        this.preSend()
      }
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
        this.pwsConnect = true
        this.onPWSConnect()
      }

      cable.onclose = event => {
        this.pwsConnect = false
        console.log('ws close')
      }

      cable.onmessage = event => {
        try {
          const msg = JSON.parse(event.data)
          if (msg.sdp || msg.ice) {
            this.webrtc.onMessage(event.data)
          } else if (msg.topic != null && this.$route.params.id == null) {
            // Get topic name
            this.onTopic(msg.topic)
          } else if (msg.req != null) {
            // Server
            if (this.file.name) {
              this.putPeerList()
            }
            this.webrtc.offer()
          } else if (msg.res != null) {
            // Client
            this.file = msg.res[0]
            this.isReceiver = true
          } else if (msg.checksum != null) {
            this.isComplete = true
            this.transfer.verify(msg.checksum)
          } else {
            console.log(msg)
          }
        } catch (e) {
          console.log(e)
        }
      }
    },
    onTopic(topic) {
      const address = document.location.href + 't/' + topic
      this.address = address
      QRCode.toCanvas(this.$refs.qrcode, address, {
        width: 400
      }, error => {
        if (error) console.error(error)
        console.log('Create QRCode:', address)
      })
    },
    confirmGet() {
      this.preRecv()
    },
    fileList() {
      const list = []
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
    preSend() {
      this.transfer = new Sender(this.file, this.dataChannel)
      this.transfer.onProgress = progress => { this.progress = progress }
      this.transfer.onComplete = checksum => {
        this.checksum = checksum
      }
    },
    preRecv() {
      this.transfer = new Recver(this.file, this.dataChannel)
      this.transfer.onProgress = progress => { this.progress = progress }
      this.transfer.onComplete = checksum => {
        this.dataChannel.send(JSON.stringify({ checksum: checksum }))
      }
      this.transfer.start()
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
