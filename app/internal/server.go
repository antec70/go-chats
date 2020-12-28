package internal

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chats/app/internal/auth"
	"go-chats/app/internal/config"

	"time"

	socketio "github.com/googollee/go-socket.io"
)

type Server struct {
	config config.ParamsLocal
	router *gin.Engine
	socket *socketio.Server
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
			s.Emit("disconnect")
			fmt.Println("bye non auth user")
		}
	}

}

func (ws *Server) NewServer() error {

	/*if er != nil {
		log.Fatal(er)
	}*/
	ch := make(chan string)
	//server, _ := socketio.NewServer(nil)

	ws.socket.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		go ttl(ch, s)
		return nil
	})

	ws.socket.OnEvent("/", "authenticate", func(s socketio.Conn, msg map[string]string) {

		ch <- msg["token"]
		user, er := auth.GetUser(msg["token"], ws.config)

		if er != nil {
			fmt.Println(er)
			s.Emit("disconnect")
		}
		fmt.Println("User: ", user)

	})
	ws.socket.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)

		s.Emit("reply", "have "+msg)
	})

	ws.socket.OnEvent("/", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)

		return "recv " + msg
	})

	ws.socket.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	ws.socket.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	ws.socket.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})

	go ws.socket.Serve()
	defer ws.socket.Close()

	ws.router.Use(GinMiddleware("http://localhost:3000"))
	ws.router.GET("/socket.io/*any", gin.WrapH(ws.socket))

	return ws.router.Run(":5000")
}
