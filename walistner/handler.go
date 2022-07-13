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
	SenderType, MessageGroup, MessageSender, SenderName, MessageTime string
	MessageID, MessageType, MessageText, MessageCaption, MessageUri  string
}

func eventHandler(evt interface{}) {
	// Using goroutine to absorb all incoming messages with no thread block
	go func() {
		switch v := evt.(type) {
		case *events.Message:
			sender := v.Info.Chat.User

			r := Response{
				SenderType:     "",
				MessageGroup:   strconv.FormatBool(v.Info.IsGroup),
				MessageSender:  "",
				SenderName:     v.Info.PushName,
				MessageTime:    v.Info.Timestamp.String(),
				MessageID:      v.Info.ID,
				MessageType:    "",
				MessageText:    "",
				MessageCaption: "",
				MessageUri:     "",
			}

			switch sender {
			case "status":
				r.SenderType = "status"
				r.MessageSender = v.Info.PushName
			default:
				r.SenderType = "message"
				r.MessageSender = sender
			}
			//fmt.Println("Received a message!", v.Message.GetConversation())
			switch {
			// Conversation: Text
			case v.Message.Conversation != nil:
				r.MessageType = "Text"
				r.MessageText = v.Message.GetConversation()

			// Conversation: Text
			case v.Message.ExtendedTextMessage != nil:
				r.MessageType = "Text"
				info, err := json.MarshalIndent(v.Message.ExtendedTextMessage.GetText(), "", "\t")
				if err != nil {
					fmt.Println(err)
				}
				r.MessageText = string(info)

			// Conversation: Text
			case v.Message.DeviceSentMessage != nil:
				r.MessageType = "Text"
				r.MessageText = v.RawMessage.String()

			// Conversation: Text
			case v.Message.Chat != nil:
				r.MessageType = "Text"
				msgReceived := v.Message.GetChat()
				r.MessageText = fmt.Sprintf("%v", msgReceived)

			// Conversation: Image
			case v.Message.ImageMessage != nil:
				r.MessageType = "Image"
				img := v.Message.GetImageMessage()
				if v.Message.ImageMessage.Caption != nil {
					r.MessageCaption = *v.Message.ImageMessage.Caption
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
					r.MessageUri = path
				}

			// Conversation: Sticker
			case v.Message.StickerMessage != nil:
				r.MessageType = "Sticker"
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
				r.MessageType = "Audio"
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
					r.MessageUri = path
				}

			// Conversation: Video
			case v.Message.VideoMessage != nil:
				r.MessageType = "Video"
				video := v.Message.GetVideoMessage()
				if v.Message.VideoMessage.Caption != nil {
					r.MessageCaption = *v.Message.VideoMessage.Caption
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
					r.MessageUri = path
				}

			// Conversation: Document
			case v.Message.DocumentMessage != nil:
				r.MessageType = "Document"
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
					r.MessageUri = path
				}

			// Conversation: Contact
			case v.Message.ContactMessage != nil:
				r.MessageType = "Contact"
				Contact := v.Message.GetContactMessage()
				//fmt.Println(Contact.GetDisplayName())
				r.MessageText = Contact.GetVcard()

			// Conversation: Location
			case v.Message.LocationMessage != nil:
				r.MessageType = "Location"
				Location := v.Message.GetLocationMessage()
				fmt.Println(Location.GetDegreesLatitude())
				fmt.Println(Location.GetDegreesLongitude())
				fmt.Println(Location.GetAddress())

				latitude := Location.GetDegreesLatitude()
				longitud := Location.GetDegreesLongitude()
				//	address := Location.GetAddress()
				link := fmt.Sprintf("<a href='https://www.google.com/maps/@%f,%f,15z' target='_blank'>Open map</a>", latitude, longitud)
				r.MessageText = "Location: " + link

			// Conversation: Buttons Response
			case v.Message.ButtonsResponseMessage != nil:
				r.MessageType = "ButtonsResponse"
				//	fmt.Println("Button responce pressed")
				ButtonResponse := v.Message.GetButtonsResponseMessage()
				// id, _ := strconv.Atoi(ButtonResponse.GetSelectedButtonId())
				r.MessageText = ButtonResponse.GetSelectedButtonId() // string(id)

			// Conversation: List Response
			case v.Message.ListResponseMessage != nil:
				r.MessageType = "ListResponse"
				// fmt.Println("List responce pressed")
				ListResponse := v.Message.GetListResponseMessage()
				// id, _ := strconv.Atoi(ListResponse.SingleSelectReply.GetSelectedRowId())
				r.MessageText = ListResponse.SingleSelectReply.GetSelectedRowId()
			}
			fmt.Println("hi")
			b, err := json.Marshal(r)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(b))
			passer.data <- sseData{
				event:   "message", // also it can be empty to be managed by source.onmessage
				message: string(b),
			}
		}
	}()
}
