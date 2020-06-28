package dev

import (
	"errors"
	"fmt"
	"image"
	"pi/driver"
	"time"
)

//https://github.com/google/periph
//https://github.com/google/periph/blob/master/devices/ssd1306/doc.go
const (
	//spi0Dev            = "/dev/spidev0.0"
	spiDefaultBus      = 0
	spiDefaultChip     = 0
	spiDefaultMode     = 0
	spiDefaultMaxSpeed = 500000
	spiDefaultBits     = 8
	// default values
	ssd1306RstPin = 19 // for raspberry pi
	ssd1306DcPin  = 16 // for raspberry pi
	ssd1306Width  = 128
	ssd1306Height = 64

	ssd1306ExternalVcc  = false
	ssd1306SetStartLine = 0x40
	// fundamental commands
	ssd1306SetContrast          = 0x81
	ssd1306DisplayOnResumeToRAM = 0xA4
	ssd1306DisplayOnResume      = 0xA5
	ssd1306SetDisplayNormal     = 0xA6
	ssd1306SetDisplayInverse    = 0xA7
	ssd1306SetDisplayOff        = 0xAE
	ssd1306SetDisplayOn         = 0xAF
	// scrolling commands
	ssd1306RightHorizontalScroll            = 0x26
	ssd1306LeftHorizontalScroll             = 0x27
	ssd1306VerticalAndRightHorizontalScroll = 0x29
	ssd1306VerticalAndLeftHorizontalScroll  = 0x2A
	ssd1306DeactivateScroll                 = 0x2E
	ssd1306ActivateScroll                   = 0x2F
	ssd1306SetVerticalScrollArea            = 0xA3
	// addressing settings commands
	ssd1306SetMemoryAddressingMode = 0x20
	ssd1306ColumnAddr              = 0x21
	ssd1306PageAddr                = 0x22
	// hardware configuration commands
	ssd1306SetSegmentRemap0   = 0xA0
	ssd1306SetSegmentRemap127 = 0xA1
	ssd1306SetMultiplexRatio  = 0xA8
	ssd1306ComScanInc         = 0xC0
	ssd1306ComScanDec         = 0xC8
	ssd1306SetDisplayOffset   = 0xD3
	ssd1306SetComPins         = 0xDA
	// timing and driving scheme commands
	ssd1306SetDisplayClock      = 0xD5
	ssd1306SetPrechargePeriod   = 0xD9
	ssd1306SetVComDeselectLevel = 0xDB
	ssd1306NOOP                 = 0xE3
	// charge pump command
	ssd1306ChargePumpSetting = 0x8D
)

// DisplayBuffer represents the display buffer intermediate memory
type DisplayBuffer struct {
	width, height, pageSize int
	buffer                  []byte
}

// NewDisplayBuffer creates a new DisplayBuffer
func NewDisplayBuffer(width, height, pageSize int) *DisplayBuffer {
	s := &DisplayBuffer{
		width:    width,
		height:   height,
		pageSize: pageSize,
	}
	s.buffer = make([]byte, s.Size())
	return s
}

// Size returns the memory size of the display buffer
func (d *DisplayBuffer) Size() int {
	return (d.width * d.height) / d.pageSize
}

// Clear the contents of the display buffer
func (d *DisplayBuffer) Clear() {
	d.buffer = make([]byte, d.Size())
}

// SetPixel sets the x, y pixel with c color
func (d *DisplayBuffer) SetPixel(x, y, c int) {
	idx := x + (y/d.pageSize)*d.width
	bit := uint(y) % uint(d.pageSize)
	if c == 0 {
		d.buffer[idx] &= ^(1 << bit)
	} else {
		d.buffer[idx] |= (1 << bit)
	}
}

// Set sets the display buffer with the given buffer
func (d *DisplayBuffer) Set(buf []byte) {
	d.buffer = buf
}

type SSD1306 struct {
	connection    driver.Connection
	name          string
	dcDriver      *driver.DigitalPin
	rstDriver     *driver.DigitalPin
	pageSize      int
	DisplayWidth  int
	DisplayHeight int
	DCPin         int
	RSTPin        int
	ExternalVcc   bool
	buffer        *DisplayBuffer
}

// NewSSD1306 creates a new NewSSD1306.
func NewSSD1306() *SSD1306 {
	// cast adaptor to spi connector since we also need the adaptor for gpio
	c, err := driver.GetSpiConnection(
		spiDefaultBus,
		spiDefaultChip,
		spiDefaultMode,
		spiDefaultBits,
		spiDefaultMaxSpeed)
	if err != nil {
		panic("unable to get connector for ssd1306")
	}
	fmt.Println("SpiConnection = ", c)
	s := &SSD1306{
		name:          "SSD1306",
		connection:    c,
		DisplayWidth:  ssd1306Width,
		DisplayHeight: ssd1306Height,
		DCPin:         ssd1306DcPin,
		RSTPin:        ssd1306RstPin,
		ExternalVcc:   ssd1306ExternalVcc,
	}
	s.dcDriver = driver.NewDigitalPin(s.DCPin)
	s.dcDriver.Export()
	s.dcDriver.Direction(driver.OUT)
	s.rstDriver = driver.NewDigitalPin(s.RSTPin)
	s.rstDriver.Export()
	s.rstDriver.Direction(driver.OUT)
	s.pageSize = s.DisplayHeight / 8
	s.buffer = NewDisplayBuffer(s.DisplayWidth, s.DisplayHeight, s.pageSize)
	s.Reset()
	s.ssd1306Init()
	return s
}

