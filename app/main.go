package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

var _ = net.Listen
var _ = os.Exit

func main(){
	//var l net.Listener
	//var err error
	
	l, err := net.Listen("tcp", "0.0.0.0:6379") //This creates a TCP listener that binds to port 6379
	if err != nil {
		fmt.Println("Error Listening", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil{
			fmt.Println("Error Accepting Connection", err.Error())
			os.Exit(1)
		}
	go handleConnection(conn)
	}
	
//conn.Write([]byte("+PONG\r\n"))
}
func handleConnection(conn net.Conn){ //net.Conn returned by l.Accept() as a connection type
		defer conn.Close()
		reader := bufio.NewReader(conn)
		store := make(map[string]string)
	for{
		line, err :=reader.ReadString('\n')
		if err!= nil{
			if err ==io.EOF{
				fmt.Printf("Client Disconnected")
			}else{
				fmt.Printf("Error while reading input")
			}
			return
		}
	//input = strings.TrimSpace(input)
		//if input== "PING"{
	//		conn.Write([]byte("+PONG\r\n"))
	//	} else{
	//		conn.Write([]byte("-ERR unknown command\r\n"))
	//	}
		line = strings.TrimSpace(line)
		var parts []string

		if strings.HasPrefix(line,"*"){ // only resp configuration
			numofele, err := strconv.Atoi(line[1:])
			if err!=nil {
				fmt.Printf("Not able to convert to int - the no of elements")
				return
			}
			for i:= 0; i<numofele; i++{
				linelength, err := reader.ReadString('\n')
				if err!= nil {
					fmt.Printf("Not able to read the length of the line")
					return
				}
				linelength = strings.TrimSpace(linelength)
				strlength, err := strconv.Atoi(linelength[1:]) //skips the first $ and goes to 4 for ECHOhg4
				if err!=nil {
					fmt.Printf("Not able to read the string bulk length")
					return
				}
				buf := make([]byte,strlength)  // make a buffer bucket to read only of strlength i.e $4 only 4 of ECHO
				_, err = reader.Read(buf)	// auto reads buf without storing it. now buf = ECHO
				if err!=nil {
					fmt.Printf("Error reading bytes of strlength in to buf")
					return
				}
				reader.Discard(2) // discarding the last 2 \r\n of any word that we read
				parts = append(parts, string(buf)) // now my vector or slice is formed as parts = {"ECHO", "HELLO"}
			}
			cmd :=strings.ToUpper(parts[0]) // convert the ping to PING or Echo to ECHO

			// instead of conn.Write([]byte (text)) we can use fmt.Fprintf that directly 
			switch cmd {
			case "PING":
				conn.Write([]byte("+PONG\r\n")) // alernative - fmt.Fprintf(conn, "+PONG\r\n")
			case "ECHO":
				if len(parts) > 1 {
					echomsg := strings.Join(parts[1:], " ") // automatically loops and join strings inside parts i.e slice using space
					conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(echomsg), echomsg))) // sending back $4\r\nHELLO\r\n
				}else {
					conn.Write([]byte("$0\r\n\r\n"))
				}
			case "SET":
				if len(parts)!=3 {
					conn.Write([]byte("Error Wrong number of arguments for SET Command\r\n"))
				}else{
					key :=parts[1]
					value := parts[2]
					store[key] = value
					conn.Write([]byte("+OK\r\n"))
				}
			case "GET":
				if len(parts)!=2 {
					conn.Write([]byte("Error Wrong number of arguments for GET Command\r\n"))
				}else{
					key := parts[1]
					getvalue, exists := store[key]
					if !exists{
						conn.Write([]byte("$-1\r\n"))
					}else{
					conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n",len(getvalue), getvalue)))
					}
				}
			default:
				fmt.Fprintf(conn, "Error Unknown Command %s\r\n", cmd)  //sprint used to store into a variable instead of printf that prints directly to screen
			}
		}
	}                                                                          	
}