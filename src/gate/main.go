package main

import (
	"fmt"
	"net"
	"reflect"
)

func procductor(writer *net.Conn) {
	fmt.Println(reflect.TypeOf(writer).String())

	// for {
	// 	b := make([]byte, 1000)
	// 	writer.Write(b)
	// }
}

func customer(reader *net.Conn) {
	fmt.Println(reflect.TypeOf(reader).String())
	// for {
	// 	bytes := make([]byte, 2048)
	// 	reader.Read(bytes)
	// }

}

func main() {
	conn, error := net.Dial("tcp", "google.com:https")
	fmt.Println(conn.RemoteAddr(), error)

	writer, reader := net.Pipe()
	fmt.Println(reflect.TypeOf(writer).String(), reflect.TypeOf(reader).String())

	go procductor(&writer)
	go customer(&reader)

	//select {}
}
