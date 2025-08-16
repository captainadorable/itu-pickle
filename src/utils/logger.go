package utils

import (
	"bytes"
	"context"
	"fmt"
	views "itu-pickle/views/index"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Logger struct {
	Broadcast chan string
	Clients   map[*websocket.Conn]chan string
	Messages  []string
	mu        sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var Logcu *Logger

func StartLogger() {
	Logcu = &Logger{
		Broadcast: make(chan string),
		Clients:   make(map[*websocket.Conn]chan string),
	}
	go Broadcaster()
}

func Broadcaster() {
	for msg := range Logcu.Broadcast {
		log.Println("Broadcasting")
		Logcu.mu.Lock()
		for client, ch := range Logcu.Clients {
			select {
			case ch <- msg:
        log.Println("Message sent to client: ", client.RemoteAddr())
			default:
				log.Printf("Client channel full, disconnecting client: %v\n", client.RemoteAddr())
				close(ch)
				client.Close()
				delete(Logcu.Clients, client)
			}
		}
		Logcu.mu.Unlock()
	}
}

func handleClient(conn *websocket.Conn) {
	msgChan := make(chan string, 200) // fixed size to prevent slow clients
	Logcu.mu.Lock()
	Logcu.Clients[conn] = msgChan
	Logcu.mu.Unlock()
	log.Println("Client connected")
  Log("Sunucuya bağlanıldı")

	go func() {
		defer func() {
      disconnectClient(conn)
		}()

		for msg := range msgChan {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				log.Printf("Error writing to client: %v: %v", conn.RemoteAddr(), err)
        disconnectClient(conn)
				return
			}
		}
	}()
}
func disconnectClient(conn *websocket.Conn) {
	conn.Close()
	Logcu.mu.Lock()
	delete(Logcu.Clients, conn)
	Logcu.mu.Unlock()
  log.Println("Client disconnected: ", conn.RemoteAddr())
}
func WebSocketHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	handleClient(conn)
	return nil
}

func Log(msg string) {
	// log current time and msg
	message := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), msg)
	Logcu.Messages = append(Logcu.Messages, message)

	var buffer bytes.Buffer
	err := SendComponent(views.Messages([]string{message}))
	if err != nil {
    log.Println("Error broadcasting component")
		return
	}

	Logcu.Broadcast <- buffer.String()
}

func SendComponent(component templ.Component) error {
	var buffer bytes.Buffer
	err := component.Render(context.Background(), &buffer)

	if err != nil {
		return err
	}

	Logcu.Broadcast <- buffer.String()
  return nil
}
