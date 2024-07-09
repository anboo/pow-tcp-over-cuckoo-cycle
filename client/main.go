package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/blake2b"
)

const (
	serverAddress = "127.0.0.1:12345"
	difficulty    = 5
)

func main() {
	var (
		startAt  = time.Now()
		solution string
	)
	defer func() {
		slog.Info("done", "solution", solution, "duration", time.Since(startAt))
	}()

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		slog.Error("failed to connect to server", "error", err)
		return
	}
	defer conn.Close()

	data := "Hello, Server!"
	hash, nonce := performPoW(data)
	message := fmt.Sprintf("%s:%d:%s", data, nonce, hash)

	_, err = conn.Write([]byte(message))
	if err != nil {
		slog.Error("failed to write message", "error", err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		slog.Error("failed to read message", "error", err)
		return
	}
	solution = string(buffer[:n])
}

func performPoW(data string) (string, int) {
	var (
		nonce  = 0
		hash   string
		target = strings.Repeat("0", difficulty)
	)

	for {
		nonce++
		record := fmt.Sprintf("%s%d", data, nonce)
		h := blake2b.Sum256([]byte(record))
		hash = hex.EncodeToString(h[:])

		if hash[:difficulty] == target {
			break
		}
	}

	return hash, nonce
}
