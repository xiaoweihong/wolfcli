package controller

import (
	"fmt"
	"testing"
	"time"
)

func TestGetVolumeInfo(t *testing.T) {
	a := time.Now().Format("2006-01-02")
	fmt.Println(a)
	parse, _ := time.Parse("2006-01-02", a)
	fmt.Println(parse)
	//time.Unix()
}
