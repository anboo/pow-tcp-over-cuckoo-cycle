package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"

	"golang.org/x/crypto/blake2b"
)

const (
	port = ":12345"
	// можно увеличивать в реальном времени во время ддос атак, например
	difficulty = 5
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
	defer conn.Close()

	challenge, err := generateChallenge()
	if err != nil {
		fmt.Println("Error generating challenge:", err)
		return
	}

	_, err = conn.Write([]byte(fmt.Sprintf("%s:%d", challenge, difficulty)))
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	message := string(buffer[:n])
	parts := strings.Split(message, ":")
	if len(parts) != 3 {
		fmt.Println("Invalid message format")
		return
	}

	receivedChallenge := parts[0]
	nonce, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Println("Invalid nonce")
		return
	}
	hash := parts[2]

	if receivedChallenge == challenge && verifyPoW(challenge, nonce, hash, difficulty) {
		conn.Write([]byte(randomQuota()))
	} else {
		response := "Invalid PoW"
		conn.Write([]byte(response))
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
	h := blake2b.Sum256([]byte(record))
	calculatedHash := hex.EncodeToString(h[:])

	target := strings.Repeat("0", difficulty)
	return calculatedHash[:difficulty] == target && calculatedHash == hash
}
