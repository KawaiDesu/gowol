package main

import (
	"bytes"
	"flag"
	"net"
	"os"
)

const ( 
	// This should never appear
	ERR_INTERNAL = 255
	// Something wrong in flags
	ERR_PARAM    = 1
	// Network connection problem
	ERR_NETWORK  = 2
	
)

var (
	port  string
	dhost string
	mac   net.HardwareAddr
	//debug bool
)

func GenMagickPacket([]byte) []byte {
	sync_chain := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
// TODO: convert typcal mac to hex-alike format	
/*
	mac16, err := hex.DecodeString(strings.Repeat(mac1, 16))
	if err != nil {
		println(err.Error())
		os.Exit(ERR_INTERNAL)
	}
*/	
	return append(sync_chain, bytes.Repeat(mac, 16)...)
}

func flags() {
	flag.StringVar(&port, "port", "9", "set custom destination port")
	flag.StringVar(&dhost, "host", "", "set destination host")
	macstring := flag.String("mac", "", "set MAC address")
	/* Not implemented yet
	flag.StringVar(&laddr, "bind", "", "bind to specific network interface")
	flag.BoolVar(&debug, "debug", false, "turn on debug output")
	*/
	flag.Parse()

	if len(dhost) == 0 || len(*macstring) == 0 {
		flag.PrintDefaults()
		os.Exit(ERR_PARAM)
	}

	parsedmac, err := net.ParseMAC(*macstring)
        if err != nil {
                println(err.Error())
                os.Exit(ERR_PARAM)
        }
	mac = parsedmac

}

func main() {
	flags()
	
	target, err := net.Dial("udp", net.JoinHostPort(dhost, port))
	if err != nil {
		println(err.Error())
		os.Exit(ERR_NETWORK)
	}
	
	// IRC-guys said that packet will be sended every Write() call.
	// So while data < MTU (as our magick packet) WOL will be sended in single packet as we need.
	target.Write(GenMagickPacket(mac)) 
	target.Close()
}
