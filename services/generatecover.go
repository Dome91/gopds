package services

import (
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/mholt/archiver/v3"
	"github.com/mustafaturan/bus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopds/domain"
	"image"
	"image/jpeg"
	"io"
	"path"
)

type GenerateCover func(e *bus.Event)

var ErrEmptyCatalogEntry = errors.New("catalog entry is empty")

func RegisterGenerateCover(b domain.Bus, generateCover GenerateCover) {
	handler := bus.Handler{Matcher: domain.GenerateCoverTopic, Handle: func(e *bus.Event) {
		go generateCover(e)
	}}

	b.RegisterHandler("GenerateCoverHandler", &handler)
}

func GenerateCoverProvider(fs afero.Fs, repository domain.CatalogRepository) GenerateCover {
	return func(e *bus.Event) {
		event, ok := e.Data.(domain.GenerateCoverEvent)
		if !ok {
			log.Error("malformed event body. Expected GenerateCoverEvent")
			return
		}

		entry, err := repository.FindByID(event.ID)
		if err != nil {
			log.Error(err)
			return
		}

		coverFilenameInCatalogEntry, err := getNameOfCoverFile(entry)
		if err != nil {
			log.Error(err)
			return
		}

		err = domain.ForEveryFileInCatalogEntryDO(entry, domain.OnlyImages, func(file archiver.File) error {
			if file.FileInfo.Name() != coverFilenameInCatalogEntry {
				return nil
			}

			cover, err := resizeImage(file, 800)
			if err != nil {
				return err
			}

			filename, err := saveCover(cover, fs)
			if err != nil {
				return err
			}

			entry.Cover = filename
			return archiver.ErrStopWalk
		})

		if err != nil {
			log.Error(err)
			return
		}

		err = repository.UpdateSetCoverByID(entry.ID, entry.Cover)
		if err != nil {
			log.Error(err)
			return
		}

		log.Infof("generated cover for %s", entry.Name)
	}
}

func getNameOfCoverFile(entry domain.CatalogEntry) (string, error) {
	filepaths, err := domain.GetFilenamesInCatalogEntryInAlphabeticalOrder(entry, domain.OnlyImages)
	if err != nil {
		return "", err
	}

	if len(filepaths) == 0 {
		return "", errors.New(fmt.Sprintf("%s :%v", entry.Name, ErrEmptyCatalogEntry))
	}

	return filepaths[0], nil
}

func resizeImage(r io.Reader, targetHeight int) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	return imaging.Resize(img, 0, targetHeight, imaging.Box), nil
}

func saveCover(img image.Image, fs afero.Fs) (string, error) {
	filename := uuid.New().String() + ".jpg"
	filepath := path.Join("data", "covers", filename)
	coverFile, err := fs.Create(filepath)
	if err != nil {
		return "", err
	}
	defer coverFile.Close()

	err = jpeg.Encode(coverFile, img, nil)
	if err != nil {
		return "", err
	}

	return filename, nil
}
