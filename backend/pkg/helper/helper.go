package helper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func SliceAtoi(str string) []int {
	fmt.Println("STR: ", str)
	if str == "" || str == "[]" {
		return []int{}
	}

	// Remove the brackets and split the string into individual elements
	str = strings.Trim(str, "[]")
	elements := strings.Split(str, ",")

	// Convert each element to an integer and store it in a slice
	var intSlice []int
	for _, elem := range elements {
		elem = strings.Trim(elem, "\"")
		num, err := strconv.Atoi(elem)
		if err != nil {
			fmt.Println("Error converting to int SliceAtoi: ", err)
			return nil
		}
		intSlice = append(intSlice, num)
	}

	return intSlice
}

func MakeInt(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.Wrap(err, "failed from makeInt")
	}
	return v, nil
}
