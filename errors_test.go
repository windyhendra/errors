package errors

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	err := New(errors.New("Error Message"))
	if err == nil {
		t.Error("[TestError] Failed")
		return
	}

	err = New(err)
	if err == nil {
		t.Error("[TestError] Failed")
		return
	}

	t.Logf("[TestError] Success with err: %s", err.Error())
}

func TestErrorWithFields(t *testing.T) {
	err := New("Error Message", Fields{
		"Field 1": "Value 1",
		"Field 2": "Value 2",
	})
	if err != nil {
		t.Log("[TestErrorWithFields] Success with err: %s", err.Error())
		return
	}
	t.Error("[TestErrorWithFields] Failed")
}

func TestErrorWithHTTPError(t *testing.T) {
	err := New("Error Message", WithHTTPError(200, "OK"))
	if err != nil {
		t.Logf("[TestErrorWithHTTPError] Success with err: %s", err.Error())
		return
	}
	t.Error("[TestErrorWithHTTPError] Failed")
}

func TestBadError(t *testing.T) {
	err := New(1)
	if err != nil {
		t.Error("[TestBadError] Failed")
		return
	}
	t.Log("[TestBadError] Success")
}
