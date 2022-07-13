package main

import (
	"context"
	"fmt"
)

var err error

func Connect() {
	// fmt.Println("Connected")
	if client.IsConnected() {
		client.Disconnect()
		passer.data <- sseData{
			event:   "notification",
			message: "Reconnecting to WhatsApp server ...",
		}
	}

	if client.Store.ID == nil {
		// No ID stored, new login
	GetQR:
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			//	panic(err)
			//passer.logs <- "Can not connect with WhatApp server, try again later"
			passer.data <- sseData{
				event:   "notification",
				message: "Can not connect with WhatApp server, try again later",
			}
			// fmt.Println("Sorry", err)

		}

		for evt := range qrChan {
			switch evt.Event {
			case "success":
				{
					//passer.logs <- "success"
					passer.data <- sseData{
						event:   "notification",
						message: "success",
					}
					// fmt.Println("Login event: success")
				}
			case "timeout":
				{
					//passer.logs <- "timeout/Refreshing"
					passer.data <- sseData{
						event:   "notification",
						message: "Timeout or error reading from WhatsApp websocket, trying refreshing ...",
					}
					// fmt.Println("Login event: timeout")
					goto GetQR
				}
			case "code":
				{
					//fmt.Println("new code recieved")
					fmt.Println(evt.Code)
					//passer.logs <- evt.Code
					passer.data <- sseData{
						event:   "qrCode",
						message: evt.Code,
					}
				}
			}
		}
	} else {
		// Already logged in, just connect
		//passer.logs <- "Already logged"
		passer.data <- sseData{
			event:   "notification",
			message: "Already logged in",
		}
		//fmt.Println("Already logged")
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}
}
