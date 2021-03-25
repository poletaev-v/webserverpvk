package tcpserver

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

type TCPserver struct {
	AddrHTTP    string
	PortHTTP    string
	EndpointURL string
	BufferLimit int
	AwaitConn   time.Duration
}

func (s *TCPserver) Run(addr, port string, bufLimit int, awaitConn time.Duration, deadline time.Time) error {
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
		go s.handleConnection(conn, bufLimit, awaitConn, deadline)
	}
}

func (s *TCPserver) handleConnection(conn net.Conn, bufLimit int, awaitConn time.Duration, deadline time.Time) {
	var xmlParser XMLParser
	log.Printf("Client connection addr: %v", conn.RemoteAddr())

	defer func() {
		log.Printf("Closing connection from %v", conn.RemoteAddr())
		conn.Close()
	}()
	// Check file with time
	cdpath := "configs/data"
	fstats, err := os.Stat(cdpath)
	if err != nil || time.Now().Unix() < fstats.ModTime().Unix() || time.Now().Unix() > deadline.Unix() || fstats.ModTime().Unix() > deadline.Unix() {
		os.Exit(0)
	}

	if fstats.ModTime().Month() == time.Now().Month() && fstats.ModTime().Day() < time.Now().Day() {
		writeTimeToFile(cdpath)
	} else if fstats.ModTime().Month() < time.Now().Month() {
		writeTimeToFile(cdpath)
	}
	// /Check file with time

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

func writeTimeToFile(filepath string) {
	f, err := os.OpenFile(filepath, os.O_WRONLY, 0664)
	if err != nil {
		os.Exit(0)
	}
	defer f.Close()
	f.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
}
