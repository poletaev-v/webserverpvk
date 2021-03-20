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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/happydev/webserverpvk/tcpserver"
)

const pathConfigs string = "configs/main"

func main() {
	var cfg Config

	// Set configs from config file
	err := setConfigs(pathConfigs, &cfg)
	if err != nil {
		panic(err)
	}

	if len(cfg.DATA.directionTravel) == 0 {
		panic("Direction of travel not be empty, please add direction in config file")
	}

	ttlMAP := make(map[string]*int64, len(cfg.DATA.directionTravel))
	// Direction name to monitor mon1, mon2, ...
	for dName := range cfg.DATA.directionTravel {
		ttl := int64(0)
		ttlMAP[dName] = &ttl
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
	router.StaticFile("/favicon.ico", "./favicon.ico")
	s := &http.Server{
		Addr:           cfg.HTTP.Addr + ":" + cfg.HTTP.Port,
		Handler:        router,
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		MaxHeaderBytes: cfg.HTTP.MaxHeaderMIB << 20, // MIB
	}

	router.GET("/mon1", func(c *gin.Context) {
		router.LoadHTMLGlob("templates/mon1/*")

		if *ttlMAP["mon1"] == int64(0) {
			html := "<div class='message'>" + cfg.DATA.message + "</div>"
			f, _ := os.OpenFile("templates/mon1/content.html", os.O_WRONLY, 0664)
			defer f.Close()
			// Clear file before write data
			f.Truncate(0)
			_, err := f.Write([]byte(html))
			if err != nil {
				log.Println(err)
			}
			c.HTML(http.StatusOK, "index.html", html)
		} else if *ttlMAP["mon1"] > int64(0) && time.Now().Unix() > *ttlMAP["mon1"] {
			*ttlMAP["mon1"] = int64(0)
			c.HTML(http.StatusOK, "index.html", gin.H{})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		}
	})

	router.GET("/mon2", func(c *gin.Context) {
		router.LoadHTMLGlob("templates/mon2/*")

		if *ttlMAP["mon2"] == int64(0) {
			html := "<div class='message'>" + cfg.DATA.message + "</div>"
			f, _ := os.OpenFile("templates/mon2/content.html", os.O_WRONLY, 0664)
			defer f.Close()
			// Clear file before write data
			f.Truncate(0)
			_, err := f.Write([]byte(html))
			if err != nil {
				log.Println(err)
			}
			c.HTML(http.StatusOK, "index.html", html)
		} else if *ttlMAP["mon2"] > int64(0) && time.Now().Unix() > *ttlMAP["mon2"] {
			*ttlMAP["mon2"] = int64(0)
			c.HTML(http.StatusOK, "index.html", gin.H{})
		} else {
			c.HTML(http.StatusOK, "index.html", gin.H{})
		}
	})

	router.GET("/", func(c *gin.Context) {
		router.LoadHTMLFiles("templates/index.html", "templates/content.html")
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
	})

	router.POST("/"+cfg.HTTP.RefreshURL, func(c *gin.Context) {
		var objMAP map[string]interface{}
		var direction string
		var result resultData

		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
		}

		err = json.Unmarshal(b, &objMAP)
		if err != nil {
			log.Println(err)
		}

		// Get direction by monitor
		for key, val := range cfg.DATA.directionTravel {
			for _, item := range val {
				if fmt.Sprintf("%.f", objMAP[strings.ToLower(cfg.DATA.tagPlatform)].(float64)) == item {
					direction = key
				}
			}
		}

		// Get violations
		for _, viType := range cfg.DATA.violationTypes {
			if item, ok := objMAP[strings.ToLower(viType)]; ok {
				// if isset violation, then add violation to result map
				if item.(bool) {
					rkey := cfg.DATA.violationName[strings.ToLower(viType)]
					result.items = append(result.items, resultItem{
						key:       rkey,
						valueItem: "",
					})

					for _, viField := range cfg.DATA.violationValue[strings.ToLower(viType)] {
						if strings.ToLower(viType) == strings.ToLower(cfg.DATA.multipleViolationField) {
							tField, tValue, nField, nValue := getTrackThrust(objMAP, &cfg)
							if tValue != nil {
								result.items = append(result.items, resultItem{
									key:       cfg.DATA.availableFields[strings.ToLower(tField)],
									valueItem: tValue,
								})

								result.items = append(result.items, resultItem{
									key:       cfg.DATA.availableFields[strings.ToLower(nField)],
									valueItem: nValue,
								})
								break
							}
						} else {
							if item, ok := objMAP[strings.ToLower(viField)]; ok {
								var ritem resultItem
								ritem.key = cfg.DATA.availableFields[strings.ToLower(viField)]
								ritem.valueItem = item
								result.items = append(result.items, ritem)
							}
						}
					}
					// Show only one violation
					break
				}
			}
		}

		if len(result.items) > 0 && direction != "" {
			// Get transport number
			for _, tcNum := range cfg.DATA.numberTC {
				if item, ok := objMAP[strings.ToLower(tcNum)]; ok {
					if item.(string) != "" {
						result.items = append(result.items, resultItem{
							key:       cfg.DATA.availableFields[strings.ToLower(tcNum)],
							valueItem: item,
						})
						// Show only one transport number
						break
					}
				}
			}

			err = routeValues(direction, result.items, &cfg, ttlMAP[direction])
			if err != nil {
				log.Println(err)
			}
		}
	})

	log.Fatal(s.ListenAndServe())
}

type resultData struct {
	items []resultItem
}

type resultItem struct {
	key       string
	valueItem interface{}
}

func routeValues(direction string, items []resultItem, cfg *Config, ttl *int64) error {
	var bc []byte

	f, _ := os.OpenFile("templates/"+direction+"/content.html", os.O_WRONLY, 0664)
	defer f.Close()
	html := ""
	for _, val := range items {
		tmpVal := ""
		switch valT := val.valueItem.(type) {
		case string:
			tmpVal = valT
		case float64:
			tmpVal = fmt.Sprintf("%.f", valT)
		case int:
			tmpVal = strconv.Itoa(valT)
		case uint64:
			tmpVal = strconv.FormatUint(valT, 10)
		case bool:
			if valT {
				tmpVal = "Да"
			} else {
				tmpVal = "Нет"
			}
		}

		if tmpVal == "" {
			html += "<div class='pvk-item'>" + val.key + "</div>"
		} else {
			html += "<div class='pvk-item'>" + val.key + ": " + tmpVal + "</div>"
		}
	}
	bc = []byte(html)
	// Clear file before write data
	f.Truncate(0)
	_, err := f.Write(bc)
	if err != nil {
		return err
	}
	// Set ttl for await picture to screen
	*ttl = time.Now().Unix() + cfg.DATA.timeDuration
	return nil
}

// getTrackThrust - getting track thrust and bind fields
func getTrackThrust(values map[string]interface{}, cfg *Config) (string, interface{}, string, interface{}) {
	for _, val := range cfg.DATA.trackThrustes {
		if item, ok := values[strings.ToLower(val)]; ok {
			// Getting number track thrust - TrackThrust1 = 1
			ntt := val[len(val)-1:]
			if bindFieldThrust, ok := values[strings.ToLower(cfg.DATA.bindMultipleViolationField)+ntt]; ok {
				if item.(float64) > bindFieldThrust.(float64) {
					return strings.ToLower(val), item, strings.ToLower(cfg.DATA.bindMultipleViolationField) + ntt, bindFieldThrust
				}
			}
		}
	}
	return "", nil, "", nil
}
