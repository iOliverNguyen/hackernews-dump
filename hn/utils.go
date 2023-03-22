package hn

import "fmt"

func wrapf(err error, s string, args ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(spr(s, args...) + ": " + err.Error())
}

func spr(s string, args ...interface{}) string {
	return fmt.Sprintf(s, args...)
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
