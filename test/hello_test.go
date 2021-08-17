package test

import "testing"

func TestHello2(t *testing.T) {
	got := Hello()
	want :="Hello,world"

	if got != want{
		t.Errorf("got '%q' want '%q'", got, want)
	}
}
