package servers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hfpublic/webbasic/configs"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	Config     *configs.Server
	HTTPServer http.Handler
	RPCServer  http.Handler
	server     *http.Server
	CloseFunc  func(*Server)
}

func (s *Server) handlerFunc(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		s.RPCServer.ServeHTTP(w, r)
	} else {
		s.HTTPServer.ServeHTTP(w, r)
	}
}

func (s *Server) Start() error {
	servURL := fmt.Sprintf("%s:%s", s.Config.HTTPServer.Host, strconv.FormatInt(int64(s.Config.HTTPServer.Port), 10))
	// 服务器添加http2支持
	h2Handler := h2c.NewHandler(http.HandlerFunc(s.handlerFunc), &http2.Server{})
	s.server = &http.Server{
		Addr:           servURL,
		Handler:        h2Handler,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s.server.ListenAndServe()
}

func (s *Server) WatchSignal() {
	defer s.CloseFunc(s)

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTSTP)
	sig := <-ch

	fmt.Println("got a signal", sig)
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := s.server.Shutdown(cxt)
	if err != nil {
		fmt.Println("err", err)
	}

	// 看看实际退出所耗费的时间
	log.Println("exited")
}
