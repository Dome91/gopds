package web

import (
	"encoding/xml"
	"github.com/gofiber/fiber/v2"
	"gopds/domain"
	"gopds/services"
)

type OPDSHandler struct {
	generateOPDSRootFeed        services.GenerateOPDSRootFeed
	generateOPDSAllFeed         services.GenerateOPDSAllFeed
	generateOPDSDirectoriesFeed services.GenerateOPDSDirectoriesFeed
	generateOPDSFeedByID        services.GenerateOPDSFeedByID
	fetchCatalogEntryByID       services.FetchCatalogEntryByID
}

func NewOPDSHandler(
	generateOPDSRootFeed services.GenerateOPDSRootFeed,
	generateOPDSAllFeed services.GenerateOPDSAllFeed,
	generateOPDSFoldersFeed services.GenerateOPDSDirectoriesFeed,
	generateOPDSFeedByID services.GenerateOPDSFeedByID,
	fetchCatalogEntryByID services.FetchCatalogEntryByID,
) *OPDSHandler {
	return &OPDSHandler{
		generateOPDSRootFeed:        generateOPDSRootFeed,
		generateOPDSAllFeed:         generateOPDSAllFeed,
		generateOPDSDirectoriesFeed: generateOPDSFoldersFeed,
		generateOPDSFeedByID:        generateOPDSFeedByID,
		fetchCatalogEntryByID:       fetchCatalogEntryByID,
	}
}

func (o *OPDSHandler) Register(app *fiber.App, authorization *Authorization) {
	group := app.Group("/opds")
	group.Use(authorization.BasicAuth)
	group.Get("", o.root)
	group.Get("/all", o.all)
	group.Get("/directories", o.directories)
	group.Get("/:id", o.byID)
	group.Get("/:id/download", o.download)
}

func (o *OPDSHandler) root(ctx *fiber.Ctx) error {
	feed := o.generateOPDSRootFeed()
	return sendFeedAsXML(ctx, feed)
}

func (o *OPDSHandler) all(ctx *fiber.Ctx) error {
	feed, err := o.generateOPDSAllFeed()
	if err != nil {
		return err
	}

	return sendFeedAsXML(ctx, feed)
}

func (o *OPDSHandler) directories(ctx *fiber.Ctx) error {
	feed, err := o.generateOPDSDirectoriesFeed()
	if err != nil {
		return err
	}

	return sendFeedAsXML(ctx, feed)
}

func (o *OPDSHandler) byID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	feed, err := o.generateOPDSFeedByID(id)
	if err != nil {
		return err
	}

	return sendFeedAsXML(ctx, feed)
}

func (o *OPDSHandler) download(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	catalogEntry, err := o.fetchCatalogEntryByID(id)
	if err != nil {
		return err
	}

	ctx.Set(fiber.HeaderContentDisposition, "filename="+catalogEntry.Name)
	return ctx.SendFile(catalogEntry.Path)
}

func sendFeedAsXML(ctx *fiber.Ctx, feed domain.Feed) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationXML)
	bytes, err := xml.Marshal(feed)
	if err != nil {
		return err
	}

	return ctx.Send(bytes)
}
