package webrtc

import (
	"context"
	"testing"
	"time"

	"github.com/pion/datachannel"
	"github.com/pion/logging"
	"github.com/pion/sctp"
	"github.com/pion/transport/test"

	log "github.com/sirupsen/logrus"
)

// https://github.com/pion/datachannel/blob/1a83011a96d548655174f44c092a8161dea2b18b/datachannel_test.go#L18-L28
func bridgeProcessAtLeastOne(br *test.Bridge) {
	nSum := 0
	for {
		time.Sleep(10 * time.Millisecond)
		n := br.Tick()
		nSum += n
		if br.Len(0) == 0 && br.Len(1) == 0 && nSum > 0 {
			break
		}
	}
}

// https://github.com/pion/datachannel/blob/1a83011a96d548655174f44c092a8161dea2b18b/datachannel_test.go#L30-L81
func createNewAssociationPair(br *test.Bridge) (*sctp.Association, *sctp.Association, error) {
	var a0, a1 *sctp.Association
	var err0, err1 error
	loggerFactory := logging.NewDefaultLoggerFactory()

	handshake0Ch := make(chan bool)
	handshake1Ch := make(chan bool)

	go func() {
		a0, err0 = sctp.Client(sctp.Config{
			NetConn:       br.GetConn0(),
			LoggerFactory: loggerFactory,
		})
		handshake0Ch <- true
	}()
	go func() {
		a1, err1 = sctp.Client(sctp.Config{
			NetConn:       br.GetConn1(),
			LoggerFactory: loggerFactory,
		})
		handshake1Ch <- true
	}()

	a0handshakeDone := false
	a1handshakeDone := false
loop1:
	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		br.Tick()

		select {
		case a0handshakeDone = <-handshake0Ch:
			if a1handshakeDone {
				break loop1
			}
		case a1handshakeDone = <-handshake1Ch:
			if a0handshakeDone {
				break loop1
			}
		default:
		}
	}

	if err0 != nil {
		return nil, nil, err0
	}
	if err1 != nil {
		return nil, nil, err1
	}

	return a0, a1, nil
}

// https://github.com/pion/datachannel/blob/1a83011a96d548655174f44c092a8161dea2b18b/datachannel_test.go#L83-L117
func closeAssociationPair(br *test.Bridge, a0, a1 *sctp.Association) {
	close0Ch := make(chan bool)
	close1Ch := make(chan bool)

	go func() {
		//nolint:errcheck,gosec
		a0.Close()
		close0Ch <- true
	}()
	go func() {
		//nolint:errcheck,gosec
		a1.Close()
		close1Ch <- true
	}()

	a0closed := false
	a1closed := false
loop1:
	for i := 0; i < 100; i++ {
		time.Sleep(10 * time.Millisecond)
		br.Tick()

		select {
		case a0closed = <-close0Ch:
			if a1closed {
				break loop1
			}
		case a1closed = <-close1Ch:
			if a0closed {
				break loop1
			}
		default:
		}
	}
}

func TestWebrtc(t *testing.T) {
	// Limit runtime in case of deadlocks
	lim := test.TimeOut(time.Second * 10)
	defer lim.Stop()

	br := test.NewBridge()
	loggerFactory := logging.NewDefaultLoggerFactory()

	a0, a1, err := createNewAssociationPair(br)
	if err != nil {
		t.Error(err)
	}

	cfg := &datachannel.Config{
		// https://pkg.go.dev/github.com/pion/datachannel#ChannelType
		ChannelType:          datachannel.ChannelTypeReliable,
		ReliabilityParameter: 0,
		Label:                "data",
		LoggerFactory:        loggerFactory,
	}

	dc0, err := datachannel.Dial(a0, 100, cfg)
	if err != nil {
		t.Error(err)
	}
	bridgeProcessAtLeastOne(br)

	dc1, err := datachannel.Accept(a1, &datachannel.Config{
		LoggerFactory: loggerFactory,
	})
	if err != nil {
		t.Error(err)
	}
	bridgeProcessAtLeastOne(br)

	// Application Test
	log.SetLevel(log.TraceLevel)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	head, body := []byte("hello"), []byte("world world !!")

	cc := New(dc0)
	cs := New(dc1)

	sign := make(chan bool)
	cs.SetOnRecv(func(h, b []byte) {
		if string(h) != string(head) {
			t.Error(h)
		}
		if string(b) != string(body) {
			t.Error(b)
		}
		sign <- true
	})

	go cc.Run(ctx)
	go cs.Run(ctx)

	if err := cc.Send(head, body); err != nil {
		t.Error(err)
	}
	bridgeProcessAtLeastOne(br)

	dc0.Close()
	dc1.Close()
	bridgeProcessAtLeastOne(br)

	closeAssociationPair(br, a0, a1)

	<-sign
}
