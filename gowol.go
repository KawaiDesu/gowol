package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"net"
	"strings"
	"os"
)

var (
	port  string
	dhost string
	mac   string
	//debug bool
)

const (
	// This should not appear at all while program compiles correctly.
	ERR_INTERNAL = -1
	// Something wrong in flags
	ERR_PARAM    = 1
	// Network connection problem
	ERR_NETWORK  = 2
	
)

func GenMagickPacket(mac string) []byte {
	sync_chain, err := hex.DecodeString("FFFFFFFFFFFF")
	if err != nil {
		println(err.Error())
		os.Exit(ERR_INTERNAL)
	}
	
	mac16, err := hex.DecodeString(strings.Repeat(mac, 16))
	if err != nil {
		println(err.Error())
		os.Exit(ERR_INTERNAL)
	}
	
	// Im so noob that i dont know any over fancy way to concat two []byte
	buf := bytes.NewBuffer(sync_chain)
	// First returned value is the length of mac16 (written bytes) - i don't need it; err is always nil.
	_, _ = buf.Write(mac16)

	return buf.Bytes()
}

func flags() {
	flag.StringVar(&port, "port", "9", "set custom destination port")
	flag.StringVar(&dhost, "host", "", "set destination host")
	flag.StringVar(&mac, "mac", "", "set MAC address")
	/* Not implemented yet
	flag.StringVar(&laddr, "bind", "", "bind to specific network interface")
	flag.BoolVar(&debug, "debug", false, "turn on debug output")
	*/
	flag.Parse()
	if len(dhost) == 0 || len(mac) == 0 {
		flag.PrintDefaults()
		os.Exit(ERR_PARAM)
	}
}

func main() {
	flags()
	
	target, err := net.Dial("udp", net.JoinHostPort(dhost, port))
	if err != nil {
		println(err.Error())
		os.Exit(ERR_NETWORK)
	}
	
	// IRC-guys said that packet will be sended every Write() call.
	// So while data < MTU (as our magick packet) it [data] will be sended in single packet (as we need for WOL).
	target.Write(GenMagickPacket(mac)) 
	target.Close()
}
