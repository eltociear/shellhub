package utils

import (
	"errors"
	"fmt"
	"regexp"
)

func GetFieldValueFromJSON(fieldName string, byteStream string) (string, error) {
	re := regexp.MustCompile(fmt.Sprintf(`(?m)%s":"(.*?)",`, fieldName))

	submatches := re.FindStringSubmatch(byteStream)
	if len(submatches) < 2 {
		return "", errors.New("couldnt match")
	}

	return submatches[1], nil
}