func (s *SSD1306) ssd1306Init() {
	s.command(ssd1306SetDisplayOff)
	s.command(ssd1306SetDisplayClock)
	if s.DisplayHeight == 16 {
		s.command(0x60)
	} else {
		s.command(0x80)
	}
	s.command(ssd1306SetMultiplexRatio)
	s.command(uint8(s.DisplayHeight) - 1)
	s.command(ssd1306SetDisplayOffset)
	s.command(0x0)
	s.command(ssd1306SetStartLine)
	s.command(0x0)
	s.command(ssd1306ChargePumpSetting)
	if s.ExternalVcc {
		s.command(0x10)
	} else {
		s.command(0x14)
	}
	s.command(ssd1306SetMemoryAddressingMode)
	s.command(0x00)
	s.command(ssd1306SetSegmentRemap0)
	s.command(0x01)
	s.command(ssd1306ComScanInc)
	s.command(ssd1306SetComPins)
	if s.DisplayHeight == 64 {
		s.command(0x12)
	} else {
		s.command(0x02)
	}
	s.command(ssd1306SetContrast)
	if s.DisplayHeight == 64 {
		if s.ExternalVcc {
			s.command(0x9F)
		} else {
			s.command(0xCF)
		}
	} else {
		s.command(0x8F)
	}
	s.command(ssd1306SetPrechargePeriod)
	if s.ExternalVcc {
		s.command(0x22)
	} else {
		s.command(0xF1)
	}
	s.command(ssd1306SetVComDeselectLevel)
	s.command(0x40)
	s.command(ssd1306DisplayOnResumeToRAM)
	s.command(ssd1306SetDisplayNormal)
	s.command(ssd1306DeactivateScroll)
	s.command(ssd1306SetDisplayOn)
}

// Halt returns true if device is halted successfully.
func (s *SSD1306) Halt() (err error) {
	s.Reset()
	s.Off()
	return nil
}

// On turns on the display.
func (s *SSD1306) On() (err error) {
	return s.command(ssd1306SetDisplayOn)
}

// Off turns off the display.
func (s *SSD1306) Off() (err error) {
	return s.command(ssd1306SetDisplayOff)
}

// Clear clears the display buffer.
func (s *SSD1306) Clear() (err error) {
	s.buffer.Clear()
	return nil
}

// Set sets a pixel in the display buffer.
func (s *SSD1306) Set(x, y, c int) {
	s.buffer.SetPixel(x, y, c)
}

// Reset re-initializes the device to a clean state.
func (s *SSD1306) Reset() (err error) {
	s.rstDriver.Write(driver.HIGH)
	time.Sleep(50 * time.Millisecond)
	s.rstDriver.Write(driver.LOW)
	time.Sleep(50 * time.Millisecond)
	s.rstDriver.Write(driver.HIGH)
	return nil
}

// SetBufferAndDisplay sets the display buffer with the given buffer and displays the image.
func (s *SSD1306) SetBufferAndDisplay(buf []byte) (err error) {
	s.buffer.Set(buf)
	return s.Display()
}

// SetContrast sets the display contrast (0-255).
func (s *SSD1306) SetContrast(contrast byte) (err error) {
	if contrast < 0 || contrast > 255 {
		return fmt.Errorf("contrast value must be between 0-255")
	}
	if err = s.command(ssd1306SetContrast); err != nil {
		return err
	}
	return s.command(contrast)
}

// Display sends the memory buffer to the display.
func (s *SSD1306) Display() (err error) {
	s.command(ssd1306ColumnAddr)
	s.command(0)
	s.command(uint8(s.DisplayWidth) - 1)
	s.command(ssd1306PageAddr)
	s.command(0)
	s.command(uint8(s.pageSize) - 1)
	if err = s.dcDriver.Write(driver.HIGH); err != nil {
		return err
	}
	return s.connection.Tx(append([]byte{0x40}, s.buffer.buffer...), nil)
}

// ShowImage takes a standard Go image and shows it on the display in monochrome.
func (s *SSD1306) ShowImage(img image.Image) (err error) {
	if img.Bounds().Dx() != s.DisplayWidth || img.Bounds().Dy() != s.DisplayHeight {
		return errors.New("Image must match the display width and height")
	}

	s.Clear()
	for y, w, h := 0, img.Bounds().Dx(), img.Bounds().Dy(); y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.At(x, y)
			if r, g, b, _ := c.RGBA(); r > 0 || g > 0 || b > 0 {
				s.Set(x, y, 1)
			}
		}
	}
	return s.Display()
}

// command sends a unique command
func (s *SSD1306) command(b byte) (err error) {
	if err = s.dcDriver.Write(driver.LOW); err != nil {
		return err
	}
	err = s.connection.Tx([]byte{b}, nil)
	return err
}
