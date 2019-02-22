package network

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/ntfox0001/svrLib/commonError"
	"github.com/ntfox0001/svrLib/network/networkInterface"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ntfox0001/svrLib/log"

	"github.com/gorilla/mux"
)

type Server struct {
	ip                string
	port              string
	router            *mux.Router
	routerMap         map[string]networkInterface.IRouterHandler
	quitChan          chan interface{}
	ListenAndServeCh  chan error
	server            *http.Server
	ssl               bool
	certFile, keyFile string
	twoWay            bool
	caPath            []string
	ReadWriteTimeout  time.Duration
	IdleTimeout       time.Duration
}

func NewServer(ip, port string) *Server {

	svr := NewServerSsl(ip, port, "", "")

	return svr
}
func NewServerSsl(ip, port, certFile, keyFile string) *Server {
	svr := &Server{
		ip:               ip,
		port:             port,
		ssl:              certFile != "" && keyFile != "",
		twoWay:           false,
		certFile:         certFile,
		keyFile:          keyFile,
		router:           mux.NewRouter(),
		routerMap:        make(map[string]networkInterface.IRouterHandler),
		quitChan:         make(chan interface{}, 1),
		ListenAndServeCh: make(chan error),
		caPath:           make([]string, 0, 5),
		ReadWriteTimeout: time.Second * 5,
		IdleTimeout:      time.Second,
	}

	return svr
}

// 开始运行服务器，不阻塞
func (s *Server) Start() error {
	log.Info("Server", "listen", s.ip, "port", s.port)
	addr := fmt.Sprintf("%s:%s", s.ip, s.port)

	if s.twoWay {
		if config, err := s.loadCaFiles(); err != nil {
			return err
		} else {
			s.server = &http.Server{
				Addr:              addr,
				Handler:           s.router,
				IdleTimeout:       s.IdleTimeout,
				ReadTimeout:       s.ReadWriteTimeout,
				ReadHeaderTimeout: s.ReadWriteTimeout,
				WriteTimeout:      s.ReadWriteTimeout,
				TLSConfig:         config,
			}
		}
	} else {
		s.server = &http.Server{
			Addr:              addr,
			Handler:           s.router,
			ReadTimeout:       s.ReadWriteTimeout,
			ReadHeaderTimeout: s.ReadWriteTimeout,
			WriteTimeout:      s.ReadWriteTimeout,
			IdleTimeout:       s.IdleTimeout}
	}

	go func() {
		go func() {
			if s.ssl {
				s.ListenAndServeCh <- s.server.ListenAndServeTLS(s.certFile, s.keyFile)
			} else {
				s.ListenAndServeCh <- s.server.ListenAndServe()
			}
		}()

		select {
		case <-s.quitChan:
		case err := <-s.ListenAndServeCh:
			log.Error("network", "ListenAndServe", err.Error())
		}

		if err := s.server.Shutdown(context.Background()); err != nil {
			log.Error("network", "shutdown:", err.Error())
		}
	}()

	return nil
}

func (s *Server) loadCaFiles() (*tls.Config, error) {
	if s.twoWay {
		pool := x509.NewCertPool()

		for _, f := range s.caPath {
			caCrt, err := ioutil.ReadFile(f)
			if err != nil {
				log.Error("loadCaFile error", "err", err.Error(), "file", f)
				return nil, err
			}
			pool.AppendCertsFromPEM(caCrt)
		}

		config := &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
		return config, nil
	}
	return nil, commonError.NewStringErr2("two way is false.")
}

func (s *Server) AddCaFile(ca string) {
	s.twoWay = true
	s.caPath = append(s.caPath, ca)
}

func (s *Server) Close() {
	s.quitChan <- struct{}{}
}

func (s *Server) RegisterRouter(router string, handler networkInterface.IRouterHandler) {
	s.router.Handle(router, handler)
}

func (s *Server) Router() *mux.Router {
	return s.router
}

func (s *Server) Ip() string {
	return s.ip
}
func (s *Server) Port() string {
	return s.port
}

// 在start前设置有效
func (s *Server) SetTimeout(rw time.Duration, idle time.Duration) {
	s.IdleTimeout = idle
	s.ReadWriteTimeout = rw
}
