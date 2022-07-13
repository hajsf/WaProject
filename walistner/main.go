package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client
var passer *DataPasser

const maxClients = 1

func init() {
	passer = &DataPasser{
		data:       make(chan sseData),
		logs:       make(chan string),
		connection: make(chan struct{}, maxClients),
	}
}

func init() {
	go func() {
		store.DeviceProps.Os = proto.String("Cocoon App")
		dbLog := waLog.Stdout("Database", "ERROR", true) // "DEBUG"
		// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
		container, err := sqlstore.New("sqlite3", "file:datastore.db?_foreign_keys=on", dbLog)
		if err != nil {
			panic(err)
		}
		// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
		deviceStore, err := container.GetFirstDevice()
		if err != nil {
			panic(err)
		}

		//clientLog := waLog.Stdout("Client", "ERROR", true)
		//client = whatsmeow.NewClient(deviceStore, clientLog)
		clientLog := LogText("Client", "DEBUG", true)
		client = whatsmeow.NewClient(deviceStore, clientLog)

		client.AddEventHandler(eventHandler)
	}()

	/* Trying to catch the error:
	cmd := exec.Command("tail", "-f",  "/usr/local/var/log/redis.log"

	    // create a pipe for the output of the script
		cmdReader, err :=  cmd.StdoutPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			return
		}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t > Output is: %s\n", scanner.Text())
		}
	}()
	*/

}

func main() {

	http.HandleFunc("/sse", passer.HandleSignal)
	go http.ListenAndServe(":4004", nil)

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	if client.IsConnected() {
		passer.data <- sseData{
			event:   "notification",
			message: "Server is shut down at the host machine...",
		}
		client.Disconnect()
	}
}
