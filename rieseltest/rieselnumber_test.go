package rieseltest

import "testing"

// Test that NewRieselNumber accepts only h >= 1 and n >= 2
func TestNewRieselNumberErrors(t *testing.T) {
	var testCases = []struct {
		h int64
		n int64
	}{
		{0, 152},
		{-1, 252352},
		{-5, 2456464},
		{423423, 1},
		{9401, 0},
		{77, -2},
		{0, 0},
	}

	for _, c := range testCases {
		_, err := NewRieselNumber(c.h, c.n)
		if err == nil {
			t.Errorf("NewRieselNumber(%v, %v) should return an error, but it didn't", c.h, c.n)
		}
	}
}

// Test that NewRieselNumber forces h to be odd when it is not
func TestNewRieselNumber(t *testing.T) {
	var testCases = []struct {
		h int64
		n int64
		expected_h int64
		expected_n int64
	}{
		{1, 2, 1, 2},
		{773, 9768731, 773, 9768731},
		{6, 152, 3, 152},
		{224, 252352, 7, 252352},
	}

	for _, c := range testCases {
		R, err := NewRieselNumber(c.h, c.n)
		if err != nil && R.h != c.expected_h && R.n != c.expected_n {
			t.Errorf("NewRieselNumber(%v, %v) should result in R.h = %v and R.n = %v, but actually " +
				"we got R.h = %v and R.n = %v", c.h, c.n, c.expected_h, c.expected_n, R.h, R.n)
		}
	}
}