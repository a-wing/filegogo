package webrtc

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"

	"github.com/pion/datachannel"
	"github.com/pion/webrtc/v3"
)

//const messageSize = 15

const (
	TextMessage   = 1
	BinaryMessage = 2
)

type sign struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func SignToData(action string, data []byte) []byte {
	s := &sign{
		Type: action,
		Data: json.RawMessage(data),
	}
	b, _ := json.Marshal(s)
	return b
}

func DataToSign(raw []byte) (string, []byte) {
	s := &sign{}
	json.Unmarshal(raw, s)
	return s.Type, s.Data
}

type Conn struct {
	pc   *webrtc.PeerConnection
	conn datachannel.ReadWriteCloser

	config *webrtc.Configuration

	OnSignSend func([]byte)

	OnOpen    func()
	OnClose   func()
	OnError   func(error)
	OnMessage func([]byte, bool)
}

func NewConn(config *webrtc.Configuration) *Conn {
	return &Conn{
		config:     config,
		OnSignSend: func([]byte) {},

		OnOpen:    func() {},
		OnClose:   func() {},
		OnError:   func(error) {},
		OnMessage: func([]byte, bool) {},
	}
}

func (c *Conn) SignRecv(raw []byte) {
	log.Debug(string(raw))
	sign, data := DataToSign(raw)
	switch sign {
	case "sdp":
		c.RecvSdp(data)
	case "ice":
		candidate := &webrtc.ICECandidateInit{}
		json.Unmarshal(data, candidate)
		c.pc.AddICECandidate(*candidate)
	default:
		log.Warn(string(raw))
	}
}

func (c *Conn) RecvSdp(data []byte) {
	sdp := &webrtc.SessionDescription{}
	json.Unmarshal(data, sdp)
	switch sdp.Type {
	case webrtc.SDPTypeOffer:
		c.getPeerConnection()
		c.recvOffer(sdp)
		c.sendAnswer()
	case webrtc.SDPTypeAnswer:
		c.recvAnswer(sdp)
	case webrtc.SDPTypePranswer:
		// TODO:
		panic("webrtc.SDPTypePranswer")
	case webrtc.SDPTypeRollback:
		// TODO:
		panic("webrtc.SDPTypeRollback")
	default:
		panic(string(data))
	}
}

func (c *Conn) Start() {
	c.getPeerConnection()
	c.sendOffer()
}

func (c *Conn) sendOffer() {
	offer, err := c.pc.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	err = c.pc.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	data, _ := json.Marshal(offer)
	c.OnSignSend(SignToData("sdp", data))
}

func (c *Conn) recvOffer(offer *webrtc.SessionDescription) {
	if err := c.pc.SetRemoteDescription(*offer); err != nil {
		panic(err)
	}
}

func (c *Conn) sendAnswer() {
	answer, err := c.pc.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	err = c.pc.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	data, _ := json.Marshal(answer)
	c.OnSignSend(SignToData("sdp", data))
}

func (c *Conn) recvAnswer(answer *webrtc.SessionDescription) {
	if err := c.pc.SetRemoteDescription(*answer); err != nil {
		panic(err)
	}
}

func (c *Conn) getPeerConnection() *webrtc.PeerConnection {
	s := webrtc.SettingEngine{}
	s.DetachDataChannels()
	api := webrtc.NewAPI(webrtc.WithSettingEngine(s))

	// Create a new RTCPeerConnection using the API object
	peerConnection, err := api.NewPeerConnection(*c.config)
	if err != nil {
		panic(err)
	}
	c.pc = peerConnection

	c.pc.OnICECandidate(func(ice *webrtc.ICECandidate) {
		if ice == nil {
			// TODO
		} else {
			data, _ := json.Marshal(ice.ToJSON())
			//c.OnSignSend(data)
			c.OnSignSend(SignToData("ice", data))
		}
	})

	// Create a datachannel with label 'data'
	negotiated := true
	id := uint16(0)
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

		c.conn = raw
		c.OnOpen()
	})

	return peerConnection
}

func (c *Conn) Send(data []byte, typ bool) error {
	_, err := c.conn.WriteDataChannel(data, typ)
	return err
}

func (c *Conn) Run() {
	for {
		data := make([]byte, 1024*64)
		count, isString, err := c.conn.ReadDataChannel(data)
		if err != nil {
			c.OnError(err)
		} else {
			log.Trace(string(data[:count]), isString)
			c.OnMessage(data[:count], isString)
		}
	}
}

func (w *Conn) Close() error {
	return w.pc.Close()
}
