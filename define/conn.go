package define

import (
	"context"
	"net"
	"time"

	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/login"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
	"github.com/Yeah114/gophertunnel/minecraft/resource"
)

// Conn is an interface that abstracts the minecraft.Conn type.
type Conn interface {
	Authenticated() bool
	ChunkRadius() int
	ClientCacheEnabled() bool
	ClientData() login.ClientData
	Close() error
	Context() context.Context
	DoSpawn() error
	DoSpawnContext(ctx context.Context) error
	DoSpawnTimeout(timeout time.Duration) error
	Flush() error
	GameData() minecraft.GameData
	IdentityData() login.IdentityData
	Latency() time.Duration
	LocalAddr() net.Addr
	Proto() minecraft.Protocol
	Read(b []byte) (n int, err error)
	ReadBytes() ([]byte, error)
	ReadPacket() (pk packet.Packet, err error)
	RemoteAddr() net.Addr
	ResourcePacks() []*resource.Pack
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	StartGame(data minecraft.GameData) error
	StartGameContext(ctx context.Context, data minecraft.GameData) error
	StartGameTimeout(data minecraft.GameData, timeout time.Duration) error
	Write(b []byte) (n int, err error)
	WritePacket(pk packet.Packet) error
}

var _ Conn = (*minecraft.Conn)(nil)

type EchoConn struct {
	*minecraft.Conn

	packets chan packet.Packet
}

func (c *EchoConn) WritePacket(pk packet.Packet) error {
	c.packets <- pk
	return nil
}

func (c *EchoConn) ReadPacket() (packet.Packet, error) {
	pk := <-c.packets
	return pk, nil
}

func NewEchoConn(conn *minecraft.Conn) Conn {
	return &EchoConn{Conn: conn, packets: make(chan packet.Packet, 10)}
}
