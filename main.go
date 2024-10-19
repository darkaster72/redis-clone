package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := createConnection()
	if err != nil {
		return
	}

	defer conn.Close() // close connection once finished

	listen(conn)
}

func listen(conn net.Conn) {
	for {
		resp := NewResp(conn)

		value, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(value)

		conn.Write([]byte("+OK\r\n"))
	}
}

func createConnection() (net.Conn, error) {
	listener, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Listening on Port :6379")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return conn, nil
}
