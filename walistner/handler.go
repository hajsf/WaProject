package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"go.mau.fi/whatsmeow/types/events"
)

type Response struct {
	senderType, messageGroup, messageSender, senderName, messageTime string
	messageID, messageType, messageText, messageCaption, messageUri  string
}

func eventHandler(evt interface{}) {
	// Using goroutine to absorb all incoming messages with no thread block
	go func() {
		switch v := evt.(type) {
		case *events.Message:
			sender := v.Info.Chat.User

			r := Response{
				senderType:     "",
				messageGroup:   strconv.FormatBool(v.Info.IsGroup),
				messageSender:  "",
				senderName:     v.Info.PushName,
				messageTime:    v.Info.Timestamp.String(),
				messageID:      v.Info.ID,
				messageType:    "",
				messageText:    "",
				messageCaption: "",
				messageUri:     "",
			}

			switch sender {
			case "status":
				r.senderType = "status"
				r.messageSender = v.Info.PushName
			default:
				r.senderType = "message"
				r.messageSender = sender
			}
			//fmt.Println("Received a message!", v.Message.GetConversation())
			switch {
			// Conversation: Text
			case v.Message.Conversation != nil:
				r.messageType = "Text"
				r.messageText = v.Message.GetConversation()

			// Conversation: Text
			case v.Message.ExtendedTextMessage != nil:
				r.messageType = "Text"
				info, err := json.MarshalIndent(v.Message.ExtendedTextMessage.GetText(), "", "\t")
				if err != nil {
					fmt.Println(err)
				}
				r.messageText = string(info)

			// Conversation: Text
			case v.Message.DeviceSentMessage != nil:
				r.messageType = "Text"
				r.messageText = v.RawMessage.String()

			// Conversation: Text
			case v.Message.Chat != nil:
				r.messageType = "Text"
				msgReceived := v.Message.GetChat()
				r.messageText = fmt.Sprintf("%v", msgReceived)

			// Conversation: Image
			case v.Message.ImageMessage != nil:
				r.messageType = "Image"
				img := v.Message.GetImageMessage()
				if v.Message.ImageMessage.Caption != nil {
					r.messageCaption = *v.Message.ImageMessage.Caption
				}

				if img != nil {
					file, err := client.Download(img)
					if err != nil {
						log.Printf("Failed to download image: %v", err)
						return
					}
					exts, _ := mime.ExtensionsByType(img.GetMimetype())
					path := fmt.Sprintf("D:/Deployment/DigitalAssistance/Downloads/Image/%s-%s%s", sender, v.Info.ID, exts[0])
					err = os.WriteFile(path, file, 0600)
					if err != nil {
						log.Printf("Failed to save image: %v", err)
						return
					}
					log.Printf("Saved image in message to %s", path)
					r.messageUri = path
				}

			// Conversation: Sticker
			case v.Message.StickerMessage != nil:
				r.messageType = "Sticker"
				sticker := v.Message.GetStickerMessage()
				if sticker.Url != nil {
					fmt.Println(sticker)
					/*	data, err := global.Cli.Download(audio)
						if err != nil {
							log.Printf("Failed to download audio: %v", err)
							return
						} */
					//	global.Passer.Logs <- fmt.Sprintf("Sticker: <a href='%v' target='_blank'>Open</a>", sticker.GetUrl())
				}

			// Conversation: Audio
			case v.Message.AudioMessage != nil:
				r.messageType = "Audio"
				audio := v.Message.GetAudioMessage()
				if audio != nil {
					file, err := client.Download(audio)
					if err != nil {
						log.Printf("Failed to download audio: %v", err)
						return
					}
					exts, _ := mime.ExtensionsByType(audio.GetMimetype())
					path := fmt.Sprintf("D:/Deployment/DigitalAssistance/Downloads/Audio/%s-%s%s", sender, v.Info.ID, exts[0])
					err = os.WriteFile(path, file, 0600)
					if err != nil {
						log.Printf("Failed to save audio: %v", err)
						return
					}
					log.Printf("Saved audio in message to %s", path)
					r.messageUri = path
				}

			// Conversation: Video
			case v.Message.VideoMessage != nil:
				r.messageType = "Video"
				video := v.Message.GetVideoMessage()
				if v.Message.VideoMessage.Caption != nil {
					r.messageCaption = *v.Message.VideoMessage.Caption
				}

				if video != nil {
					file, err := client.Download(video)
					if err != nil {
						log.Printf("Failed to download video: %v", err)
						return
					}
					exts, _ := mime.ExtensionsByType(video.GetMimetype())
					path := fmt.Sprintf("D:/Deployment/DigitalAssistance/Downloads/Video/%s-%s%s", sender, v.Info.ID, exts[0])
					err = os.WriteFile(path, file, 0600)
					if err != nil {
						log.Printf("Failed to save video: %v", err)
						return
					}
					log.Printf("Saved video in message to %s", path)
					r.messageUri = path
				}

			// Conversation: Document
			case v.Message.DocumentMessage != nil:
				r.messageType = "Document"
				document := v.Message.GetDocumentMessage()
				if document != nil {
					file, err := client.Download(document)
					if err != nil {
						log.Printf("Failed to download audio: %v", err)
						return
					}
					exts, _ := mime.ExtensionsByType(document.GetMimetype())
					path := fmt.Sprintf("D:/Deployment/DigitalAssistance/Downloads/Documents/%s-%s%s", sender, v.Info.ID, exts[0])
					err = os.WriteFile(path, file, 0600)
					if err != nil {
						log.Printf("Failed to save document: %v", err)
						return
					}
					log.Printf("Saved document in message to %s", path)
					r.messageUri = path
				}

			// Conversation: Contact
			case v.Message.ContactMessage != nil:
				r.messageType = "Contact"
				Contact := v.Message.GetContactMessage()
				//fmt.Println(Contact.GetDisplayName())
				r.messageText = Contact.GetVcard()

			// Conversation: Location
			case v.Message.LocationMessage != nil:
				r.messageType = "Location"
				Location := v.Message.GetLocationMessage()
				fmt.Println(Location.GetDegreesLatitude())
				fmt.Println(Location.GetDegreesLongitude())
				fmt.Println(Location.GetAddress())

				latitude := Location.GetDegreesLatitude()
				longitud := Location.GetDegreesLongitude()
				//	address := Location.GetAddress()
				link := fmt.Sprintf("<a href='https://www.google.com/maps/@%f,%f,15z' target='_blank'>Open map</a>", latitude, longitud)
				r.messageText = "Location: " + link

			// Conversation: Buttons Response
			case v.Message.ButtonsResponseMessage != nil:
				r.messageType = "ButtonsResponse"
				//	fmt.Println("Button responce pressed")
				ButtonResponse := v.Message.GetButtonsResponseMessage()
				// id, _ := strconv.Atoi(ButtonResponse.GetSelectedButtonId())
				r.messageText = ButtonResponse.GetSelectedButtonId() // string(id)

			// Conversation: List Response
			case v.Message.ListResponseMessage != nil:
				r.messageType = "ListResponse"
				// fmt.Println("List responce pressed")
				ListResponse := v.Message.GetListResponseMessage()
				// id, _ := strconv.Atoi(ListResponse.SingleSelectReply.GetSelectedRowId())
				r.messageText = ListResponse.SingleSelectReply.GetSelectedRowId()
			}
		}
		passer.data <- sseData{
			event:   "message", // also it can be empty to be managed by source.onmessage
			message: "Reconnecting to WhatsApp server ...",
		}
	}()
}
