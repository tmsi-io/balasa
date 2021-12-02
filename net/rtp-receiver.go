package net

import (
	"fmt"
	"net"
)

type Receiver struct {
	TCPListener *net.TCPListener
	TCPPort     int
	Stoped      bool
}

var receiver *Receiver = &Receiver{
	Stoped:  true,
	TCPPort: 0,
}

func GetRTPReceiver() *Receiver {
	return receiver
}

func (rec *Receiver) Start() (err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Panic in *Receiver.Start()", err)
		}
	}()
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", rec.TCPPort))
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	rec.Stoped = false
	rec.TCPListener = listener
	networkBuffer := 1024 * 1024 * 1
	for !rec.Stoped {
		conn, err := rec.TCPListener.Accept()
		if err != nil {
			fmt.Printf("RTSP Accept Data Error: %s \n", err)
			continue
		}
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			if err := tcpConn.SetReadBuffer(networkBuffer); err != nil {
				fmt.Printf("RTSP Receiver conn set read buffer error, %v \n", err)
			}
			if err := tcpConn.SetWriteBuffer(networkBuffer); err != nil {
				fmt.Printf("RTSP Receiver conn set write buffer error, %v \n", err)
			}
		}
		NewUploader(conn)
	}
	return
}

func (rec *Receiver) Stop() {
	rec.Stoped = true
	if rec.TCPListener != nil {
		_ = rec.TCPListener.Close()
		rec.TCPListener = nil
	}
}
