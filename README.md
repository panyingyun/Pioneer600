# Pioneer600

Learn hardware by golang on Raspberry Pi 4

### 1、Hardware 
- Raspberry Pi 4 Model B Rev 1.2
- Pioneer600 (http://www.waveshare.net/wiki/Pioneer600)

### 2、Rom & OS 
- Raspberry Pi OS (32-bit) Lite (Raspbian GNU/Linux 10)

### 3、Boot config file

*open i2c、spi、1-wire interface*

*config.txt*
```shell
# Uncomment some or all of these to enable the optional hardware interfaces
dtparam=i2c_arm=on
#dtparam=i2s=on
dtparam=spi=on

# Additional overlays and parameters are documented /boot/overlays/README

# Enable audio (loads snd_bcm2835)
dtparam=audio=on

[pi4]
# Enable DRM VC4 V3D driver on top of the dispmanx display stack
dtoverlay=vc4-fkms-v3d
max_framebuffers=2

[all]
#dtoverlay=vc4-fkms-v3d
dtoverlay=w1-gpio,gpio_pin=4
```

### 4、How to build 
```shell

//ssh to your board and prepare Pioneer600
//and install git tools 
sudo apt-get update
sudo apt-get upgrade
sudo apt install -y vim git build-essential 
sudo apt install -y golang

//clone code 
git clone git@github.com:panyingyun/Pioneer600.git

//build 
cd Pioneer600/src/pi/cmd

//build arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build  -o Pioneer600
```

### 5、How to run 

- Test Gpio LED (D1)
```shell
sudo ./Pioneer600 -f 0
```

- I2C Test PCF8574 LED(D2)
```shell
sudo ./Pioneer600 -f 1
```

- I2C Test PCF8574 Beep
```shell
sudo ./Pioneer600 -f 2
```

- 1-Wire Test DS18b20.
```shell
sudo ./Pioneer600 -f 3
```

- I2C RTC Test DS3231.
```shell
sudo ./Pioneer600 -f 4
```

- "SPI Test SSD1306.
```shell
sudo ./Pioneer600 -f 5
```

### 6、Auto run when reboot OS 

Run build_arm64.sh, it will autorun Pioneer600 when reboot os
```shell
chmod 755 build_arm64.sh 
sudo sh build_arm64.sh 
```

Other shell cmd 
```shell
// enable 
sudo systemctl enable Pioneer600.service

// disable 
sudo systemctl disable Pioneer600.service

// start service 
sudo systemctl start  Pioneer600.service

// stop 
sudo systemctl stop Pioneer600.service

// restart 
sudo systemctl restart Pioneer600.service
```


### 7、Thanks To

```shell
 	logger.Info("Raspberry Pi 4 and Pioneer600")
	logger.Info("Learn how to use golang control devices.")
	logger.Info("Thanks to https://gobot.io")
	logger.Info("Thanks to http://www.waveshare.net.")
	logger.Info("Thanks to https://periph.io/project/goals.")
```