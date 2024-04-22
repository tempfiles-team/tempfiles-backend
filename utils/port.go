package utils

import (
	"fmt"
	"net"
	"strconv"
)

func CheckPortAvailable(port string) string {
	// Parse the port number
	p, err := strconv.Atoi(port)
	if err != nil {
		return ""
	}

	// Check if the port is available
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", p))
	if err != nil {
		return ""
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		// If the port is not available, increment the port number and try again
		return CheckPortAvailable(fmt.Sprintf("%d", p+1))
	}
	defer l.Close()

	// If the port is available, return it

	fmt.Println("Port is available:", port)
	return port
}
