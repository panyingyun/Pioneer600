package dev

import (
	"pi/driver"
	"pi/log"
)

const (
	//PCF8574
	I2cAddrPcf8574 int    = 0x20
	I2cLedTwoOn    byte   = 0xEF
	I2cLedTwoOff   byte   = 0xFF
	I2cDev         string = "/dev/i2c-1"
)

const (
	StatusOffLedTwo int = 0
	StatusOnLedTwo  int = 1
)

type PCF8574LED struct {
	i2c          *driver.I2CDevice
	ledTwoStatus int
}

func NewPCF8574LED() *PCF8574LED {
	dev, err := driver.NewI2cDevice(I2cDev)
	if err != nil {
		log.Default().Infof("err: %v", err)
		return nil
	}
	dev.SetAddress(I2cAddrPcf8574)
	return &PCF8574LED{
		i2c:          dev,
		ledTwoStatus: StatusOffLedTwo,
	}
}

func (p *PCF8574LED) LED2On() error {
	log.Default().Info("LED2 On ...")
	err := p.i2c.WriteByte(I2cLedTwoOn)
	if err != nil {
		return err
	}
	p.ledTwoStatus = StatusOnLedTwo
	return nil
}

func (p *PCF8574LED) LED2Off() error {
	log.Default().Info("LED2 Off ...")
	err := p.i2c.WriteByte(I2cLedTwoOff)
	if err != nil {
		return err
	}
	p.ledTwoStatus = StatusOffLedTwo
	return nil
}

func (p *PCF8574LED) Toggle() error {
	log.Default().Info("LED2 Toggle ...")
	if p.ledTwoStatus == StatusOffLedTwo {
		return p.LED2On()
	} else {
		return p.LED2Off()
	}
}
