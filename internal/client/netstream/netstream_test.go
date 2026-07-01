package netstream

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSenderAndReceiver(t *testing.T) {
	dir := t.TempDir()

	srcPath := filepath.Join(dir, "source.txt")
	content := []byte("hello netstream")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	serverDone := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			serverDone <- err
			return
		}
		defer conn.Close()
		serverDone <- FileReceiver("source.txt", &conn)
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	if err := FileSender(srcPath, &conn); err != nil {
		t.Fatalf("FileSender: %v", err)
	}
	conn.Close()

	if err := <-serverDone; err != nil {
		t.Fatalf("FileReceiver: %v", err)
	}

	got, err := os.ReadFile("source.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove("source.txt")

	if string(got) != string(content) {
		t.Errorf("got %q, want %q", got, content)
	}
}
