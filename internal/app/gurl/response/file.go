package response

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type FileResponse interface {
	Save(url string) error
}

func SaveFile(contentType string, data []byte, url string) error {
	if contentType == "image/jpeg" {
		image := JPEGImageRespones{
			Data: data,
		}

		if err := image.Save(url); err != nil {
			return fmt.Errorf("failed to save jpeg image: %s", err.Error())
		}
	} else if contentType == "image/png" {
		image := PNGImageRespones{
			Data: data,
		}

		if err := image.Save(url); err != nil {
			return fmt.Errorf("failed to save png image: %s", err.Error())
		}
	} else if contentType == "application/pdf" {
		document := PDFResponse{
			Data: data,
		}

		if err := document.Save(url); err != nil {
			return fmt.Errorf("failed to save pdf: %s", err.Error())
		}
	} else if contentType == "video/mp4" {
		video := MP4VideoResponse{
			Data: data,
		}

		if err := video.Save(url); err != nil {
			return fmt.Errorf("failed to save video: %s", err.Error())
		}
	} else {
		return fmt.Errorf("invalid file format : %s\n", contentType)
	}

	return nil
}

type PDFResponse struct {
	Data []byte
}

func (p *PDFResponse) Save(url string) error {
	if err := os.Mkdir("Documents", 0755); err != nil {
		log.Printf("pdf mkdir failed: %s", err.Error())
	}

	filename := path.Base(url)

	if filename == "." || filename == "/" || !strings.Contains(filename, ".pdf") {
		filename = fmt.Sprintf("%d.pdf", time.Now().Unix())
	}

	err := ioutil.WriteFile(fmt.Sprintf("Documents/%s", filename), p.Data, 0755)
	if err != nil {
		return err
	}

	return nil
}

type JPEGImageRespones struct {
	Data []byte
}

func (p *JPEGImageRespones) Save(url string) error {
	// we should actually first chech if directory exists
	// instead I just ignore error altogether.
	if err := os.Mkdir("Photos", 0755); err != nil {
		log.Printf("photo mkdir failed: %s", err.Error())
	}

	filename := path.Base(url)

	if filename == "." || filename == "/" || !strings.Contains(filename, ".jpg") || strings.Contains(filename, ".jpeg") {
		filename = fmt.Sprintf("%d.jpg", time.Now().Unix())
	}

	err := ioutil.WriteFile(fmt.Sprintf("Photos/%s", filename), p.Data, 0755)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

type PNGImageRespones struct {
	Data []byte
}

func (p *PNGImageRespones) Save(url string) error {
	// we should actually first chech if directory exists
	// instead I just ignore error altogether.
	if err := os.Mkdir("Photos", 0755); err != nil {
		log.Printf("photo mkdir failed: %s", err.Error())
	}

	filename := path.Base(url)

	if filename == "." || filename == "/" || !strings.Contains(filename, ".png") {
		filename = fmt.Sprintf("%d.png", time.Now().Unix())
	}

	err := ioutil.WriteFile(fmt.Sprintf("Photos/%s", filename), p.Data, 0755)
	if err != nil {
		return err
	}

	return nil
}

type MP4VideoResponse struct {
	Data []byte
}

func (p *MP4VideoResponse) Save(url string) error {
	if err := os.Mkdir("Videos", 0755); err != nil {
		log.Printf("video mkdir failed: %s", err.Error())
	}

	filename := path.Base(url)

	if filename == "." || filename == "/" || !strings.Contains(filename, ".mp4") {
		filename = fmt.Sprintf("%d.mp4", time.Now().Unix())
	}

	err := ioutil.WriteFile(fmt.Sprintf("Videos/%s", filename), p.Data, 0755)
	if err != nil {
		return err
	}

	return nil
}
