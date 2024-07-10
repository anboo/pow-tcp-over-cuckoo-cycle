package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	serverAddress = "127.0.0.1:12345"
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

	buffer := make([]byte, 34)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	msg := string(buffer[:n])
	parts := strings.Split(msg, ":")
	if len(parts) != 2 {
		slog.Error("failed to parse message")
		return
	}

	difficulty, err := strconv.Atoi(parts[1])
	if err != nil {
		slog.Error("failed to parse difficulty", "msg", msg)
		return
	}

	hash, nonce := performPoW(parts[0], difficulty)
	message := fmt.Sprintf("%s:%d:%s", parts[0], nonce, hash)

	_, err = conn.Write([]byte(message))
	if err != nil {
		slog.Error("failed to write message", "error", err)
		return
	}

	buffer = make([]byte, 1024)
	n, err = conn.Read(buffer)
	if err != nil {
		slog.Error("failed to read message", "error", err)
		return
	}
	solution = string(buffer[:n])
}

func performPoW(challenge string, difficulty int) (string, int) {
	nonce := 0
	target := strings.Repeat("0", difficulty)
	var hash string

	for {
		nonce++
		record := fmt.Sprintf("%s%d", challenge, nonce)
		h := sha256.Sum256([]byte(record))
		hash = hex.EncodeToString(h[:])

		if hash[:difficulty] == target {
			break
		}
	}

	return hash, nonce
}
