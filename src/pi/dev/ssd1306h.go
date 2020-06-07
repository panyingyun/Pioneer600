package dev

import (
	"image"
	"image/draw"
	"pi/driver"
	"pi/log"
	"time"
	"unicode/utf8"

	"periph.io/x/periph/conn/display"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/devices/ssd1306/image1bit"
	"periph.io/x/periph/host"
)

type SSD1306Pos int

const (
	PosTopCenter SSD1306Pos = iota
	PosTopLeft
	PosTopRight
	PosBottomLeft
	PosBottomRight
	PosBottomCenter
)

type SSD1306H struct {
	spidev        string
	dc            string
	height        int
	width         int
	rotated       bool
	sequential    bool
	swapTopBottom bool
	dev           *ssd1306.Dev
	dcpin         gpio.PinOut
	rsPinNo       int
	rstDriver     *driver.DigitalPin
}

func NewSSD1306H() *SSD1306H {
	s := &SSD1306H{
		spidev:        "/dev/spidev0.0",
		dc:            "16",
		height:        64,
		width:         128,
		rotated:       false,
		sequential:    false,
		swapTopBottom: false,
		rsPinNo:       19,
	}
	_, err := host.Init()
	if err != nil {
		log.Default().Error("Host init Fail!")
		return nil
	}

	s.dcpin = gpioreg.ByName(s.dc)
	s.rstDriver = driver.NewDigitalPin(s.rsPinNo)
	s.rstDriver.Export()
	s.rstDriver.Direction(driver.OUT)
	s.Reset()

	opts := ssd1306.Opts{W: s.width, H: s.height, Rotated: s.rotated, Sequential: s.sequential, SwapTopBottom: s.swapTopBottom}
	c, err := spireg.Open(s.spidev)
	if err != nil {
		log.Default().Error("SPI Open Fail!")
		return nil
	}

	s.dev, err = ssd1306.NewSPI(c, s.dcpin, &opts)
	if err != nil {
		log.Default().Error("Connect to SPI Dev Fail!")
		return nil
	}
	return s
}

func (ssd *SSD1306H) DrawText(pos SSD1306Pos, text string) error {
	log.Default().Info("dev Bounds = ", ssd.dev.Bounds())
	src := image1bit.NewVerticalLSB(ssd.dev.Bounds())
	img := convert(ssd.dev, src)
	switch pos {
	case PosTopCenter:
		drawTextTopCenter(img, text)
	case PosTopLeft:
		drawTextTopLeft(img, text)
	case PosTopRight:
		drawTextTopRight(img, text)
	case PosBottomLeft:
		drawTextBottomLeft(img, text)
	case PosBottomRight:
		drawTextBottomRight(img, text)
	case PosBottomCenter:
		drawTextBottomCenter(img, text)
	default:

	}

	if err := ssd.dev.Draw(ssd.dev.Bounds(), img, image.Point{}); err != nil {
		log.Default().Error("Draw error!")
		return err
	}
	return nil
}

// Reset SSD1306H
func (ssd *SSD1306H) Reset() (err error) {
	ssd.rstDriver.Write(driver.HIGH)
	time.Sleep(50 * time.Millisecond)
	ssd.rstDriver.Write(driver.LOW)
	time.Sleep(50 * time.Millisecond)
	ssd.rstDriver.Write(driver.HIGH)
	return nil
}

func (ssd *SSD1306H) DrawImage() error {
	return nil
}

func convert(display display.Drawer, src image.Image) *image1bit.VerticalLSB {
	screenBounds := display.Bounds()
	size := screenBounds.Size()
	src = resize(src, size)
	img := image1bit.NewVerticalLSB(screenBounds)
	r := src.Bounds()
	r = r.Add(image.Point{X: (size.X - r.Max.X) / 2, Y: (size.Y - r.Max.Y) / 2})
	draw.Draw(img, r, src, image.Point{}, draw.Src)
	return img
}

func resize(src image.Image, size image.Point) *image.NRGBA {
	srcMax := src.Bounds().Max
	dst := image.NewNRGBA(image.Rectangle{Max: size})
	for y := 0; y < size.Y; y++ {
		sY := (y*srcMax.Y + size.Y/2) / size.Y
		for x := 0; x < size.X; x++ {
			dst.Set(x, y, src.At((x*srcMax.X+size.X/2)/size.X, sY))
		}
	}
	return dst
}

func drawTextBottomRight(img draw.Image, text string) {
	advance := utf8.RuneCountInString(text) * 7
	bounds := img.Bounds()
	if advance > bounds.Dx() {
		advance = 0
	} else {
		advance = bounds.Dx() - advance
	}
	drawText(img, image.Point{X: advance, Y: bounds.Dy() - 1 - 13}, text)
}

func drawTextBottomLeft(img draw.Image, text string) {
	bounds := img.Bounds()
	drawText(img, image.Point{X: 0, Y: bounds.Dy() - 1 - 13}, text)
}

func drawTextBottomCenter(img draw.Image, text string) {
	advance := utf8.RuneCountInString(text) * 7
	bounds := img.Bounds()
	if advance > bounds.Dx() {
		advance = 0
	} else {
		advance = (bounds.Dx() - advance) / 2
	}
	drawText(img, image.Point{X: advance, Y: bounds.Dy() - 1 - 13}, text)
}

func drawTextTopLeft(img draw.Image, text string) {
	drawText(img, image.Point{X: 0, Y: 0}, text)
}

func drawTextTopRight(img draw.Image, text string) {
	advance := utf8.RuneCountInString(text) * 7
	bounds := img.Bounds()
	if advance > bounds.Dx() {
		advance = 0
	} else {
		advance = bounds.Dx() - advance
	}
	drawText(img, image.Point{X: advance, Y: 0}, text)
}

func drawTextTopCenter(img draw.Image, text string) {
	advance := utf8.RuneCountInString(text) * 7
	bounds := img.Bounds()
	if advance > bounds.Dx() {
		advance = 0
	} else {
		advance = (bounds.Dx() - advance) / 2
	}
	drawText(img, image.Point{X: advance, Y: 0}, text)
}
