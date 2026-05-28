package hysteria

import "testing"

func TestUDPMessageSerializeWritesSessionID(t *testing.T) {
	msg := &UDPMessage{SessionID: 0x01020304, PacketID: 0x0506, FragCount: 1, Addr: "example.com:443", Data: []byte("payload")}
	buf := make([]byte, msg.Size())
	if n := msg.Serialize(buf); n != len(buf) {
		t.Fatalf("Serialize() length = %d, want %d", n, len(buf))
	}
	parsed, err := ParseUDPMessage(buf)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.SessionID != msg.SessionID {
		t.Fatalf("SessionID = %#x, want %#x", parsed.SessionID, msg.SessionID)
	}
}
