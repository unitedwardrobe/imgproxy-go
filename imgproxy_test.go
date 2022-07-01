package imgproxy

import (
	"encoding/hex"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewImgproxy(t *testing.T) {
	Convey("NewImgproxy()", t, func() {
		Convey("Retunrs error if the signature is not valid", func() {
			_, err := NewImgproxy(Config{
				BaseURL:       "http://localhost",
				SignatureSize: 33,
				Key:           hex.EncodeToString([]byte("key")),
				Salt:          hex.EncodeToString([]byte("salt")),
				Encode:        true,
			})
			So(errors.Cause(err), ShouldResemble, ErrInvalidSignature)
		})
	})
}

func Test_ImgproxyBuilder(t *testing.T) {
	Convey("Imgproxy.Builder()", t, func() {
		Convey("Returns the url with the uri encoded and sign when Encode is true and key and salt are not empty", func() {
			ip, err := NewImgproxy(Config{
				BaseURL:       "http://localhost",
				SignatureSize: 15,
				Key:           hex.EncodeToString([]byte("key")),
				Salt:          hex.EncodeToString([]byte("salt")),
				Encode:        true,
			})
			So(err, ShouldBeNil)

			url, err := ip.Builder().Generate("my/image.jpg")
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "http://localhost/6wIzqvuZtfHT1LL3J_z0/bXkvaW1hZ2UuanBn")
		})

		Convey("Returns the url without signature when key and salt are empty", func() {
			ip, err := NewImgproxy(Config{
				BaseURL:       "http://localhost",
				SignatureSize: 15,
				Key:           "",
				Salt:          "",
				Encode:        false,
			})
			So(err, ShouldBeNil)

			url, err := ip.Builder().Generate("my/image.jpg")
			So(err, ShouldBeNil)
			So(url, ShouldEqual, "http://localhost/insecure/plain/my/image.jpg")
		})

		Convey("With key salt and no encoded", func() {
			ip, err := NewImgproxy(Config{
				BaseURL:       "http://localhost",
				SignatureSize: 15,
				Key:           hex.EncodeToString([]byte("key")),
				Salt:          hex.EncodeToString([]byte("salt")),
				Encode:        false,
			})

			So(err, ShouldBeNil)

			Convey("With resize", func() {
				Convey("Sets fit option", func() {
					url, err := ip.Builder().
						Resize(ResizingTypeFit, 123, 456, true, false).
						Generate("my/image.jpg")

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "http://localhost/GmhHRFFOqT3jETNsx2ag/rs:fit:123:456:1/plain/my/image.jpg")
				})

				Convey("Sets fill option", func() {
					url, err := ip.Builder().
						Resize(ResizingTypeFill, 123, 456, true, false).
						Generate("my/image.jpg")

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "http://localhost/ya-KQkCvu-ZWLI2hZ3IT/rs:fill:123:456:1/plain/my/image.jpg")
				})
			})

			Convey("Size sets size option", func() {
				url, err := ip.Builder().
					Size(1, 2, true).
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/QvjH30YVgilqJVE3wjlj/s:1:2:1/plain/my/image.jpg")
			})

			Convey("ResizingType sets resizing type option", func() {
				url, err := ip.Builder().
					ResizingType(ResizingTypeFill).
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/YohdbDjMgBhATlRO7Ifs/rs:fill/plain/my/image.jpg")
			})

			Convey("Width sets width option", func() {
				url, err := ip.Builder().
					Width(1).
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/196LdHe9OIT7BZBGvnHF/w:1/plain/my/image.jpg")
			})

			Convey("Height sets height option", func() {
				url, err := ip.Builder().
					Height(1).
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/jR-ekA_-dQRu4KZ4ULEC/h:1/plain/my/image.jpg")
			})

			Convey("DPR", func() {
				Convey("With zero skips option", func() {
					url, err := ip.Builder().
						DPR(0).
						Generate("my/image.jpg")

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "http://localhost/s-cFqOcqN4HMtEZQwoyp/plain/my/image.jpg")
				})

				Convey("Higher than zero sets the option", func() {
					url, err := ip.Builder().
						DPR(10).
						Generate("my/image.jpg")

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "http://localhost/H3-NasL_V_EQt3f84ocr/dpr:10/plain/my/image.jpg")
				})
			})

			Convey("Enlarge sets enlarge option", func() {
				url, err := ip.Builder().
					Enlarge(1).
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/Kh3OA5md8aQEj5l9oc4t/el:1/plain/my/image.jpg")
			})

			Convey("Gravity", func() {
				Convey("With GravityEnum.* it sets the option", func() {
					url, err := ip.Builder().
						Gravity(GravityEnumCenter).
						Generate("my/image.jpg")

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "http://localhost/W0gScceNWUH3OUnYcHWI/g:ce/plain/my/image.jpg")
				})

				Convey("With OffsetGravity it sets the option", func() {
					url, err := ip.Builder().
						Gravity(OffsetGravity{
							Type:    GravityEnumNorth,
							XOffset: 10,
							YOffset: 20,
						}).
						Generate("my/image.jpg")

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "http://localhost/Y4L9ShmQTBTGIIINzZ3S/g:no:10:20/plain/my/image.jpg")
				})

				Convey("With FocusPoint it sets the option", func() {
					url, err := ip.Builder().
						Gravity(FocusPoint{
							X: 10,
							Y: 20,
						}).
						Generate("my/image.jpg")

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "http://localhost/tzY1UfBRkno8WSwTFnsN/g:fp:10:20/plain/my/image.jpg")
				})
			})

			Convey("Quality sets the quality option", func() {
				url, err := ip.Builder().
					Quality(10).
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/MBiaiGY_V7KY3L0IU4OO/q:10/plain/my/image.jpg")
			})

			Convey("Background", func() {
				Convey("With HexColor sets the option", func() {
					url, err := ip.Builder().
						Background(HexColor("#000000")).
						Generate("my/image.jpg")

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "http://localhost/s_b6BgKbDAPpika-5B-H/bg:#000000/plain/my/image.jpg")
				})

				Convey("With RGBColor sets the option", func() {
					url, err := ip.Builder().
						Background(RGBColor{
							R: 1,
							G: 2,
							B: 3,
						}).
						Generate("my/image.jpg")

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "http://localhost/cZhqkP4TlQRjki_sH00q/bg:1:2:3/plain/my/image.jpg")
				})
			})

			Convey("Blur sets the blur option", func() {
				url, err := ip.Builder().
					Blur(10).
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/5FzJ1UlJmL4Sy47brR4Q/bl:10/plain/my/image.jpg")
			})

			Convey("Sharpen sets the sharpen option", func() {
				url, err := ip.Builder().
					Sharpen(10).
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/gs5haFXa0LuNuFr7K88J/sh:10/plain/my/image.jpg")
			})

			Convey("Watermark", func() {
				Convey("With offset sets the option", func() {
					url, err := ip.Builder().
						Watermark(
							1,
							WatermarkPositionWest,
							&WatermarkOffset{
								X: 1,
								Y: 2,
							},
							3,
						).
						Generate("my/image.jpg")

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "http://localhost/2DBmLVDEF_HSDTlksj1c/wm:1:we:1:2:3/plain/my/image.jpg")
				})

				Convey("Without offset sets the option", func() {
					url, err := ip.Builder().
						Watermark(
							1,
							WatermarkPositionWest,
							nil,
							3,
						).
						Generate("my/image.jpg")

					So(err, ShouldBeNil)
					So(url, ShouldEqual, "http://localhost/Kj5PQr1LcllLJp39EZhf/wm:1:we:3/plain/my/image.jpg")
				})
			})

			Convey("Preset sets the preset option", func() {
				url, err := ip.Builder().
					Preset("foo", "bar").
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/16Ns_TrC2dRfO8dxVohb/pr:foo:bar/plain/my/image.jpg")
			})

			Convey("CacheBuster sets the cacheBuster option", func() {
				url, err := ip.Builder().
					CacheBuster("foo").
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/6FP7ES0ITuA5lFKCUjV4/cb:foo/plain/my/image.jpg")
			})

			Convey("Format sets the format option", func() {
				url, err := ip.Builder().
					Format("png").
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/jXuXqfAktdBIyinMAcf8/f:png/plain/my/image.jpg")
			})

			Convey("Crop sets the crop option", func() {
				url, err := ip.Builder().
					Crop(1, 2, GravityEnumCenter).
					Generate("my/image.jpg")

				So(err, ShouldBeNil)
				So(url, ShouldEqual, "http://localhost/00J_9T9UyVpOBQkQbodf/c:1:2:ce/plain/my/image.jpg")
			})
		})
	})
}
