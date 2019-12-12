package main

import (
	"testing"
	"time"
)

// Patch to correct the breaking change in Golang 13
// c.f. https://stackoverflow.com/a/58192326/3893444
// c.f. https://github.com/golang/go/issues/31859#issuecomment-489889428
var _ = func() bool {
	testing.Init()
	return true
}()

func TestHumanDate(t *testing.T) {

	// Initialize a new time.Time object and pass it to the humanDate function
	tm := time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC)
	hd := humanDate(tm)

	// Check that the output from the humanDate function is in the format we // expect. If it isn't what we expect, use the t.Errorf() function to
	// indicate that the test has failed and log the expected and actual
	// values.
	if hd != "17 Dec 2020 at 10:00" {
		t.Errorf("want %q; got %q", "17 Dec 2020 at 10:00", hd)
	}
}
