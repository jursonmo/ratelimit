package ratelimit

import (
	"io"
	"net"
)

type conn struct {
	net.Conn
	r io.Reader
	w io.Writer
}

//example:
// func handleConn(conn net.Conn) {
//	var buf []byte = make([]byte, 1024)
//	var rate int64 = 10000
// 	conn = ConnWithRateLimit(conn, rate)
// 	conn.Read(buf)
//	......
// 	conn.Write(buf)
// }

func ConnWithRateLimit(conn net.Conn, rate int64) net.Conn {
	if rate == 0 {
		return conn
	}
	cap := rate / 10 //默认有十分之一的capacity
	rb := NewBucketWithRate(float64(rate), cap)
	wb := NewBucketWithRate(float64(rate), cap)
	return ConnWithRateBucket(conn, rb, wb)
}

func ConnWithRateBucket(c net.Conn, rb, wb *Bucket) net.Conn {
	conn := &conn{Conn: c}
	if rb != nil {
		conn.r = Reader(c, rb)
	} else {
		conn.r = c
	}
	if wb != nil {
		conn.w = Writer(c, wb)
	} else {
		conn.w = c
	}
	return conn
}

func (c *conn) Read(buf []byte) (int, error) {
	return c.r.Read(buf)
}

func (c *conn) Write(buf []byte) (int, error) {
	return c.w.Write(buf)
}
