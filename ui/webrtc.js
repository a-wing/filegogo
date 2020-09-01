'use strict';

const vid = document.getElementsByTagName('video')[0];

let cable

let lightcable = (addr) => {
  cable = new WebSocket(addr);

  cable.onopen = event => {
    console.log("ws open")
    cable.send("SESSION_OK");
  }

  cable.onclose = event => {
    console.log("ws close")
  }

  cable.onmessage = event => {
    try {
      let msg = JSON.parse(event.data);

      if (msg.sdp != null) {
        rtcLink(msg.sdp)
      } else if (msg.ice != null) {
        onIncomingICE(msg.ice);
      } else {
        console.log("RECV: EEEEEEEEEEEEEEEEEEEEEE")
        console.log(msg)
      }
    } catch (e) {
      console.log(e)
    }
  }

}

document.getElementById('btn-go').addEventListener('click', async () => {
  console.log("start ws")

  let addr = document.getElementById('address').value
  console.log(addr)
  lightcable(addr)
})

let configuration = {iceServers: [
  {urls: "stun:stun.l.google.com:19302"}
]};

const pc = new RTCPeerConnection(configuration);
//const pc = new RTCPeerConnection();

let onIncomingICE = ice => {
  let candidate = new RTCIceCandidate(ice);
  pc.addIceCandidate(candidate).catch(setError);
}

let setError = msg => {
  console.log("EEEEEEEEEEEEEEEEEEEEEE")
  console.log(msg)
}

let rtcLink = (sdp) => {
  try {
    pc.addTransceiver('video', { direction: 'recvonly' });
    pc.addEventListener('icecandidate', ev => {
      console.log('icecandidate', ev.candidate);
      if (ev.candidate === null) {
        console.log("icecandidate is DONE")
      } else {
        cable.send(JSON.stringify({'ice': ev.candidate }));
      }
    });
    pc.addEventListener('iceconnectionstatechange', () => {
      console.log('iceconnectionstatechange', pc.iceConnectionState);
    });
    pc.addEventListener('signalingstatechange', () => {
      console.log('signalingstatechange', pc.signalingState);
    });
    pc.addEventListener('track', ev => {
      console.log('track', ev.track);
      vid.srcObject = ev.streams[0];
    });
    pc.setRemoteDescription(sdp);
    document.getElementById('remote-sdp').value = sdp.sdp;
    pc.createAnswer().then(answer => {
      pc.setLocalDescription({ type: 'answer', sdp: answer.sdp });
      cable.send(JSON.stringify({'sdp': { type: 'answer', sdp: answer.sdp }}));
      document.getElementById('local-sdp').value = answer.sdp;
    });
  } catch (e) {
    alert(e.name);
  }
}
