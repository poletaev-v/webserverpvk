package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.100.4:4545")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	for {
		var filename string
		XMLData := make([]byte, 1<<20)

		fmt.Println("Введите название xml файла: ")
		_, err := fmt.Scanln(&filename)
		if err != nil {
			fmt.Println("Некорректный ввод", err)
			continue
		}
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println(err)
			break
		}

		cnt, err := file.Read(XMLData)
		if cnt == 0 || err != nil {
			fmt.Println(err)
			break
		}

		// Send message to server
		if n, err := conn.Write([]byte(XMLData)); n == 0 || err != nil {
			fmt.Println(err)
			return
		}
		// Get response
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			break
		}
		fmt.Println(string(buff[:n]))
	}
}
