package typeconv

import (
	"errors"
	"strconv"
)

func ToArray(in interface{}) ([]interface{}, error) {
	var out []interface{}
	// fmt.Println("Type:", reflect.TypeOf(in).String())

	switch t := in.(type) {
	case []interface{}:
		for _, value := range t {
			out = append(out, value)
		}
	case *[]string:
		for _, value := range *t {
			out = append(out, value)
		}
	default:
		return nil, errors.New("Incompatible type")
	}

	return out, nil
}

func ToIntArray(in []interface{}) []int {
	var out []int

	for _, value := range in {
		out = append(out, value.(int))
	}

	return out
}

func ToStringArray(in []interface{}) []string {
	var out []string

	for _, value := range in {
		out = append(out, value.(string))
	}

	return out
}

func StringArray2IntArray(in []string) []int {
	var out []int

	for _, value := range in {
		i, _ := strconv.Atoi(value)
		out = append(out, i)
	}

	return out
}
