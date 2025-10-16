package main

import (
	"fmt"
	"net"
	"os"
)

var _ = net.Listen
var _ = os.Exit

func main(){
	//var l net.Listener
	//var err error
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Error Listening", err.Error())
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil{
		fmt.Println("Error Accepting Connection", err.Error())
		os.Exit(1)
	}
conn.Write([]byte("+PONG\r\n"))
}