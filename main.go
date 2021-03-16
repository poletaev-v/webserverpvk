// Copyright (c) 2021 Vitaliy Poletaev <hapdev22@gmail.com>. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:

//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/happydev/webserverpvk/tcpserver"
)

const pathConfigs string = "configs/main"

func main() {
	var cfg Config
	var ttl int64

	err := setConfigs(pathConfigs, &cfg)
	if err != nil {
		panic(err)
	}

	// Serve TCP server
	tcpsrv := tcpserver.TCPserver{
		AddrHTTP:    cfg.HTTP.Addr,
		PortHTTP:    cfg.HTTP.Port,
		EndpointURL: cfg.HTTP.RefreshURL,
		BufferLimit: cfg.TCP.BufferLimit,
		AwaitConn:   cfg.TCP.AwaitConn,
	}

	go func() {
		log.Fatal(tcpsrv.Run(cfg.TCP.Addr, cfg.TCP.Port, cfg.TCP.BufferLimit, cfg.TCP.AwaitConn))
	}()

	// Serve HTTP server
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.Static("/fonts", "./fonts")
	router.Static("/templates", "./templates")
	router.StaticFile("/favicon.ico", "./favicon.ico")
	s := &http.Server{
		Addr:           cfg.HTTP.Addr + ":" + cfg.HTTP.Port,
		Handler:        router,
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		MaxHeaderBytes: cfg.HTTP.MaxHeaderMIB << 20, // MIB
	}

	router.GET("/", func(c *gin.Context) {
		router.LoadHTMLGlob("templates/*")

		if ttl == 0 {
			html := "<div class='message'>" + cfg.DATA.message + "</div>"
			f, _ := os.OpenFile("templates/content.html", os.O_WRONLY, 0664)
			defer f.Close()
			// Clear file before write data
			f.Truncate(0)
			_, err := f.Write([]byte(html))
			if err != nil {
				log.Println(err)
			}
			c.HTML(http.StatusOK, "index.html", html)
		} else if ttl > 0 && time.Now().Unix() > ttl {
			ttl = 0
			c.HTML(http.StatusOK, "index.html", gin.H{})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		}
	})

	router.POST("/"+cfg.HTTP.RefreshURL, func(c *gin.Context) {
		var objMAP map[string]interface{}
		var bc []byte
		result := make(map[string]interface{})

		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(b, &objMAP)
		if err != nil {
			log.Println(err)
		}

		for key, val := range cfg.DATA.availableFields {
			if item, ok := objMAP[key]; ok {
				result[val] = item
			}
		}

		if len(result) > 0 {
			f, _ := os.OpenFile("templates/content.html", os.O_WRONLY, 0664)
			defer f.Close()
			html := ""
			for key, val := range result {
				tmpVal := ""
				switch valT := val.(type) {
				case string:
					log.Println("string", valT)
					tmpVal = valT
				case float64:
					tmpVal = fmt.Sprintf("%.f", valT)
				case int:
					log.Println("int", valT)
					tmpVal = strconv.Itoa(valT)
				case uint64:
					log.Println("uint64", valT)
					tmpVal = strconv.FormatUint(valT, 10)
				case bool:
					if valT {
						tmpVal = "Да"
					} else {
						tmpVal = "Нет"
					}
				}
				html += "<div class='pvk-item'>" + key + ": " + tmpVal + "</div>"
			}
			bc = []byte(html)
			// Clear file before write data
			f.Truncate(0)
			_, err = f.Write(bc)
			if err != nil {
				log.Println(err)
			}
			// Set ttl for await picture to screen
			ttl = time.Now().Unix() + cfg.DATA.timeDuration
		}
	})

	log.Fatal(s.ListenAndServe())
}
