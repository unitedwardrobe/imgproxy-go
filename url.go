package imgproxy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ImgproxyURLData is a struct that contains the data required for generating an imgproxy URL.
type ImgproxyURLData struct {
	*Imgproxy
	Options map[string]string
}

const insecureSignature = "insecure"

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

	if len(i.salt) == 0 && len(i.key) == 0 {
		return i.cfg.BaseURL + insecureSignature + uriWithOptions, nil
	}

	signature, err := getSignatureHash(i.key, i.salt, i.cfg.SignatureSize, uriWithOptions)
	if err != nil {
		return "", err
	}

	return i.cfg.BaseURL + signature + uriWithOptions, nil
}

func getSignatureHash(key []byte, salt []byte, signatureSize int, payload string) (string, error) {
	signature := hmac.New(sha256.New, key)

	if _, err := signature.Write(salt); err != nil {
		return "", errors.WithStack(err)
	}

	if _, err := signature.Write([]byte(payload)); err != nil {
		return "", errors.WithStack(err)
	}

	sha := base64.RawURLEncoding.EncodeToString(signature.Sum(nil)[:signatureSize])

	return sha, nil
}

// resizingType enum.
type resizingType string

// ResizingType enum.
const (
	ResizingTypeFit      = resizingType("fit")       // resizes the image while keeping aspect ratio to fit a given size.
	ResizingTypeFill     = resizingType("fill")      // resizes the image while keeping aspect ratio to fill a given size and crops projecting parts.
	ResizingTypeFillDown = resizingType("fill-down") // the same as fill, but if the resized image is smaller than the requested size, imgproxy will crop the result to keep the requested aspect ratio.
	ResizingTypeForce    = resizingType("force")     // resizes the image without keeping the aspect ratio.
	ResizingTypeAuto     = resizingType("auto")      // if both source and resulting dimensions have the same orientation (portrait or landscape), imgproxy will use fill. Otherwise, it will use fit.
)

// Resize resizes the image.
func (i *ImgproxyURLData) Resize(resizingType resizingType, width int, height int, enlarge bool) *ImgproxyURLData {
	return i.SetOption("rs", fmt.Sprintf(
		"%s:%d:%d:%s",
		resizingType,
		width, height,
		boolAsNumberString(enlarge),
	))
}

func (i *ImgproxyURLData) Size(width int, height int, enlarge bool) *ImgproxyURLData {
	return i.SetOption("s", fmt.Sprintf(
		"%d:%d:%s",
		width, height,
		boolAsNumberString(enlarge),
	))
}

func (i *ImgproxyURLData) ResizingType(resizingType resizingType) *ImgproxyURLData {
	return i.SetOption("rs", string(resizingType))
}

// Width defines the width of the resulting image.
// When set to 0, imgproxy will calculate width using the defined height and source aspect ratio.
// When set to 0 and resizing type is force, imgproxy will keep the original width.
func (i *ImgproxyURLData) Width(width int) *ImgproxyURLData {
	return i.SetOption("w", strconv.Itoa(width))
}

// Height defines the height of the resulting image.
// When set to 0, imgproxy will calculate resulting height using the defined width and source aspect ratio.
// When set to 0 and resizing type is force, imgproxy will keep the original height.
func (i *ImgproxyURLData) Height(height int) *ImgproxyURLData {
	return i.SetOption("h", strconv.Itoa(height))
}

// DPR controls the output density of your image.
func (i *ImgproxyURLData) DPR(dpr int) *ImgproxyURLData {
	if dpr > 0 {
		return i.SetOption("dpr", strconv.Itoa(dpr))
	}

	return i
}

// Enlarge enlarges the image.
func (i *ImgproxyURLData) Enlarge(enlarge int) *ImgproxyURLData {
	return i.SetOption("el", strconv.Itoa(enlarge))

}

type GravitySetter interface {
	SetGravityOption(i *ImgproxyURLData) *ImgproxyURLData
	GetStringOption() string
}

type OffsetGravity struct {
	Type    gravityEnum
	XOffset int
	YOffset int
}

func (o OffsetGravity) SetGravityOption(i *ImgproxyURLData) *ImgproxyURLData {
	return i.SetOption("g", o.GetStringOption())
}

func (o OffsetGravity) GetStringOption() string {
	return fmt.Sprintf("%s:%d:%d", o.Type, o.XOffset, o.YOffset)
}

type FocusPoint struct {
	X int64
	Y int64
}

func (f FocusPoint) SetGravityOption(i *ImgproxyURLData) *ImgproxyURLData {
	return i.SetOption("g", f.GetStringOption())
}

func (f FocusPoint) GetStringOption() string {
	return fmt.Sprintf("fp:%d:%d", f.X, f.Y)
}

type gravityEnum string

const (
	GravityEnumCenter    = gravityEnum("ce")   // default
	GravityEnumNorth     = gravityEnum("no")   // (top edge)
	GravityEnumSouth     = gravityEnum("so")   // (bottom edge)
	GravityEnumEast      = gravityEnum("ea")   // (right edge)
	GravityEnumWest      = gravityEnum("we")   // (left edge)
	GravityEnumNorthEast = gravityEnum("noea") // (top-right corner)
	GravityEnumNorthWest = gravityEnum("nowe") // (top-left corner)
	GravityEnumSouthEast = gravityEnum("soea") // (bottom-right corner)
	GravityEnumSouthWest = gravityEnum("sowe") // (bottom-left corner)
	GravityEnumSmart     = gravityEnum("sm")   // libvips detects the most "interesting" section of the image and considers it as the center of the resulting image
)

