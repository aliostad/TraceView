package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/aliostad/TraceView/tracing"
)

func main() {

	udpPortPtr := flag.Int("uport", 1969, "UDP port")
	httpPortPtr := flag.Int("hport", 8969, "HTTP port")
	hostPtr := flag.String("host", "0.0.0.0", "host")
	timestampFieldNamesPtr := flag.String("tfn", "", "timestamp field names, comma separated")
	messageFieldNamesPtr := flag.String("mfn", "", "message field names, comma separated")
	levelFieldNamesPtr := flag.String("lfn", "", "levle field names, comma separated")
	corridFieldNamesPtr := flag.String("cfn", "", "correlation Id field names, comma separated")
	indexableFieldNamesPtr := flag.String("ifn", "", "indexable field names, comma separated")
	keepOriginalPayloadPtr := flag.Bool("keep-original-payload", false, "keep original payload")

	flag.Parse()

	fmt.Println("This is the UDP port: ", *udpPortPtr)
	fmt.Println("This is the HTTP port: ", *httpPortPtr)
	fmt.Println("This is host", *hostPtr)

	config := tracing.Config{
		TimestampFieldNames:     splitNames(timestampFieldNamesPtr),
		MessageFieldNames:       splitNames(messageFieldNamesPtr),
		LevelFieldNames:         splitNames(levelFieldNamesPtr),
		CorrelationIdFieldNames: splitNames(corridFieldNamesPtr),
		IndexableFieldNames:     splitNames(indexableFieldNamesPtr),
		KeepOriginalPayload:     *keepOriginalPayloadPtr,
	}

	store, err := tracing.NewInMemoryStore(&config)
	handleErrorNot(err)
	parser := tracing.NewPayloadParserWithConfig(&config)

	dispatch := make(chan string, 200)
	defer close(dispatch)
	go listenUdp(*udpPortPtr, *hostPtr, dispatch)
	go readFrom(store, parser, dispatch)
	api := NewTraceApi(*httpPortPtr, *hostPtr, &config, store)
	api.Start()
	defer api.Stop(context.Background())
	_, _ = fmt.Scanln() // wait for user input
}

// comma
func splitNames(cfg *string) []string {
	if cfg == nil || *cfg == "" {
		return []string{}
	}
	return strings.Split(*cfg, ",")
}

func readFrom(store tracing.TraceStore,
	parser *tracing.PayloadParser,
	dispatch <-chan string) {
	for dispatchData := range dispatch {
		trc, err := parser.Parse(dispatchData)
		if err != nil {
			fmt.Println("Could not parse: ", dispatchData, err.Error())
		} else {
			err = store.Store(trc, dispatchData)
			if err != nil {
				fmt.Println("Could not store: ", trc, err.Error())
			}
		}
	}
}

func listenUdp(port int, host string, dispatch chan<- string) {

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(host),
	})

	handleErrorNot(err)

	defer conn.Close()
	fmt.Printf("UDP listening at %s\n", conn.LocalAddr().String())

	buffer := make([]byte, 64*1024)
	for {
		len, _, err := conn.ReadFromUDP(buffer[:])
		handleErrorNot(err)

		data := strings.TrimSpace(string(buffer[:len]))
		dispatch <- data
	}
}

// for now we just panic
func handleErrorNot(e error) {
	if e != nil {
		panic(e)
	}
}
