package controller

import (
	"fmt"
	"testing"
	"time"
)

func TestGetVolumeInfo(t *testing.T) {
	fmt.Println(time.Now())
	zero := TimeNowZero()
	fmt.Println(zero)
	//time.Unix()
}

func TestDisplay(t *testing.T) {
	var s [][]string
	a1 := []string{"a"}
	s = append(s, a1)
	fmt.Println(s)
}
