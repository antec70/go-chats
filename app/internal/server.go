package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func ttl(dataChan <-chan string, s socketio.Conn) bool {
	for afterCh := time.After(10 * time.Second); ; {
		select {
		case d := <-dataChan:
			fmt.Println("Got:", d)
			return true
		case <-afterCh:
			fmt.Println("Time's up!")
			s.Emit("disconnect")
			return false
		}
	}

}

func NewServer() {
	router := gin.New()
	ch := make(chan string)
	server, _ := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		go ttl(ch, s)
		return nil
	})

	server.OnEvent("/", "authenticate", func(s socketio.Conn, msg string) {

		ch <- msg

		/*if msg == "2" {
			fmt.Println("connected: peter", s.ID())
			s.Emit("reply", "have "+msg)
		} else {
			s.Emit("disconnect")
		}*/

	})
	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)

		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)

		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})

	go server.Serve()
	defer server.Close()

	router.Use(GinMiddleware("http://localhost:3000"))
	router.GET("/socket.io/*any", gin.WrapH(server))

	router.Run()
}
