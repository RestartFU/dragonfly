package tests

import (
	"io"
	"log/slog"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/df-mc/dragonfly/mocks"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"go.uber.org/mock/gomock"
)

func TestSessionSendMessageWritesTextPacket(t *testing.T) {
	ctrl := gomock.NewController(t)
	conn := mocks.NewMockSessionConn(ctrl)

	var (
		mu      sync.Mutex
		packets []packet.Packet
	)
	conn.EXPECT().ChunkRadius().Return(8).AnyTimes()
	conn.EXPECT().IdentityData().Return(login.IdentityData{
		Identity:    "00000000-0000-0000-0000-000000000001",
		DisplayName: "gomock-test",
	}).AnyTimes()
	conn.EXPECT().RemoteAddr().Return(&net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 19132}).AnyTimes()
	conn.EXPECT().WritePacket(gomock.Any()).DoAndReturn(func(pk packet.Packet) error {
		mu.Lock()
		packets = append(packets, pk)
		mu.Unlock()
		return nil
	}).AnyTimes()
	conn.EXPECT().Close().Return(nil).AnyTimes()

	s := session.Config{
		Log:            slog.New(slog.NewTextHandler(io.Discard, nil)),
		MaxChunkRadius: 8,
	}.New(conn)
	defer s.CloseConnection()

	s.SendMessage("hello from gomock")

	deadline := time.After(time.Second)
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()

	for {
		mu.Lock()
		for _, pk := range packets {
			text, ok := pk.(*packet.Text)
			if ok && text.TextType == packet.TextTypeRaw && text.Message == "hello from gomock" {
				mu.Unlock()
				return
			}
		}
		mu.Unlock()

		select {
		case <-deadline:
			t.Fatal("session did not write expected raw text packet")
		case <-tick.C:
		}
	}
}
