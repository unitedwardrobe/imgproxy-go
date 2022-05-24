package imgproxy

import (
	"encoding/hex"
	"strings"

	"github.com/pkg/errors"
)

// Imgproxy is a URL builder helper for imgproxy.
type Imgproxy struct {
	cfg  Config
	key  []byte
	salt []byte
}

// NewImgproxy returns a new *Imgproxy.
func NewImgproxy(cfg Config) (*Imgproxy, error) {
	if !strings.HasSuffix(cfg.BaseURL, "/") {
		cfg.BaseURL = cfg.BaseURL + "/"
	}

	key, err := hex.DecodeString(cfg.Key)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	salt, err := hex.DecodeString(cfg.Salt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Imgproxy{
		cfg:  cfg,
		salt: salt,
		key:  key,
	}, nil
}

// Builder returns a *ImgproxyURLData that can be used to construct an imgproxy URL.
func (i *Imgproxy) Builder() *ImgproxyURLData {
	return &ImgproxyURLData{
		Imgproxy: i,
		Options:  make(map[string]string, 0),
	}
}
