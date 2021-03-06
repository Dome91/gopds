package web

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	log "github.com/sirupsen/logrus"
	"gopds/configuration"
	"gopds/public"
	"net/http"
	"os"
	"os/signal"
	"path"
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

	serveCovers(app)
	serveUI(app)
	return &Server{app: app}
}

func (s *Server) Start() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Shutting down GOPDS...")
		_ = s.app.Shutdown()
	}()

	err := s.app.Listen(configuration.Address + ":" + configuration.Port)
	if err != nil {
		log.Warn(err)
	}
}

func serveCovers(app *fiber.App) {
	app.Use("/covers/*", filesystem.New(filesystem.Config{
		Root: http.Dir(path.Join("data")),
	}))
}

func serveUI(app *fiber.App) {
	app.Use("/*", filesystem.New(filesystem.Config{
		Root:         http.FS(public.FS),
		NotFoundFile: "index.html",
	}))
}
