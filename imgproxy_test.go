package imgproxy

import (
	"encoding/hex"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Imgproxy(t *testing.T) {
	Convey("Imgproxy", t, func() {
		ip, err := NewImgproxy(Config{
			BaseURL:       "http://localhost",
			SignatureSize: 15,
			Key:           hex.EncodeToString([]byte{1, 2, 3, 4}),
			Salt:          hex.EncodeToString([]byte{5, 6, 7, 8}),
			Encode:        true,
		})
		So(err, ShouldBeNil)

		url, err := ip.Builder().Resize(ResizingTypeFit, 123, 456, true).Generate("my/image.jpg")
		So(err, ShouldBeNil)
		So(url, ShouldEqual, "http://localhost/ed574597a1570b1/rs:fit:123:456:1/bXkvaW1hZ2UuanBn")
	})
}
