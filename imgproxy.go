package imgproxy

import (
	"encoding/hex"
	stdErrs "errors"
	"strings"

	"github.com/pkg/errors"
)

// Imgproxy is a URL builder helper for imgproxy.
type Imgproxy struct {
	cfg  Config
	key  []byte
	salt []byte
}

// ErrInvalidSignature error.
var ErrInvalidSignature = stdErrs.New("invalid signature size")

// NewImgproxy returns a new *Imgproxy.
func NewImgproxy(cfg Config) (*Imgproxy, error) {
	if !strings.HasSuffix(cfg.BaseURL, "/") {
		cfg.BaseURL = cfg.BaseURL + "/"
	}

	if cfg.SignatureSize < 1 || cfg.SignatureSize > 32 {
		return nil, errors.WithStack(ErrInvalidSignature)
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
