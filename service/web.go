package service

import (
	"net/http"
	//"net/http/httputil"
	//"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"

	"github.com/binacsgo/log"

	"github.com/OpenKikCoc/raftkv/config"
)

// WebService the web service
type WebService interface {
	Serve() error
}

// WebServiceImpl inplement of WebService
type WebServiceImpl struct {
	Config *config.Config `inject-name:"Config"`
	Logger log.Logger     `inject-name:"WebLogger"`

	r *gin.Engine
	s *http.Server
}

// AfterInject inject
func (ws *WebServiceImpl) AfterInject() error {
	ws.r = gin.New()
	ws.r.Use(gin.Recovery())
	ws.r.Use(ws.tlsTransfer())
	ws.setRouter(ws.r)
	ws.s = &http.Server{
		Addr:           ":" + ws.Config.WebConfig.HTTPPort,
		Handler:        ws.r,
		ReadTimeout:    time.Second,
		WriteTimeout:   time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return nil
}

// Serve start web serve
func (ws *WebServiceImpl) Serve() error {
	ws.Logger.Info("WebService Serve", "HTTPPort", ws.Config.WebConfig.HTTPPort, "HttpsPort", ws.Config.WebConfig.HTTPSPort)
	go func() {
		if err := ws.s.ListenAndServe(); err != nil {
			ws.Logger.Error("WebService Serve", "ListenAndServe err", err)
		}
	}()
	err := ws.r.RunTLS(":"+ws.Config.WebConfig.HTTPSPort, ws.Config.WebConfig.CertPath, ws.Config.WebConfig.KeyPath)
	if err != nil {
		ws.Logger.Error("WebService Serve", "ListenAndServeTLS err", err)
		return err
	}
	return nil
}

func (ws *WebServiceImpl) tlsTransfer() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     ws.Config.WebConfig.Host + ":" + ws.Config.WebConfig.HTTPSPort,
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			ws.Logger.Error("WebService tlsTransfer", "Process err", err)
			c.Abort()
			return
		}
		c.Next()
	}
}

// ------------------ Gin Router ------------------

// setRouter set all router
func (ws *WebServiceImpl) setRouter(r *gin.Engine) {
	//ws.setBasicRouter(r)
}
