<template>
  <div class="container">
    <b-field label="address">
      <b-input v-model="address"></b-input>
    </b-field>
    <b-field label="message">
      <b-input v-model="message"></b-input>
    </b-field>
    <b-field label="Message">
      <b-input maxlength="200" type="textarea"></b-input>
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
      <b-button type="is-danger is-light" @click="offer">offer</b-button>
      <b-button type="is-warning is-light" @click="sendMsg">Send</b-button>
      <b-button type="is-warning is-light" @click="ff">F</b-button>
    </div>
  </div>
</template>

<script>
export default {
  data: () => ({
    name: 'John Silver',
    address: 'ws://localhost:8033/ws/1234',
    pc: {},
    channel: {},
    receiveChannel: {},
    cable: {},
    message: "",
    dropFiles: [],
  }),
  created() {
    this.connect(() => {

      if (this.$route.params.id) {
        this.offer()
      }
    })

  },
  mounted() {},
  methods: {
    connect(callback) {
      console.log("connect")

      let cable = new WebSocket(this.address)
      this.cable = cable

      cable.onopen = event => {
        console.log("ws open")
        //cable.send("SESSION_OK");
        callback()
      }

      cable.onclose = event => {
        console.log("ws close")
      }

      cable.onmessage = event => {
        try {
          let msg = JSON.parse(event.data);

          if (msg.sdp != null) {
            //rtcLink(msg.sdp)
            console.log(msg.type)
            if (msg.type == "offer") {
              this.answer(msg)
            } else {
              this.onAnswer(msg)
            }
          } else if (msg.ice != null) {
            //onIncomingICE(msg.ice);
            this.onIncomingICE(msg.ice);
          } else {
            console.log("RECV: EEEEEEEEEEEEEEEEEEEEEE")
            console.log(msg)
          }
        } catch (e) {
          console.log(e)
        }
      }
    },
    init() {
      let configuration = {iceServers: [
        {urls: "stun:stun.l.google.com:19302"}
      ]};

      const pc = new RTCPeerConnection(configuration);
      this.pc = pc

      pc.addEventListener('iceconnectionstatechange', () => {
        console.log('iceconnectionstatechange', pc.iceConnectionState);
      });
      pc.addEventListener('icecandidate', ev => {
        console.log('icecandidate', ev.candidate);
        if (ev.candidate === null) {
          console.log(pc)
          console.log("icecandidate is DONE")
        } else {
          //cable.send(JSON.stringify({'ice': ev.candidate }));
          let msg = JSON.stringify({'ice': ev.candidate });

          //this.cable.send(JSON.stringify(ev));
          this.cable.send(JSON.stringify({'ice': ev.candidate }));

          console.log(msg)
        }
      });
    },
    offer() {
      this.init()
      let pc = this.pc

      this.channel = pc.createDataChannel("sendDataChannel", { reliable: true });

      this.channel.onopen = () => {
        console.log("data channel open")
      }
      this.channel.onclose = () => {
        console.log("data channel close")
      }
      pc.createOffer().then(offer => {
        console.log(offer)
        pc.setLocalDescription(offer);
        //pc.setLocalDescription({ type: 'answer', sdp: answer.sdp });
        //cable.send(JSON.stringify({'sdp': { type: 'answer', sdp: answer.sdp }}));
        this.cable.send(JSON.stringify(offer));
        //document.getElementById('local-sdp').value = answer.sdp;
      });


    },
    onAnswer(sdp) {
      this.pc.setRemoteDescription(sdp);
    },
    answer(sdp) {
      this.init()
      let pc = this.pc

      pc.ondatachannel = event => {
        this.receiveChannel = event.channel;
        this.receiveChannel.onmessage = ev => {
          console.log(ev.data)
        }
      }

      this.pc.setRemoteDescription(sdp);

      this.pc.createAnswer().then(answer => {
        pc.setLocalDescription(answer);
        this.cable.send(JSON.stringify(answer));
      });
    },
    onIncomingICE(ice) {
      let candidate = new RTCIceCandidate(ice);
      console.log(ice)
      //this.pc.addIceCandidate(candidate).catch(setError);
      this.pc.addIceCandidate(candidate).then(r => {
        console.log("candidate set success")
        console.log(r)
      }).catch(ev => {
        console.log("candidate set failure")
        console.log(ev)
      });
    },
    sendMsg() {
      console.log(this.message)
      this.channel.send(this.message)
    },
    deleteDropFile(index) {
      this.dropFiles.splice(index, 1)
    },
    ff() {
      if (this.dropFiles.length != 0) {
        //console.log(this.dropFiles[0].stream().getReader())
        //let _this = this
        this.dropFiles[0].text().then(v => {
          console.log(v)
          this.channel.send(v)
        }).catch(
          d => console.log(d)
        )
      }
    },
  }
}
</script>

