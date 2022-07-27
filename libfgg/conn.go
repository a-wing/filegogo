package libfgg

import (
	"context"
	"encoding/json"
	"errors"

	dcConn "filegogo/libfgg/transport/webrtc"
	wsConn "filegogo/libfgg/transport/websocket"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	log "github.com/sirupsen/logrus"
)

const (
	methodWebrtcUp  = "webrtc-up"
	methodWebrtcIce = "webrtc-ice"
	methodWebrtcSdp = "webrtc-sdp"
)

func (t *Fgg) UseWebsocket(addr string) error {
	log.Debug("websocket connect: ", addr)

	ws, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}

	conn := wsConn.New(ws)
	t.AddConn(conn)

	go conn.Run(context.Background())
	return nil
}

func (t *Fgg) UseWebRTC(iceServers []webrtc.ICEServer) error {
	setting := webrtc.SettingEngine{}
	setting.DetachDataChannels()
	api := webrtc.NewAPI(webrtc.WithSettingEngine(setting))

	peerConnection, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: iceServers,
	})
	if err != nil {
		return err
	}

	// === Up ===
	t.rpc[methodWebrtcUp] = func(data []byte) (interface{}, error) {
		offer, err := peerConnection.CreateOffer(nil)
		if err != nil {
			return nil, err
		}
		if err := peerConnection.SetLocalDescription(offer); err != nil {
			return nil, err
		}

		req, err := json.Marshal(offer)
		if err != nil {
			return nil, err
		}

		res, _, err := t.call(methodWebrtcSdp, req)
		sdp := webrtc.SessionDescription{}
		if err := json.Unmarshal(res, &sdp); err != nil {
			return nil, err
		}

		if err := peerConnection.SetRemoteDescription(sdp); err != nil {
			return nil, err
		}
		return nil, nil
	}
	// === Up ===

	// === SDP ===
	t.rpc[methodWebrtcSdp] = func(data []byte) (interface{}, error) {
		sdp := webrtc.SessionDescription{}
		if err := json.Unmarshal(data, &sdp); err != nil {
			return nil, err
		}

		switch sdp.Type {
		case webrtc.SDPTypeOffer:
			if err := peerConnection.SetRemoteDescription(sdp); err != nil {
				return nil, err
			}

			answer, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				return answer, err
			}
			err = peerConnection.SetLocalDescription(answer)
			return answer, err
		case webrtc.SDPTypeAnswer:
			// TODO
		case webrtc.SDPTypePranswer:
			// TODO:
			panic("webrtc.SDPTypePranswer")
		case webrtc.SDPTypeRollback:
			// TODO:
			panic("webrtc.SDPTypeRollback")
		default:
			panic(string(data))
		}
		return nil, errors.New("Unknown Error")
	}
	// === SDP ===

	// === ICECandidate ===
	t.rpc[methodWebrtcIce] = func(data []byte) (interface{}, error) {
		log.Tracef("RECV ICE: %s\n", data)
		candidate := webrtc.ICECandidateInit{}
		err := json.Unmarshal(data, &candidate)
		if err != nil {
			log.Error(err)
		}
		peerConnection.AddICECandidate(candidate)
		return nil, err
	}

	peerConnection.OnICECandidate(func(ice *webrtc.ICECandidate) {
		if ice == nil {
			log.Debug("ICE Server already end")
		} else {
			if data, err := json.Marshal(ice.ToJSON()); err != nil {
				log.Error(err)
			} else {
				log.Tracef("SEND ICE: %s\n", data)
				t.notify(methodWebrtcIce, data)
			}
		}
	})
	// === ICECandidate ===

	// === DataChannel ===
	// Create a datachannel with label 'data'
	negotiated := true
	id := uint16(0)
	dataChannel, err := peerConnection.CreateDataChannel("data", &webrtc.DataChannelInit{
		Negotiated: &negotiated,
		ID:         &id,
	})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	var conn *dcConn.Conn
	dataChannel.OnOpen(func() {
		log.Printf("Data channel '%s'-'%d' opened.\n", dataChannel.Label(), dataChannel.ID())

		// Detach the data channel
		dc, err := dataChannel.Detach()
		if err != nil {
			log.Error(err)
		} else {
			conn = dcConn.New(dc)
			t.AddConn(conn)
			go conn.Run(ctx)
		}
	})

	// TODO:
	// This State 'disconnected' no call onClose
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		log.Printf("ICE Connection State has changed: %s\n", connectionState.String())
		if connectionState == webrtc.ICEConnectionStateDisconnected {
			cancel()
			if conn != nil {
				t.DelConn(conn)
				conn = nil
			}
		}
	})

	dataChannel.OnClose(func() {
		log.Printf("Data channel '%s'-'%d' closed.\n", dataChannel.Label(), dataChannel.ID())
		cancel()
		if conn != nil {
			t.DelConn(conn)
			conn = nil
		}
	})
	dataChannel.OnError(func(err error) {
		log.Error(err)
	})
	// === DataChannel ===

	return nil
}
