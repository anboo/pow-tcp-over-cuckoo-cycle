package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/scrypt"
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
	var hash string

	// Используем параметр N, который увеличивается с ростом сложности
	N := 1024 * (1 << uint(difficulty)) // Начальное значение 1024, увеличивается экспоненциально
	r := 8
	p := 1

	for {
		nonce++
		record := fmt.Sprintf("%s%d", challenge, nonce)
		h, _ := scrypt.Key([]byte(record), []byte(challenge), N, r, p, 32)
		hash = hex.EncodeToString(h[:])

		// Простая проверка: хэш должен начинаться с двух нулей для усложнения задачи
		if strings.HasPrefix(hash, "00") {
			break
		}
	}

	return hash, nonce
}
