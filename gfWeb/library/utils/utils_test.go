package utils

import (
	"testing"
)

func TestRandomString(t *testing.T) {
	str := RandomString(10)
	Initglog()
	//Debug("asdfasdfasdf")
	t.Error(str)
}
