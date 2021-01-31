package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
	"gopds/configuration"
	"net/http"
)

type Handler interface {
	Register(app *fiber.App, authorization *Authorization)
}

type Server struct {
	app *fiber.App
}

func NewServer(authorization *Authorization, handlers ...Handler) *Server {
	app := fiber.New()
	app.Use(compress.New())
	app.Use(recover.New())
	app.Use(logger.New())

	for _, handler := range handlers {
		handler.Register(app, authorization)
	}

	serveUI(app)
	return &Server{app: app}
}

func (s *Server) Start() {
	err := s.app.Listen(configuration.Address + ":" + configuration.Port)
	if err != nil {
		log.Warn(err)
	}
}

func serveUI(app *fiber.App) {
	app.Use("/*", filesystem.New(filesystem.Config{
		Root:         http.Dir("public"),
		NotFoundFile: "index.html",
	}))
}
