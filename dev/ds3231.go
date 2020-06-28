package dev

import (
	"fmt"
	"pi/driver"
	"pi/log"
)

//http://www.waveshare.net/study/article-623-1.html
//http://www.waveshare.net/study/article-623-1.html
//https://shumeipai.nxez.com/2019/05/08/raspberry-pi-configuration-of-ds3231-clock-module-i2c-interface.html
const (
	//DS3231
	I2cAddrDS3231 int = 0x68
)

var (
	timeNow = []byte{0x00, 0x59, 0x23, 0x04, 0x04, 0x06, 0x20}
	weekArr = [7]string{"SUN", "Mon", "Tues", "Wed", "Thur", "Fri", "Sat"}
)

type DS3231 struct {
	i2c  *driver.I2CDevice
	time string
}

func NewDS3231() *DS3231 {
	dev, err := driver.NewI2cDevice(I2cDev)
	if err != nil {
		log.Default().Error("err: ", err)
		return nil
	}
	dev.SetAddress(I2cAddrDS3231)
	return &DS3231{
		i2c:  dev,
		time: "",
	}
}

func (d *DS3231) SetTime() error {
	return d.i2c.WriteBlockData(0x00, timeNow)
}

func (d *DS3231) Time() (time string, err error) {
	sec, err := d.i2c.ReadByteData(0x00)
	sec = sec & 0x7f
	min, err := d.i2c.ReadByteData(0x01)
	min = min & 0x7f
	hour, err := d.i2c.ReadByteData(0x02)
	hour = hour & 0x3F
	week, err := d.i2c.ReadByteData(0x03)
	week = week & 0x07
	day, err := d.i2c.ReadByteData(0x04)
	day = day & 0x3F
	month, err := d.i2c.ReadByteData(0x05)
	month = month & 0x1F
	year, err := d.i2c.ReadByteData(0x06)
	d.time = fmt.Sprintf("20%x年 %x 月 %x日 %v  %x:%x:%x\n", year, month, day, weekArr[week-1], hour, min, sec)
	time = d.time
	return
}
