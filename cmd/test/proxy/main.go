package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	translator "github.com/Yeah114/gopherconvert/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
	"github.com/google/uuid"
)

const (
	listenAddress = ":19133"
	serverAddress = "127.0.0.1:19132"
)

func main() {
	cfg := minecraft.ListenConfig{
		AuthenticationDisabled: true,
		AllowUnknownPackets:    true,
		StatusProvider:         minecraft.NewStatusProvider("GopherTranslate Proxy", "GopherTranslate"),
	}
	listener, err := cfg.Listen("raknet", listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("proxy listening on %s and forwarding to %s", listenAddress, serverAddress)
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Printf("accept failed: %v", err)
			continue
		}
		go handleClient(conn.(*minecraft.Conn))
	}
}

func handleClient(clientConn *minecraft.Conn) {
	defer clientConn.Close()

	log.Printf("client %s joined from %s", clientConn.IdentityData().DisplayName, clientConn.RemoteAddr())
	dialer := minecraft.Dialer{
		IdentityData:               clientConn.IdentityData(),
		ClientData:                 clientConn.ClientData(),
		AutoProtocol:               true,
		DisconnectOnUnknownPackets: false,
		DisconnectOnInvalidPackets: false,
		DownloadResourcePack: func(_ uuid.UUID, _ string, _, _ int) bool {
			return false
		},
	}

	serverConn, err := dialer.Dial("raknet", serverAddress)
	if err != nil {
		log.Printf("dial server failed: %v", err)
		return
	}
	defer serverConn.Close()

	serverToClient, err := translator.NewMinecraftConverter(serverConn, clientConn)
	if err != nil {
		log.Printf("create server-to-client converter failed: %v", err)
		return
	}

	if err := clientConn.StartGame(serverConn.GameData()); err != nil {
		log.Printf("start client game failed: %v", err)
		return
	}
	if err := serverConn.DoSpawn(); err != nil {
		log.Printf("spawn on server failed: %v", err)
		return
	}

	clientToServer, err := translator.NewMinecraftConverter(clientConn, serverConn)
	if err != nil {
		log.Printf("create client-to-server converter failed: %v", err)
		return
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)
	go forwardPackets(&wg, errCh, "server to client", serverConn, clientConn, serverToClient)
	go forwardPackets(&wg, errCh, "client to server", clientConn, serverConn, clientToServer)

	err = <-errCh
	_ = clientConn.Close()
	_ = serverConn.Close()
	wg.Wait()
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
		log.Printf("session ended: %v", err)
	}
}

func forwardPackets(wg *sync.WaitGroup, errCh chan<- error, name string, src, dst *minecraft.Conn, converter *translator.MinecraftConverter) {
	defer wg.Done()
	for {
		if err := src.SetReadDeadline(time.Now().Add(time.Minute)); err != nil {
			errCh <- fmt.Errorf("%s: set read deadline: %w", name, err)
			return
		}
		pk, err := src.ReadPacket()
		if err != nil {
			errCh <- fmt.Errorf("%s: read packet: %w", name, err)
			return
		}
		if err := writePacket(dst, converter, pk); err != nil {
			errCh <- fmt.Errorf("%s: write %T: %w", name, pk, err)
			return
		}
	}
}

func writePacket(dst *minecraft.Conn, converter *translator.MinecraftConverter, pk packet.Packet) error {
	dstPacket, err := converter.ConvertPacket(pk)
	if err != nil {
		return err
	}
	return dst.WritePacket(dstPacket)
}
