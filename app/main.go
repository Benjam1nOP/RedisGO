package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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
	for {
		conn, err := l.Accept()
		if err != nil{
			fmt.Println("Error Accepting Connection", err.Error())
			os.Exit(1)
		}
	//defer l.Close()
	go handleConnection(conn)
	}
	
//conn.Write([]byte("+PONG\r\n"))
}
func handleConnection(conn net.Conn){
	
		defer conn.Close()
		reader := bufio.NewReader(conn)
	for{
		input, err :=reader.ReadString('\n')
		if err!= nil{
			fmt.Printf("Error while reading input")
			return
		}
	input = strings.TrimSpace(input)
		if strings.ToUpper(input)== "PING"{
			conn.Write([]byte("+PONG\r\n"))
		}else{
			conn.Write([]byte("UNKNOWN INPUT\r\n"))
		}
	}
	
}