package hysteria

import (
	stdnet "net"
	"testing"

	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/transport/internet/stat"
)

type testUserConn struct {
	stdnet.Conn
	user *protocol.MemoryUser
}

func (c *testUserConn) User() *protocol.MemoryUser {
	return c.user
}

func TestMemoryUserFromConnectionUnwrapsStatsConnection(t *testing.T) {
	client, server := stdnet.Pipe()
	defer client.Close()
	defer server.Close()

	want := &protocol.MemoryUser{Email: "hysteria@example.com", Level: 1}
	wrapped := &stat.CounterConnection{Connection: &testUserConn{Conn: server, user: want}}

	if got := memoryUserFromConnection(wrapped); got != want {
		t.Fatalf("memoryUserFromConnection() = %p, want %p", got, want)
	}
}

func TestRemoteIPStringSplitsHostPort(t *testing.T) {
	addr := &stdnet.TCPAddr{IP: stdnet.ParseIP("203.0.113.7"), Port: 443}
	if got := remoteIPString(addr); got != "203.0.113.7" {
		t.Fatalf("remoteIPString() = %q, want %q", got, "203.0.113.7")
	}
}

func TestMemoryUserFromConnectionReturnsNilWithoutUser(t *testing.T) {
	client, server := stdnet.Pipe()
	defer client.Close()
	defer server.Close()

	if got := memoryUserFromConnection(server); got != nil {
		t.Fatalf("memoryUserFromConnection() = %p, want nil", got)
	}
}
