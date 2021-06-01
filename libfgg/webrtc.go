package libfgg

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/SB-IM/jsonrpc-lite"
	"github.com/pion/datachannel"
	"github.com/pion/webrtc/v3"
)

const messageSize = 15

type WebrtcConn struct {
	sign chan bool
	conn datachannel.ReadWriteCloser
}

func NewWebrtcConn() *WebrtcConn {
	return &WebrtcConn{
		sign: make(chan bool),
	}
}

func (w *WebrtcConn) Send(t int, data []byte) error {
	isString := true
	if t == BinaryMessage {
		isString = false
	}
	_, err := w.conn.WriteDataChannel(data, isString)
	//log.Println(c, err)
	return err
}

func (w *WebrtcConn) Recv() (int, []byte, error) {
	t := TextMessage
	data := make([]byte, 1024)
	//data := make([]byte, messageSize)
	c, isString, err := w.conn.ReadDataChannel(data)
	if !isString {
		t = BinaryMessage
	}
	return t, data[:c], err
}

func (w *WebrtcConn) getPeerConnection() *webrtc.PeerConnection {
	s := webrtc.SettingEngine{}
	s.DetachDataChannels()

	// Create an API object with the engine
	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))

	// Everything below is the Pion WebRTC API! Thanks for using it ❤️.

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection using the API object
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Create a datachannel with label 'data'
	negotiated := true
	id := uint16(1234)
	dataChannel, err := peerConnection.CreateDataChannel("data", &webrtc.DataChannelInit{
		Negotiated: &negotiated,
		ID:         &id,
	})
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	dataChannel.OnOpen(func() {
		log.Printf("Data channel '%s'-'%d' open.\n", dataChannel.Label(), dataChannel.ID())

		// Detach the data channel
		raw, dErr := dataChannel.Detach()
		if dErr != nil {
			panic(dErr)
		}

		w.conn = raw
		w.sign <- true
	})

	return peerConnection
}

func (w *WebrtcConn) RunOffer(c Conn) {
	peerConnection := w.getPeerConnection()

	// Create an offer to send to the browser
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	data, err := jsonrpc.NewNotify("offer", *peerConnection.LocalDescription()).ToJSON()
	if err != nil {
		panic(err)
	}
	c.Send(TextMessage, data)

	// Wait for the answer
	_, msg, err := c.Recv()
	if err != nil {
		panic(err)
	}
	answer := webrtc.SessionDescription{}
	json.Unmarshal(*jsonrpc.ParseObject(msg).Params, &answer)

	// Apply the answer as the remote description
	err = peerConnection.SetRemoteDescription(answer)
	if err != nil {
		panic(err)
	}
}

func (w *WebrtcConn) RunAnswer(c Conn) {
	peerConnection := w.getPeerConnection()

	// Wait for the offer
	_, msg, err := c.Recv()
	if err != nil {
		panic(err)
	}
	offer := webrtc.SessionDescription{}
	json.Unmarshal(*jsonrpc.ParseObject(msg).Params, &offer)

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	data, err := jsonrpc.NewNotify("answer", *peerConnection.LocalDescription()).ToJSON()
	if err != nil {
		panic(err)
	}

	c.Send(TextMessage, data)
}
