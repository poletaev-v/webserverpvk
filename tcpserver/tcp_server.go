package tcpserver

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"
)

type TCPserver struct {
	AddrHTTP    string
	PortHTTP    string
	EndpointURL string
	BufferLimit int
	AwaitConn   time.Duration
}

func (s *TCPserver) Run(addr, port string, bufLimit int, awaitConn time.Duration) error {
	log.Println("TCP server starting...")
	listener, err := net.Listen("tcp", addr+":"+port)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("TCP server is listening %v", addr+":"+port)
	for {
		// Open connection
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			conn.Close()
			continue
		}
		go s.handleConnection(conn, bufLimit, awaitConn)
	}
}

func (s *TCPserver) handleConnection(conn net.Conn, bufLimit int, awaitConn time.Duration) {
	var xmlParser XMLParser
	log.Printf("Client connection addr: %v", conn.RemoteAddr())

	defer func() {
		log.Printf("Closing connection from %v", conn.RemoteAddr())
		conn.Close()
	}()

	// Wait n seconds and close connection with client
	conn.SetDeadline(time.Now().Add(awaitConn))
	for {
		// Get data from xml file client
		input := make([]byte, bufLimit<<20)
		n, err := conn.Read(input)
		if n == 0 || err != nil {
			break
		}
		// If not active action
		// Wait n seconds and close connection with client
		conn.SetDeadline(time.Now().Add(awaitConn))
		buffXML := input[:n]

		// Parsing data xml
		pXML, err := xmlParser.parse(buffXML)
		if err != nil {
			log.Println(err)
			break
		}

		values := xmlParser.getValues(pXML)
		b, err := json.Marshal(values)
		if err != nil {
			log.Println(err)
			break
		}

		// Prepare request for next sending
		req, err := http.NewRequest("POST", "http://"+s.AddrHTTP+":"+s.PortHTTP+"/"+s.EndpointURL, bytes.NewBuffer(b))
		if err != nil {
			log.Println(err)
			break
		}
		client := &http.Client{}
		// Post request to HTTP server
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			break
		}
		log.Println(resp.StatusCode)
	}
}
