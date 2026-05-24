package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"sync"
	"time"

	convert "github.com/Yeah114/gopherconvert/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
	"github.com/google/uuid"
)

const (
	listenAddress = ":19133"
	serverAddress = "127.0.0.1:19132"
)

func main() {
	pool := minecraft.NewBedrockProtocolPool()
	acceptedProtocols := make([]minecraft.Protocol, len(pool))
	index := 0
	for _, protocol := range pool {
		acceptedProtocols[index] = protocol
		index++
	}
	log.Printf("loaded %d protocols", len(pool))

	cfg := minecraft.ListenConfig{
		ErrorLog:               slog.Default(),
		AuthenticationDisabled: true,
		AllowUnknownPackets:    true,
		AcceptedProtocols:      acceptedProtocols,
		StatusProvider:         minecraft.NewStatusProvider("GopherConvert Proxy", "GopherConvert"),
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

	log.Printf("client %s joined from %s, version: %s", clientConn.IdentityData().DisplayName, clientConn.RemoteAddr(), clientConn.Proto().Ver())
	dialer := minecraft.Dialer{
		ErrorLog:                   slog.Default(),
		IdentityData:               clientConn.IdentityData(),
		ClientData:                 clientConn.ClientData(),
		AutoProtocol:               true,
		DisconnectOnUnknownPackets: false,
		DisconnectOnInvalidPackets: false,
		DownloadResourcePack: func(_ uuid.UUID, _ string, _, _ int) bool {
			return false
		},
	}

	log.Printf("joining server: %s", serverAddress)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	serverConn, err := dialer.DialContext(ctx, "raknet", serverAddress)
	if err != nil {
		log.Printf("dial server failed: %v", err)
		_ = clientConn.WritePacket(&packet.Disconnect{
			Reason:  packet.DisconnectReasonThirdPartyNoInternet,
			Message: err.Error(),
		})
		return
	}
	defer serverConn.Close()
	log.Printf("doing spawn on %s server...", serverConn.Proto().Ver())
	if err := serverConn.DoSpawn(); err != nil {
		log.Printf("spawn on server failed: %v", err)
		return
	}

	converter := convert.NewMinecraftConverter(clientConn, serverConn)
	log.Print("starting client game...")
	gameData := serverConn.GameData()
	if err := converter.StartGameContext(ctx, &gameData); err != nil {
		log.Printf("start client game failed: %v", err)
		_ = clientConn.WritePacket(&packet.Disconnect{
			Reason:  packet.DisconnectReasonThirdPartyNoInternet,
			Message: err.Error(),
		})
	}
	if err != nil {
		log.Printf("create converter failed: %v", err)
		return
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)
	log.Print("start converting")
	go forwardPackets(&wg, errCh, "server to client", serverConn, clientConn, converter)
	go forwardPackets(&wg, errCh, "client to server", clientConn, serverConn, converter)

	err = <-errCh
	_ = clientConn.Close()
	_ = serverConn.Close()
	wg.Wait()
	if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
		log.Printf("session ended: %v", err)
	}
}

func forwardPackets(wg *sync.WaitGroup, errCh chan<- error, name string, src, dst *minecraft.Conn, converter *convert.MinecraftConverter) {
	defer wg.Done()
	for {
		pk, err := src.ReadPacket()
		if err != nil {
			errCh <- fmt.Errorf("%s: read packet: %w", name, err)
			return
		}
		err = converter.HandlePacket(pk, src)
		if err != nil {
			errCh <- fmt.Errorf("%s: write %T: %w", name, pk, err)
			return
		}
	}
}
