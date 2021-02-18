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
