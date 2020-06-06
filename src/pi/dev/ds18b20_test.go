package dev

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestDs18b20Parser(t *testing.T) {
	data := "13 02 4b 46 7f ff 0d 10 e7 : crc=e7 YES\n13 02 4b 46 7f ff 0d 10 e7 t=33187"
	r := regexp.MustCompile("t=([0-9]+)")
	a := r.FindString(data)
	b := strings.ReplaceAll(a, "t=", "")
	temp, _ := strconv.ParseFloat(b, 64)
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(temp)
	temp = temp / 1000.0
	fmt.Println(temp)
}
