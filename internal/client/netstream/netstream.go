package netstream

import (
	"io"
	"net"
	"os"
	"path"
)

func FileSender(filename string, conn *net.Conn) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(*conn, file)
	return err
}

func FileReceiver(filename string, conn *net.Conn) error {
	safename := path.Base(filename)

	file, err := os.OpenFile(safename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, *conn)
	return err
}
