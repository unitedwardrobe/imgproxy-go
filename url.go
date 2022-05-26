package imgproxy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

// ImgproxyURLData is a struct that contains the data required for generating an imgproxy URL.
type ImgproxyURLData struct {
	*Imgproxy
	Options map[string]string
}

// Generate generates the imgproxy URL.
func (i *ImgproxyURLData) Generate(uri string) (string, error) {
	if i.cfg.Encode {
		uri = base64.RawStdEncoding.EncodeToString([]byte(uri))
	} else {
		uri = "plain/" + uri
	}

	keys := make([]string, len(i.Options))
	j := 0
	for key := range i.Options {
		keys[j] = key
		j++
	}
	sort.Strings(keys)

	options := "/"
	for _, key := range keys {
		options += key + ":" + i.Options[key] + "/"
	}

	uriWithOptions := options + uri
	signature := hmac.New(sha256.New, i.key)

	if _, err := signature.Write(i.salt); err != nil {
		return "", errors.WithStack(err)
	}

	if _, err := signature.Write([]byte(uriWithOptions)); err != nil {
		return "", errors.WithStack(err)
	}

	sha := base64.RawURLEncoding.EncodeToString(signature.Sum(nil)[:i.cfg.SignatureSize])

	return i.cfg.BaseURL + sha + uriWithOptions, nil
}

// ResizingType enum.
type ResizingType string

// ResizingType enum.
const (
	ResizingTypeFill = "fill"
	ResizingTypeFit  = "fit"
	ResizingTypeCrop = "crop"
)

// Resize resizes the image.
func (i *ImgproxyURLData) Resize(resizingType ResizingType, width int, height int, enlarge bool) *ImgproxyURLData {
	return i.SetOption("rs", fmt.Sprintf(
		"%s:%d:%d:%s",
		resizingType,
		width, height,
		boolAsNumberString(enlarge),
	))
}

// SetOption sets an option on the URL.
func (i *ImgproxyURLData) SetOption(key, value string) *ImgproxyURLData {
	i.Options[key] = value
	return i
}
