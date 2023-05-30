package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

var (
	url        string
	headerSize int64 = 512
	imageTypes       = []string{
		"image/bmp",
		"image/cis-cod",
		"image/gif",
		"image/ief",
		"image/jpeg",
		"image/pipeg",
		"image/png",
		"image/svg+xml",
		"image/tiff",
		"image/webp",
		"image/x-cmu-raster",
		"image/x-cmx",
		"image/x-icon",
	}
	imageTypeMap = make(map[string]struct{})
)

func main() {

	for _, imageType := range imageTypes {
		imageTypeMap[imageType] = struct{}{}
	}

	flag.StringVar(&url, "url", "", "url to check")
	flag.Parse()
	imageOk, kind, err := isImage(url)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	if imageOk {
		fmt.Printf("IMAGE: ")
	} else {
		fmt.Printf("NOT IMAGE: ")
	}
	fmt.Printf("file type at %s is %s\n", url, kind)

}

func fileType(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	r := http.MaxBytesReader(nil, resp.Body, headerSize)
	buf, err := io.ReadAll(r)
	if err != nil {
		if _, ok := err.(*http.MaxBytesError); !ok {
			return "", err
		}
	}
	return http.DetectContentType(buf), nil
}

func isImage(url string) (bool, string, error) {
	kind, err := fileType(url)
	if err != nil {
		return false, "", err
	}
	if _, ok := imageTypeMap[kind]; ok {
		return true, kind, nil
	}
	return false, kind, fmt.Errorf("%s does not seem like an image")
}
