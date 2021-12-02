package net

import (
	"net"
	"time"
)

type RichConn struct {
	net.Conn
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LastDataTime int64
}

// SetReadTimeout reset read timeout duration
func (conn *RichConn) SetReadTimeout(t time.Duration) {
	conn.ReadTimeout = t
}

// SetWriteTimeOut reset write timeout duration
func (conn *RichConn) SetWriteTimeOut(t time.Duration) {
	conn.WriteTimeout = t
}

// Read refresh read time
func (conn *RichConn) Read(b []byte) (n int, err error) {
	if conn.ReadTimeout > 0 {
		_ = conn.Conn.SetReadDeadline(time.Now().Add(conn.ReadTimeout))
	} else {
		var t time.Time
		_ = conn.Conn.SetReadDeadline(t)
	}
	return conn.Conn.Read(b)
}

// Write refresh write time
func (conn *RichConn) Write(b []byte) (n int, err error) {
	if conn.WriteTimeout > 0 {
		_ = conn.Conn.SetWriteDeadline(time.Now().Add(conn.WriteTimeout))
	} else {
		var t time.Time
		_ = conn.Conn.SetWriteDeadline(t)
	}
	return conn.Conn.Write(b)
}
