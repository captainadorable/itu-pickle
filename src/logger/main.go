package logger

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

func NewLogger() *Logger {
	logger := &Logger{
		Broadcast: make(chan string),
		Clients:   make(map[*websocket.Conn]chan string),
	}
	go logger.Broadcaster()

	return logger
}

func (l *Logger) Broadcaster() {
	for msg := range l.Broadcast {
		log.Println("Broadcasting")
		l.mu.Lock()
		for client, ch := range l.Clients {
			select {
			case ch <- msg:
        log.Println("Message sent to client: ", client.RemoteAddr())
			default:
				log.Printf("Client channel full, disconnecting client: %v\n", client.RemoteAddr())
				close(ch)
				client.Close()
				delete(l.Clients, client)
			}

		}
		l.mu.Unlock()
	}
}

func (l *Logger) handleClient(conn *websocket.Conn) {
	msgChan := make(chan string, 200) // fixed size to prevent slow clients
	l.mu.Lock()
	l.Clients[conn] = msgChan
	l.mu.Unlock()
	log.Println("Client connected")
  l.Log("Sunucuya bağlanıldı")

	go func() {
		defer func() {
      l.disconnectClient(conn)
		}()

		for msg := range msgChan {
			if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
				log.Printf("Error writing to client: %v: %v", conn.RemoteAddr(), err)
        l.disconnectClient(conn)
				return
			}
		}
	}()
}
func (l *Logger) disconnectClient(conn *websocket.Conn) {
	conn.Close()
	l.mu.Lock()
	delete(l.Clients, conn)
	l.mu.Unlock()
  log.Println("Client disconnected: ", conn.RemoteAddr())
}
func (l *Logger) WebSocketHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	l.handleClient(conn)
	return nil
}

func (l *Logger) Log(msg string) {
	// log current time and msg
	message := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), msg)
	l.Messages = append(l.Messages, message)

	var buffer bytes.Buffer
	err := l.SendComponent(views.Messages([]string{message}))
	if err != nil {
    log.Println("Error broadcasting component")
		return
	}

	l.Broadcast <- buffer.String()
}

func (l *Logger) SendComponent(component templ.Component) error {
	var buffer bytes.Buffer
	err := component.Render(context.Background(), &buffer)

	if err != nil {
		return err
	}

	l.Broadcast <- buffer.String()
  return nil
}
