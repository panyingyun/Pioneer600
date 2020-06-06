package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pi/dev"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

const (
	FunctionGpioLedOne    int = 0
	FunctionPCF8574LedTwo int = 1
	FunctionPCF8574Beep   int = 2
	FunctionDs18B20       int = 3
	FunctionDs3231        int = 4
	FunctionSSD1306       int = 5
)

func testGpioLEDOne() {
	fmt.Println("GPIO Test LED One.")
	led := dev.NewLEDOne()
	led.Init()
	for {
		led.Toggle()
		time.Sleep(1 * time.Second)
	}
}

func testPCF8574LedTwo() {
	fmt.Println("I2C Test PCF8574 LED Two.")
	led2 := dev.NewPCF8574LED()
	for {
		led2.Toggle()
		time.Sleep(1 * time.Second)
	}
}

func testPCF8574Beep() {
	fmt.Println("I2C Test PCF8574 Beep.")
	beep := dev.NewPCF8574Beep()
	type note struct {
		tone     float64
		duration float64
	}

	song := []note{
		{dev.C4, dev.Quarter},
		{dev.C4, dev.Quarter},
		{dev.G4, dev.Quarter},
		{dev.G4, dev.Quarter},
		{dev.A4, dev.Quarter},
		{dev.A4, dev.Quarter},
		{dev.G4, dev.Half},
		{dev.F4, dev.Quarter},
		{dev.F4, dev.Quarter},
		{dev.E4, dev.Quarter},
		{dev.E4, dev.Quarter},
		{dev.D4, dev.Quarter},
		{dev.D4, dev.Quarter},
		{dev.C4, dev.Half},
	}
	for {
		for _, val := range song {
			beep.Tone(val.tone, val.duration)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func testDS18b20() {
	fmt.Println("1-Wire Test DS18b20.")
	ds18b20 := dev.NewDS18B20()
	for {
		err := ds18b20.FetchTemperate()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Current temperate : %v ℃\n", ds18b20.Temperate())
		time.Sleep(2 * time.Second)
	}
}

func testDS3231() {
	fmt.Println("I2C RTC　Test DS3231.")
	ds3231 := dev.NewDS3231()
	ds3231.SetTime()
	for {
		t, err := ds3231.Time()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(t)
		time.Sleep(2 * time.Second)
	}
}

func testSSD1306() {
	fmt.Println("SPI Test SSD1306.")
	ssd1306 := dev.NewSSD1306H()
	for {
		fmt.Println("ssd1306 ...... ", ssd1306)
		ssd1306.DrawText(dev.PosTopLeft, "Super Google.")
		time.Sleep(1 * time.Second)
		ssd1306.DrawText(dev.PosTopCenter, "Super Google.")
		time.Sleep(1 * time.Second)
		ssd1306.DrawText(dev.PosTopRight, "Super Google.")
		time.Sleep(1 * time.Second)
		ssd1306.DrawText(dev.PosBottomLeft, "Super Google.")
		time.Sleep(1 * time.Second)
		ssd1306.DrawText(dev.PosBottomCenter, "Super Google.")
		time.Sleep(1 * time.Second)
		ssd1306.DrawText(dev.PosBottomRight, "Super Google.")
		time.Sleep(1 * time.Second)
	}
}

func run(c *cli.Context) error {

	fmt.Println("conf = ", c.String("conf"))
	config := viper.New()
	config.SetConfigFile(c.String("conf"))
	config.SetConfigType("yaml")
	config.ReadInConfig()
	opt, err := NewOptions(config)
	if err != nil {
		fmt.Println("err = ", err)
	}
	logger, err := NewLogger(opt)
	if err != nil {
		fmt.Println("err = ", err)
	}
	//Out Some Target for project
	logger.Info("Raspberry Pi 4 and Pioneer600")
	logger.Info("Learn how to use golang control devices.")
	logger.Info("Thanks to https://gobot.io")
	logger.Info("Thanks to http://www.waveshare.net.")
	logger.Info("Thanks to https://periph.io/project/goals.")
	function := c.Int("function")
	switch function {
	case FunctionGpioLedOne:
		testGpioLEDOne()
	case FunctionPCF8574LedTwo:
		testPCF8574LedTwo()
	case FunctionPCF8574Beep:
		testPCF8574Beep()
	case FunctionDs18B20:
		testDS18b20()
	case FunctionDs3231:
		testDS3231()
	case FunctionSSD1306:
		testSSD1306()
	default:
		fmt.Printf("%v is not define yet.\n", function)
	}

	//quit when receive end signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.Printf("signal received signal %v", <-sigChan)
	logger.Warn("shutting down server")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "pi"
	app.Usage = "/usr/bin/pi -c /etc/pi/prod.yml"
	app.Version = "0.0.1"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "conf,c",
			Usage:  "Set conf path here",
			Value:  "prod.yml",
			EnvVar: "APP_CONF",
		},
		cli.IntFlag{
			Name:   "function,f",
			Usage:  "Test different device function",
			Value:  FunctionGpioLedOne,
			EnvVar: "APP_CONF",
		},
	}
	app.Run(os.Args)
}
