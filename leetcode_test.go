package grabb3r

import (
	"os"
	"testing"
)

func newLeetcodeSource() SolutionSource {
	user, _ := os.LookupEnv("LEETCODE_USER")
	password, _ := os.LookupEnv("LEETCODE_PASSWORD")
	if len(user) == 0 || len(password) == 0 {
		return nil
	}
	return NewLeetCodeSource(user, password)
}

func TestLeetcodeSource_Login(t *testing.T) {
	ls := newLeetcodeSource()
	if ls == nil {
		t.Skip("credentials not provided")
	}
	err := ls.Login()
	if err != nil {
		t.Error(err)
	}
}

func TestLeetcodeSource_ListSolutions(t *testing.T) {
	ls := newLeetcodeSource()
	if ls == nil {
		t.Skip("credentials not provided")
	}
	err := ls.Login()
	if err != nil {
		t.Error(err)
	}
	solutionChan, errChan := ls.ListSolutions()
	select {
	case <-solutionChan:
	case err := <-errChan:
		t.Error(err)
	}
}
