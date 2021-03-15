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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/happydev/webserverpvk/tcpserver"
)

const (
	refreshURL  string = "VcjZt3ZMwbmzvtif3Q9WzrjhZnZefctK5Cxw0N5vzSrl5HTtrMMWTTEIHAQbh7TThzyBgUf6lNMdU6VMtKVq2TamtDBpsvcEOFggFnE5UXFZ3sWGHIxMpPGO"
	pathConfigs string = "configs/main"
)

func main() {
	var cfg Config

	err := setConfigs(pathConfigs, &cfg)
	if err != nil {
		panic(err)
	}

	// Serve TCP server
	tcpsrv := tcpserver.TCPserver{
		AddrHTTP:    cfg.HTTP.Addr,
		PortHTTP:    cfg.HTTP.Port,
		EndpointURL: refreshURL,
		BufferLimit: cfg.TCP.BufferLimit,
		AwaitConn:   cfg.TCP.AwaitConn,
	}

	go func() {
		log.Fatal(tcpsrv.Run(cfg.TCP.Addr, cfg.TCP.Port, cfg.TCP.BufferLimit, cfg.TCP.AwaitConn))
	}()

	// Serve HTTP server
	router := gin.Default()
	// router.Delims("{[{", "}]}")
	router.Static("/assets", "./assets")
	router.Static("/fonts", "./fonts")
	router.Static("/data", "./data")
	router.LoadHTMLGlob("templates/*")
	s := &http.Server{
		Addr:           cfg.HTTP.Addr + ":" + cfg.HTTP.Port,
		Handler:        router,
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		MaxHeaderBytes: cfg.HTTP.MaxHeaderMIB << 20, // MIB
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"message": "Ведется весовой и габаритный контроль",
		})
	})

	router.POST("/"+refreshURL, func(c *gin.Context) {
		var bc []byte
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
		}
		// Adding timestamp to json
		if b[len(b)-1] == 125 {
			strTime := strconv.FormatInt((time.Now().UnixNano() / int64(time.Millisecond)), 10)
			bc = append(b[:len(b)-1], []byte(`,"timeStamp":`+strTime+`}`)...)
		}

		f, _ := os.OpenFile("data/data.json", os.O_WRONLY|os.O_CREATE, 0664)
		_, err = f.Write(bc)
		if err != nil {
			log.Println(err)
		}
		f.Close()
	})

	log.Fatal(s.ListenAndServe())
}
