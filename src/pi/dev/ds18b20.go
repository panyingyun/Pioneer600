package dev

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"pi/log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//http://www.waveshare.net/study/article-607-1.html
//sudo modprobe w1-gpio
//sudo modprobe w1-therm
//cd  /sys/bus/w1/devices
//cd 28-00000xxx
//cat w1_slave
//"13 02 4b 46 7f ff 0d 10 e7 : crc=e7 YES
// 13 02 4b 46 7f ff 0d 10 e7 t=33187"
// /boot/config.txt 添加 dtoverlay=w1-gpio,gpio_pin=4
// 打开单总线，对应BCM的4脚

const (
	rootPath          = "/sys/bus/w1/devices/"
	ds18b20PrefixName = "28-00"
	ds18b20Slave      = "/w1_slave"
)

type DS18B20 struct {
	name      string
	temperate float64
}

func NewDS18B20() *DS18B20 {
	return &DS18B20{
		name:      "",
		temperate: 0.0,
	}
}

// readDirNames reads the directory named by dirname and returns
// a sorted list of directory entries.
func readDirNames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

func (d *DS18B20) FetchTemperate() (err error) {
	// find ds18b20
	names, _ := readDirNames(rootPath)
	for _, name := range names {
		if strings.HasPrefix(name, ds18b20PrefixName) {
			d.name = name
			log.Default().Info("ds18b20's name:  ", name)
			break
		}
	}
	if d.name == "" {
		err = errors.New("Can not find ds18b20.")
		return
	}
	//calculate temperate
	devPath := rootPath + d.name + (string)(filepath.Separator) + ds18b20Slave
	log.Default().Infof("devPath = ", devPath)
	data, err := ioutil.ReadFile(devPath)
	if err != nil {
		return
	}
	r := regexp.MustCompile("t=([0-9]+)")
	a := r.FindString(string(data))
	b := strings.ReplaceAll(a, "t=", "")
	temp, _ := strconv.ParseFloat(b, 64)
	temp = temp / 1000.0
	d.temperate = temp
	return
}

func (d *DS18B20) Name() string {
	return d.name
}

func (d *DS18B20) Temperate() float64 {
	return d.temperate
}
