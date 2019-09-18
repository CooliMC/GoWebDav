package gowebdav

import (
	"fmt"
	"log"
	"net"
)

func Start() {
	var flag_address string = ""
	var flag_port string = "443"

	listener, err := net.Listen("tcp", flag_address + ":" + flag_port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on", listener.Addr().String())

	if flag_port != "80" {
		//if err := http.ServeTLS()
	}
}