package utils 

import (
	"regexp"
)

func IsSafeString(input string) bool {

	injectionPattern := `(?i)(\bSELECT\b|\bINSERT\b|\bUPDATE\b|\bDELETE\b|\bDROP\b|\bUNION\b|\bEXEC\b|\bALTER\b|\bTRUNCATE\b|\bOR\b|\bAND\b|\b;|\b--\s)`
	regex := regexp.MustCompile(injectionPattern)
	match := regex.MatchString(input)

	return !match
}