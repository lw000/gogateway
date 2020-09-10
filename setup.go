package main

import (
	"context"
	"fmt"
	"github.com/lw000/gocommon/app/gin"
	"github.com/lw000/gocommon/web/gin/middleware"
	"gogateway/config"
	"gogateway/frontend"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type GinServer struct {
	serves []*http.Server
	gWait  errgroup.Group
}

func (gs *GinServer) Start(debug int64, serveConf []config.Server) error {
	for _, serveConf := range serveConf {
		serveConfCopy := serveConf

		serve := &http.Server{
			Addr:    fmt.Sprintf(":%d", serveConfCopy.Listen),
			Handler: gs.getHandler(debug, &serveConfCopy),
		}
		gs.serves = append(gs.serves, serve)

		ssl := strings.ToLower(serveConfCopy.Ssl)
		switch ssl {
		case "on":
			gs.gWait.Go(func() error {
				log.Info("Listening and serving HTTPS on ", serve.Addr)
				if err := serve.ListenAndServeTLS(serveConfCopy.SslCertfile, serveConfCopy.SslKeyfile); err != nil {
					log.Error(err)
					return err
				}
				return nil
			})
		default:
			gs.gWait.Go(func() error {
				log.Info("Listening and serving HTTP on ", serve.Addr)
				if err := serve.ListenAndServe(); err != nil {
					log.Error(err)
					return err
				}
				return nil
			})
		}
	}
	if len(serveConf) > 0 {
		go func() {
			if err := gs.gWait.Wait(); err != nil {
				log.Error(err)
			}
		}()
	}
	return nil
}

func (gs *GinServer) Stop() {
	for _, serve := range gs.serves {
		err := serve.Shutdown(context.Background())
		if err != nil {
			log.Error(err)
		}
	}
}

func (gs *GinServer) getHandler(debug int64, serveConf *config.Server) http.Handler {
	app := tygin.NewApplication(debug)
	app.Init()
	app.Engine().Use(
		// 跨域处理
		tymiddleware.CorsHandler(nil),
		// 主机域名绑定
		tymiddleware.HostBindingHandler(serveConf.Servername),
	)

	engine := app.Engine()

	frontend.RegisterFrontendService(engine)

	return engine
}
