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


// Returns the string rep of the current time
func GetTimeString() string {
	t := time.Now()
	timeString := fmt.Sprintf("%d %s %d, %d:%d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	return timeString
}