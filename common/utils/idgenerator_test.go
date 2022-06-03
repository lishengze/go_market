package utils

import (
	"fmt"
	"testing"
)

func TestCreateMemberId(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(GenOrderID())
	}
}
