package hysteria

import (
	"bytes"
	"testing"
)

func TestFragUDPMessageRejectsTooSmallMaxSize(t *testing.T) {
	msg := &UDPMessage{PacketID: 1, FragCount: 1, Addr: "example.com:443", Data: []byte("payload")}
	if _, err := FragUDPMessage(msg, msg.HeaderSize()); err == nil {
		t.Fatal("expected error when max size cannot fit payload")
	}
}

func TestDefraggerHandlesInterleavedPackets(t *testing.T) {
	msgA := &UDPMessage{PacketID: 11, FragCount: 1, Addr: "a.example:443", Data: bytes.Repeat([]byte("a"), 90)}
	msgB := &UDPMessage{PacketID: 12, FragCount: 1, Addr: "b.example:443", Data: bytes.Repeat([]byte("b"), 90)}

	fragsA, err := FragUDPMessage(msgA, msgA.HeaderSize()+30)
	if err != nil {
		t.Fatal(err)
	}
	fragsB, err := FragUDPMessage(msgB, msgB.HeaderSize()+30)
	if err != nil {
		t.Fatal(err)
	}
	if len(fragsA) < 2 || len(fragsB) < 2 {
		t.Fatal("expected fragmented messages")
	}

	df := &Defragger{}
	if got := df.Feed(&fragsA[0]); got != nil {
		t.Fatalf("unexpected complete packet: %+v", got)
	}
	if got := df.Feed(&fragsB[0]); got != nil {
		t.Fatalf("unexpected complete packet: %+v", got)
	}

	var gotA, gotB *UDPMessage
	for i := 1; i < len(fragsA); i++ {
		gotA = df.Feed(&fragsA[i])
	}
	for i := 1; i < len(fragsB); i++ {
		gotB = df.Feed(&fragsB[i])
	}

	if gotA == nil || !bytes.Equal(gotA.Data, msgA.Data) || gotA.Addr != msgA.Addr {
		t.Fatalf("packet A was not reassembled correctly: %+v", gotA)
	}
	if gotB == nil || !bytes.Equal(gotB.Data, msgB.Data) || gotB.Addr != msgB.Addr {
		t.Fatalf("packet B was not reassembled correctly: %+v", gotB)
	}
}
