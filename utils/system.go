package utils

import "net"

func GetFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	addr := listener.Addr().(*net.TCPAddr)
	assignedPort := addr.Port
	listener.Close()
	return assignedPort, nil
}
