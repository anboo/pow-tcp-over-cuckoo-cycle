package main

import (
	"crypto/rand"
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
	port = ":12345"
	// можно увеличивать в реальном времени во время ддос атак, например
	difficulty = 5

	waitResponseDeadline  = 10 * time.Second
	sendChallengeDeadline = 3 * time.Second
	sendResultDeadline    = 5 * time.Second
)

func main() {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error setting up server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			slog.Debug("close connection:", err)
		}
	}()

	challenge, err := generateChallenge()
	if err != nil {
		slog.Error("error generating challenge:", err)
		return
	}

	err = conn.SetWriteDeadline(time.Now().Add(sendChallengeDeadline))
	if err != nil {
		slog.Error("error setting write deadline:", err)
		return
	}
	resp := []byte(fmt.Sprintf("%s:%d", challenge, difficulty))
	_, err = conn.Write(resp)
	if err != nil {
		slog.Debug("send challenge:", err)
		return
	}

	err = conn.SetReadDeadline(time.Now().Add(waitResponseDeadline))
	if err != nil {
		slog.Error("error setting read deadline:", err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		slog.Error("error reading:", err)
		return
	}

	message := string(buffer[:n])
	parts := strings.Split(message, ":")
	if len(parts) != 3 {
		slog.Debug("invalid parts:", message)
		return
	}

	receivedChallenge := parts[0]
	nonce, err := strconv.Atoi(parts[1])
	if err != nil {
		slog.Debug("invalid nonce:", message)
		return
	}
	hash := parts[2]

	err = conn.SetWriteDeadline(time.Now().Add(sendResultDeadline))
	if err != nil {
		slog.Error("error setting write deadline:", err)
		return
	}
	if receivedChallenge == challenge && verifyPoW(challenge, nonce, hash, difficulty) {
		_, err = conn.Write([]byte(randomQuote()))
	} else {
		response := "Invalid PoW"
		_, err = conn.Write([]byte(response))
	}

	if err != nil {
		slog.Error("error send challenge result:", err)
	}
}

func generateChallenge() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func verifyPoW(challenge string, nonce int, hash string, difficulty int) bool {
	record := fmt.Sprintf("%s%d", challenge, nonce)
	h := sha256.Sum256([]byte(record))
	calculatedHash := hex.EncodeToString(h[:])

	target := strings.Repeat("0", difficulty)
	return calculatedHash[:difficulty] == target && calculatedHash == hash
}
