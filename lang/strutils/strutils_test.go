package strutils_test

import (
	"testing"

	"github.com/renatopp/golden/lang/strutils"
	"github.com/stretchr/testify/assert"
)

func TestEscapeNewlines(t *testing.T) {
	target := "Hello\nWorld"
	expected := "Hello\\nWorld"
	actual := strutils.EscapeNewlines(target)
	assert.Equal(t, expected, actual)
}

func TestPadLeft(t *testing.T) {
	target := "hi"
	expected := "    hi"
	actual := strutils.PadLeft(target, 6)
	assert.Equal(t, expected, actual)
}

func TestPadLeftWith(t *testing.T) {
	target := "hi"
	expected := "----hi"
	actual := strutils.PadLeftWith(target, 6, "-")
	assert.Equal(t, expected, actual)
}

func TestPadRight(t *testing.T) {
	target := "hi"
	expected := "hi    "
	actual := strutils.PadRight(target, 6)
	assert.Equal(t, expected, actual)
}

func TestPadRightWith(t *testing.T) {
	target := "hi"
	expected := "hi----"
	actual := strutils.PadRightWith(target, 6, "-")
	assert.Equal(t, expected, actual)
}

func TestPadCenter(t *testing.T) {
	target := "hi!"
	expected := "  hi! "
	actual := strutils.PadCenter(target, 6)
	assert.Equal(t, expected, actual)
}

func TestPadCenterWith(t *testing.T) {
	target := "hi!"
	expected := "--hi!-"
	actual := strutils.PadCenterWith(target, 6, "-")
	assert.Equal(t, expected, actual)
}

func TestRepeat(t *testing.T) {
	target := "hi"
	expected := "hihihi"
	actual := strutils.Repeat(target, 3)
	assert.Equal(t, expected, actual)
}

func TestReverse(t *testing.T) {
	targets := []string{"heLLo", "WORLD", "ÇÃO", "123"}
	expecteds := []string{"oLLeh", "DLROW", "OÃÇ", "321"}
	for i, target := range targets {
		expected := expecteds[i]
		actual := strutils.Reverse(target)
		assert.Equal(t, expected, actual)
	}
}

func TestReplace(t *testing.T) {
	target := "Hello World "
	expected := "Hello, World, "
	actual := strutils.Replace(target, " ", ", ")
	assert.Equal(t, expected, actual)
}

func TestToLowerCase(t *testing.T) {
	targets := []string{"heLLo", "WORLD", "ÇÃO", "123"}
	expecteds := []string{"hello", "world", "ção", "123"}
	for i, target := range targets {
		expected := expecteds[i]
		actual := strutils.ToLowerCase(target)
		assert.Equal(t, expected, actual)
	}
}

func TestToUpperCase(t *testing.T) {
	targets := []string{"heLLo", "WORLD", "çãO", "123"}
	expecteds := []string{"HELLO", "WORLD", "ÇÃO", "123"}
	for i, target := range targets {
		expected := expecteds[i]
		actual := strutils.ToUpperCase(target)
		assert.Equal(t, expected, actual)
	}
}

func TestToTitleCase(t *testing.T) {
	targets := []string{"heLLo", "WORLD", "çãO", "123"}
	expecteds := []string{"Hello", "World", "Ção", "123"}
	for i, target := range targets {
		expected := expecteds[i]
		actual := strutils.ToTitleCase(target)
		assert.Equal(t, expected, actual)
	}
}

// func TestToCamelCase(t *testing.T) {
// 	targets := []string{"heLLo woRLD", "WORLD", "çãO sample", "123"}
// 	expecteds := []string{"helloWorld", "world", "çãoSample", "123"}
// 	for i, target := range targets {
// 		expected := expecteds[i]
// 		actual := strutils.ToCamelCase(target)
// 		assert.Equal(t, expected, actual)
// 	}
// }

// func TestToPascalCase(t *testing.T) {
// 	targets := []string{"heLLo-woRLD ", "WORLD", "çãO sample", "123"}
// 	expecteds := []string{"HelloWorld", "World", "ÇãoSample", "123"}
// 	for i, target := range targets {
// 		expected := expecteds[i]
// 		actual := strutils.ToPascalCase(target)
// 		assert.Equal(t, expected, actual)
// 	}
// }
