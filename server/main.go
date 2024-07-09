package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"time"

	"github.com/AidosKuneen/cuckoo"
)

const (
	port = ":12345"
)

var algorithms = []string{
	"CuckooCycle",
}

func generateNonce(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seed.Intn(len(charset))]
	}
	return string(result)
}

func checkCuckooCycle(nonce, solution string) bool {
	graph := cuckoo.NewCuckoo()
	proof, err := hex.DecodeString(solution)
	if err != nil {
		return false
	}

	sipkey := []byte(nonce)
	generatedSolution, success := graph.PoW(sipkey)
	if !success {
		return false
	}

	generatedProof := make([]byte, len(generatedSolution)*4)
	for i, val := range generatedSolution {
		generatedProof[i*4] = byte(val >> 24)
		generatedProof[i*4+1] = byte(val >> 16)
		generatedProof[i*4+2] = byte(val >> 8)
		generatedProof[i*4+3] = byte(val)
	}

	compareProofs := func(proof1, proof2 []byte) bool {
		if len(proof1) != len(proof2) {
			return false
		}
		for i := range proof1 {
			if proof1[i] != proof2[i] {
				return false
			}
		}
		return true
	}

	return compareProofs(proof, generatedProof)
}

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

	algorithm := algorithms[rand.Intn(len(algorithms))]
	nonce := generateNonce(16)

	_, err := conn.Write([]byte(fmt.Sprintf("%s:%s", algorithm, nonce)))
	if err != nil {
		slog.Error("conn write", "err", err.Error())
		return
	}

	slog.Info("conn", "nonce", nonce, "alg", algorithm)

	buf := make([]byte, 168)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	solution := string(buf[:n])

	if checkCuckooCycle(nonce, solution) {
		quote := "This is a quote from the book of wisdom."
		conn.Write([]byte(quote))
	} else {
		conn.Write([]byte("Invalid POW solution."))
	}
}
