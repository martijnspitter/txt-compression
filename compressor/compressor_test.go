package compressor

import (
	"testing"
)

func TestNewCompressor(t *testing.T) {
	tests := []struct {
		input    string
		expected map[rune]int
	}{
		{
			input:    "",
			expected: map[rune]int{},
		},
		{
			input:    "a",
			expected: map[rune]int{'a': 1},
		},
	}

	for _, test := range tests {
		c := NewCompressor(test.input)
		c.GenerateFreqTable()
		actual := *c.GetFreqTable()

		assertMapEqual(t, test.expected, actual)
	}

}

func assertMapEqual(t *testing.T, expected, actual map[rune]int) {
	if len(expected) != len(actual) {
		t.Errorf("Expected map length %d, got %d", len(expected), len(actual))
		t.FailNow()
	}

	for k, v := range expected {
		if actual[k] != v {
			t.Errorf("Expected value %d for key %c, got %d", v, k, actual[k])
			t.FailNow()
		}
	}
}
