package utils

import "net"

func ParseReq(conn net.Conn) ([]byte, error) {
	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
