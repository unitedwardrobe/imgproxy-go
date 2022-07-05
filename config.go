package imgproxy

// Config holds the parameters for constructing an imgproxy URL builder.
type Config struct {
	BaseURL       string
	SignatureSize int
	Key           string
	Salt          string
	EncodePath    bool
}
