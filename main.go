package main

import (
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/aliostad/TraceView/tracing"
)

func main() {

	udpPortPtr := flag.Int("uport", 1969, "UDP port")
	httpPortPtr := flag.Int("hport", 11969, "HTTP port")
	hostPtr := flag.String("host", "0.0.0.0", "host")
	flag.Parse()

	fmt.Println("This is the UDP port: ", *udpPortPtr)
	fmt.Println("This is the HTTP port: ", *httpPortPtr)
	fmt.Println("This is host", *hostPtr)

	dispatch := make(chan string, 200)

	go listenUdp(*udpPortPtr, *hostPtr, dispatch)
	go readFrom(dispatch)

	var tt tracing.Trace = tracing.Trace{}
	fmt.Println(tt)
	_, _ = fmt.Scanln() // stop
	close(dispatch)
}

func readFrom(dispatch <-chan string) {
	for dispatchData := range dispatch {
		fmt.Println(dispatchData)
	}
}

func listenUdp(port int, host string, dispatch chan<- string) {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(host),
	})

	handleErrorNot(err)

	defer conn.Close()
	fmt.Printf("server listening %s\n", conn.LocalAddr().String())

	for {
		message := make([]byte, 64*1024)
		len, _, err := conn.ReadFromUDP(message[:])
		handleErrorNot(err)

		data := strings.TrimSpace(string(message[:len]))
		dispatch <- data
	}

}

// for now we just panic
func handleErrorNot(e error) {
	if e != nil {
		panic(e)
	}
}
