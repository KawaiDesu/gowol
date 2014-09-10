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
	debug bool
)

func GenMagickPacket(mac string) []byte {
	sync_chain, err := hex.DecodeString("FFFFFFFFFFFF")
	if err != nil {
		println(err.Error())
	}
	mac16, err := hex.DecodeString(strings.Repeat(mac, 16))
	if err != nil {
		println(err.Error())
	}
	buf := bytes.NewBuffer(sync_chain)
	_, err = buf.Write(mac16)
	if err != nil {
		println(err.Error())
	}
	return buf.Bytes()
}

func flags() {
	flag.StringVar(&port, "port", "9", "set custom destination port")
	flag.StringVar(&dhost, "host", "", "set destination host")
	flag.StringVar(&mac, "mac", "", "set MAC address")
	flag.StringVar(&laddr, "bind", "", "bind to specific network interface")
	flag.BoolVar(&debug, "debug", false, "turn on debug output")
	flag.Parse()
	if len(dhost) == 0 || len(mac) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func atstart() {
	flags()

}

func main() {
	atstart()
	target, err := net.Dial("udp", net.JoinHostPort(dhost, port))
	if err != nil {
		println(err.Error())
	}
	target.Write(GenMagickPacket(mac))
	target.Close()
}
