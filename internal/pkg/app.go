package pkg

import (
	"fmt"

	"patients_categories/internal/app/config"
	"patients_categories/internal/app/handler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Application struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(cfg *config.Config, router *gin.Engine, h *handler.Handler) *Application {
	return &Application{
		Config:  cfg,
		Router:  router,
		Handler: h,
	}
}

func (a *Application) RunApp() {
	logrus.Info("Server start up")
	a.Handler.RegisterHandler(a.Router)
	a.Handler.RegisterStatic(a.Router)

	addr := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	if err := a.Router.Run(addr); err != nil {
		logrus.Fatal(err)
	}
}
