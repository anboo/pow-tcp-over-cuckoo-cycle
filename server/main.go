package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/scrypt"
)

const (
	port = ":12345"

	initialDifficulty   = 1
	targetCalculateTime = 5 * time.Second
	calculateWindow     = 5

	waitResponseDeadline  = 2 * time.Minute
	sendChallengeDeadline = 3 * time.Second
	sendResultDeadline    = 5 * time.Second
)

var (
	mutex          sync.Mutex
	difficulty     = initialDifficulty
	calculateTimes []time.Duration
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
	mutex.Lock()
	currentDifficulty := difficulty
	mutex.Unlock()
	resp := []byte(fmt.Sprintf("%s:%d", challenge, currentDifficulty))
	_, err = conn.Write(resp)
	if err != nil {
		slog.Debug("send challenge:", err)
		return
	}

	networkStartTime := time.Now()
	err = conn.SetReadDeadline(time.Now().Add(waitResponseDeadline))
	if err != nil {
		slog.Error("error setting read deadline:", err)
		return
	}

	buffer := make([]byte, 110)
	n, err := conn.Read(buffer)
	if err != nil {
		slog.Error("error reading:", err)
		return
	}
	networkDuration := time.Since(networkStartTime)

	receivedMessage := string(buffer[:n])
	parts := strings.Split(receivedMessage, ":")
	if len(parts) != 3 {
		slog.Debug("invalid parts:", receivedMessage)
		return
	}

	receivedChallenge := parts[0]
	nonce, err := strconv.Atoi(parts[1])
	if err != nil {
		slog.Debug("invalid nonce:", receivedMessage)
		return
	}
	hash := parts[2]

	slog.Info("handle conn", "currentDifficulty", currentDifficulty, "nonce", nonce, "size", len(receivedMessage))

	err = conn.SetWriteDeadline(time.Now().Add(sendResultDeadline))
	if err != nil {
		slog.Error("error setting write deadline:", err)
		return
	}
	if receivedChallenge == challenge && verifyPoW(challenge, nonce, hash, currentDifficulty) {
		mutex.Lock()
		calculateTimes = append(calculateTimes, networkDuration)
		if len(calculateTimes) >= calculateWindow {
			difficulty = adjustDifficulty(calculateTimes, targetCalculateTime, currentDifficulty)
			calculateTimes = calculateTimes[1:]
		}
		mutex.Unlock()
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
	N := 1024 * (1 << uint(difficulty)) // Начальное значение 1024, увеличивается экспоненциально
	r := 8
	p := 1
	h, _ := scrypt.Key([]byte(record), []byte(challenge), N, r, p, 32)
	calculatedHash := hex.EncodeToString(h[:])

	return strings.HasPrefix(calculatedHash, "00") && calculatedHash == hash
}

func adjustDifficulty(times []time.Duration, target time.Duration, currentDifficulty int) int {
	totalTime := time.Duration(0)
	for _, t := range times {
		totalTime += t
	}
	averageTime := totalTime / time.Duration(len(times))

	if averageTime < target {
		return currentDifficulty + 1
	} else if averageTime > target {
		return currentDifficulty - 1
	}
	return currentDifficulty
}
