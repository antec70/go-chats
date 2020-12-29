package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chats/app/internal/auth"
	"go-chats/app/internal/config"
	"log"
	"strconv"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

type Server struct {
	config config.ParamsLocal
	router *gin.Engine
}

func NewWsServ(config config.ParamsLocal) *Server {
	return &Server{
		config: config,
		router: gin.New(),
	}
}

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", " GET")
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
			s.Emit("disconnect")
			s.Close()
			fmt.Println("bye non auth user")
		}
	}

}

func (ws *Server) NewServer() error {
	server, er := socketio.NewServer(nil)
	if er != nil {
		log.Fatal(er)
	}
	ch := make(chan string)

	server.OnConnect("/chat", func(s socketio.Conn) error {
		fmt.Println("connected:", s.ID())
		go ttl(ch, s)
		return nil
	})

	server.OnEvent("/chat", "authenticate", func(s socketio.Conn, msg map[string]string) {
		ch <- msg["token"]
		user, er := auth.GetUser(msg["token"], 0, ws.config)
		if er != nil {
			fmt.Println("Error: user not found")
			s.Emit("app-error", er)
			s.Emit("disconnect")
		} else {
			s.Join("messages-to-" + strconv.Itoa(user.ID))
			s.SetContext(user.ID)

		}

	})

	server.OnEvent("/chat", "message/send", func(s socketio.Conn, msg map[string]interface{}) {
		newM := NewMessage(ws.config)
		message, er := newM.save(msg, s.Context().(int))
		user, err := auth.GetUser("", s.Context().(int), ws.config)
		if err != nil {
			fmt.Println(err)
		}
		if er != nil {
			fmt.Println(er)
			s.Emit("app-error", er)
		} else {
			fmt.Println(message)
			newM.publish(message, user)
		}
	})

	server.OnEvent("/chat", "message/read", func(s socketio.Conn, msg map[string]interface{}) {

	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		fmt.Println("User: ", last)
		s.Close()
		return last
	})

	server.OnDisconnect("/chat", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})

	go server.Serve()
	defer server.Close()

	ws.router.Use(GinMiddleware("http://localhost:3000"))
	ws.router.GET("/socket.io/*any", gin.WrapH(server))

	return ws.router.Run(":5000")
}