func (g gravityEnum) SetGravityOption(i *ImgproxyURLData) *ImgproxyURLData {
	return i.SetOption("g", g.GetStringOption())
}

func (g gravityEnum) GetStringOption() string {
	return string(g)
}

// Gravity guides imgproxy when needs to cut some parts of the image.
func (i *ImgproxyURLData) Gravity(g GravitySetter) *ImgproxyURLData {
	return g.SetGravityOption(i)
}

// Quality redefines quality of the resulting image, as a percentage.
func (i *ImgproxyURLData) Quality(quality int) *ImgproxyURLData {
	return i.SetOption("q", strconv.Itoa(quality))
}

type HexColor string

func (h HexColor) SetBgOption(i *ImgproxyURLData) *ImgproxyURLData {
	return i.SetOption("bg", string(h))
}

type RGBColor struct {
	R int
	G int
	B int
}

func (rgb RGBColor) SetBgOption(i *ImgproxyURLData) *ImgproxyURLData {
	return i.SetOption("bg", fmt.Sprintf("%d:%d:%d", rgb.R, rgb.G, rgb.B))
}

type BackgroundSetter interface {
	SetBgOption(*ImgproxyURLData) *ImgproxyURLData
}

// Background fills the resulting image background with the specified color.
// RGBColor are the red, green and blue channel values of the background color (0-255).
// HexColor is a hex-coded value of the color.
// Useful when you convert an image with alpha-channel to JPEG.
func (i *ImgproxyURLData) Background(bg BackgroundSetter) *ImgproxyURLData {
	return bg.SetBgOption(i)
}

// Blur applies a gaussian blur filter to the resulting image.
// The value of sigma defines the size of the mask imgproxy will use.
func (i *ImgproxyURLData) Blur(sigma int) *ImgproxyURLData {
	return i.SetOption("bl", strconv.Itoa(sigma))
}

// Sharpen applies the sharpen filter to the resulting image.
// The value of sigma defines the size of the mask imgproxy will use.
func (i *ImgproxyURLData) Sharpen(sigma int) *ImgproxyURLData {
	return i.SetOption("sh", strconv.Itoa(sigma))
}

type watermarkPosition string

const (
	WatermarkPositionCenter    = watermarkPosition("ce")   // default
	WatermarkPositionNorth     = watermarkPosition("no")   // (top edge);
	WatermarkPositionSouth     = watermarkPosition("so")   // (bottom edge);
	WatermarkPositionEast      = watermarkPosition("ea")   // (right edge);
	WatermarkPositionWest      = watermarkPosition("we")   // (left edge);
	WatermarkPositionNorthEast = watermarkPosition("noea") // (top-right corner);
	WatermarkPositionNorthWest = watermarkPosition("nowe") // (top-left corner);
	WatermarkPositionSouthEast = watermarkPosition("soea") // (bottom-right corner);
	WatermarkPositionSouthWest = watermarkPosition("sowe") // (bottom-left corner);
	WatermarkPositionReplicate = watermarkPosition("re")   // replicate watermark to fill the whole image;
)

type WatermarkOffset struct {
	X int
	Y int
}

// Watermark places a watermark on the processed image.
func (i *ImgproxyURLData) Watermark(opacity int, position watermarkPosition, offset *WatermarkOffset, scale int) *ImgproxyURLData {
	var offsetStr string

	if offset != nil {
		offsetStr = fmt.Sprintf(":%d:%d", offset.X, offset.Y)
	}

	return i.SetOption("wm",
		fmt.Sprintf(
			"%d:%s%s:%d", opacity, position, offsetStr, scale,
		),
	)
}

// Preset defines a list of presets to be used by imgproxy.
func (i *ImgproxyURLData) Preset(presets ...string) *ImgproxyURLData {
	return i.SetOption("pr", strings.Join(presets, ":"))
}

// CacheBuster doesn’t affect image processing but its changing allows for bypassing the CDN, proxy server and browser cache.
// Useful when you have changed some things that are not reflected in the URL, like image quality settings, presets, or watermark data.
// It’s highly recommended to prefer the cachebuster option over a URL query string because that option can be properly signed.
func (i *ImgproxyURLData) CacheBuster(buster string) *ImgproxyURLData {
	return i.SetOption("cb", buster)
}

// Format specifies the resulting image format. Alias for the extension part of the URL.
func (i *ImgproxyURLData) Format(extension string) *ImgproxyURLData {
	return i.SetOption("f", extension)
}

func (i *ImgproxyURLData) Crop(width int, height int, gravity GravitySetter) *ImgproxyURLData {
	crop := fmt.Sprintf("%d:%d", width, height)

	if gravity != nil {
		crop += ":" + gravity.GetStringOption()
	}

	return i.SetOption("c", crop)
}

// SetOption sets an option on the URL.
func (i *ImgproxyURLData) SetOption(key, value string) *ImgproxyURLData {
	i.Options[key] = value
	return i
}
