package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
