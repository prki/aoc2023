package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreprocessLine(t *testing.T) {
	tIn := "xdzonexdzthree"
	expected := "xdz1xdz3"
	actual := PreprocessLine(tIn)
	assert.Equal(t, expected, actual)
}
