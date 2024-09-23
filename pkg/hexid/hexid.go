package hexid

import (
	"math/rand/v2"
	"regexp"
)

var (
	pattern = `^[0-9]{12}$`
)

// TODO: make UUID generator
func Generate() (int, error) {
	return rand.IntN(10000), nil
}

func Validate(id string) bool {
	matched, _ := regexp.MatchString(pattern, id)
	return matched
}
