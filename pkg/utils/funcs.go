package utils

import (
	"time"
	"fmt"
)

// Returns the first error that is not nil
func CheckAllErrors(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// Deletes the first instance of the pod with the given name in the list
func DeletePodNameOnce(lst []*Pod, name string) []*Pod {
	for i, p := range lst {
		if p.Name == name {
			return append(lst[:i], lst[i+1:]...)
		}
	}
	return lst
}

// Returns a list of all the values of a map
func GetValuesList(m map[string]*interface{}) []interface{} {
	r := make([]interface{}, 0, len(m))
	for _, value := range m {
		r = append(r, value)
	}
	return r
}

// Returns the string rep of the current time
func GetTimeString() string {
	t := time.Now()
	timeString := fmt.Sprintf("%d %s %d, %d:%d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return timeString
}