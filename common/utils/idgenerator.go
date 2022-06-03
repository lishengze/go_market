package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func CreateMemberID() string {
	return fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
}

func GenOrderID() string {
	timeSeq := time.Now().Format("20060102150405.000")
	return fmt.Sprintf("%s", strings.Replace(timeSeq, ".", "", 1)[2:])
}

func GenTradeID() string {
	timeSeq := time.Now().Format("20060102150405.000")
	return fmt.Sprintf("t%s", strings.Replace(timeSeq, ".", "", 1)[2:])
}
