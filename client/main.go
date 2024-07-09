package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strings"
	"syscall"
	"time"

	"github.com/AidosKuneen/cuckoo"
)

const (
	serverAddress = "127.0.0.1:12345"
)

func findCuckooCycleSolution(ctx context.Context, nonce string) (string, error) {
	solver := cuckoo.NewCuckoo()
	var proof []uint32
	var success bool

	for !success {
		proof, success = solver.PoW([]byte(nonce))
	}

	proofBytes := make([]byte, len(proof)*4)
	for i, val := range proof {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		proofBytes[i*4] = byte(val >> 24)
		proofBytes[i*4+1] = byte(val >> 16)
		proofBytes[i*4+2] = byte(val >> 8)
		proofBytes[i*4+3] = byte(val)
	}

	return hex.EncodeToString(proofBytes), nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	fmem, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer fmem.Close()

	runtime.GC()

	var (
		attempt = 1
		startAt = time.Now()
	)
	for {
		slog.Info("connecting ask challenge")
		conn, err := net.Dial("tcp", serverAddress)
		if err != nil {
			log.Fatalf("error connecting to server: %s", err)
		}
		slog.Info("done connecting ask challenge")

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalf("error reading from server: %s", err)
		}

		solution, err := startConnectionChallenge(ctx, buf, n)
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			attempt++
			conn.Close()
			slog.Warn("timeout solve challenge. try again new challenge task.")
			continue
		case err != nil:
			panic(err)
		default:
			if err := pprof.WriteHeapProfile(fmem); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}

			slog.Info("done", "solution", solution, "attempts", attempt, "duration", time.Since(startAt))
			conn.Write([]byte(solution))

			n, err = conn.Read(buf)
			if err != nil {
				log.Fatalf("error reading from server: %s", err)
			}
			conn.Close()

			slog.Info(string(buf[:n]))
			os.Exit(0)
		}
	}
}

func startConnectionChallenge(ctx context.Context, buf []byte, n int) (string, error) {
	message := string(buf[:n])
	parts := strings.Split(message, ":")
	algorithm := parts[0]
	nonce := parts[1]

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	var (
		solution string
		err      error
	)
	switch algorithm {
	case "CuckooCycle":
		solution, err = findCuckooCycleSolution(ctx, nonce)
	default:
		return "", fmt.Errorf("unknown algorithm: %s", algorithm)
	}

	return solution, err
}
