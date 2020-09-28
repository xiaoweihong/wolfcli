package utils

import (
	"fmt"
	"strings"
	"testing"
)

func TestConverArceeURLToWeedUrl(t *testing.T) {
	urlTest := "http://192.168.2.6:8501/api/v2/file/4/114e7b66d03c15"
	url := ConverArceeURLToWeedUrl(urlTest)
	fmt.Println(url)
	contains := strings.Contains(urlTest, "api/v2/file")
	fmt.Println(contains)
}
