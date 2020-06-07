package dev

import (
	"pi/driver"
	"pi/log"
)

const (
	pinLedOne int = 26

	statusOFFLedOne int = 0
	statusONLedOne  int = 1
)

type LEDOne struct {
	pin    *driver.DigitalPin
	status int
}

func NewLEDOne() *LEDOne {
	return &LEDOne{
		pin:    driver.NewDigitalPin(pinLedOne),
		status: statusOFFLedOne,
	}
}

func (led *LEDOne) Init() error {
	log.Default().Info("Init LED1 here.!")
	err := led.pin.Export()
	if err != nil {
		return err
	}
	err = led.pin.Direction(driver.OUT)
	if err != nil {
		return err
	}
	return nil
}

//低电平亮，高电平暗
func (led *LEDOne) On() error {
	log.Default().Info("Switch LED to On.!")
	err := led.pin.Write(driver.LOW)
	if err != nil {
		return err
	}
	led.status = statusONLedOne
	return nil
}

func (led *LEDOne) Off() error {
	log.Default().Info("Switch LED to Off.!")
	err := led.pin.Write(driver.HIGH)
	if err != nil {
		return err
	}
	led.status = statusOFFLedOne
	return nil
}

func (led *LEDOne) Status() int {
	return led.status
}

func (led *LEDOne) Toggle() error {
	if led.status == statusOFFLedOne {
		return led.On()
	} else {
		return led.Off()
	}
}
