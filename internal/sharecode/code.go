package sharecode

import (
	"fmt"
	"math/rand"
)

func Generate() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
