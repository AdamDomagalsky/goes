package main

import (
	"io"
	"net"
	"testing"
)

func init() {
	daemonStarted := startNetworkDaemon()
	daemonStarted.Wait()
	daemonStartedWithPool := startNetworkDaemonWithPool()
	daemonStartedWithPool.Wait()
}

// go test -benchtime=10s -bench=.
func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatalf("cannot dial host %v", err)
		}
		if _, err := io.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}

func BenchmarkNetworkRequestWithPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8081")
		if err != nil {
			b.Fatalf("cannot dial host %v", err)
		}
		if _, err := io.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}
