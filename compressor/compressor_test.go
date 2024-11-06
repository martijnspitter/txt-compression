package compressor

import (
	"testing"
)

func TestFreqTable(t *testing.T) {
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
		c := NewFrequencyTable()
		c.Create(test.input)
		actual := c.Get()

		assertMapEqual(t, test.expected, actual)
	}

}

func TestCodeTable(t *testing.T) {
	tests := []struct {
		input    string
		expected map[rune]string
	}{
		{
			input:    "",
			expected: map[rune]string{},
		},
		{
			input:    "a",
			expected: map[rune]string{'a': "0"},
		},
		{
			input:    "aa",
			expected: map[rune]string{'a': "0"},
		},
		{
			input:    "ab",
			expected: map[rune]string{'a': "0", 'b': "1"},
		},
		{
			input:    "abc",
			expected: map[rune]string{'a': "10", 'b': "11", 'c': "0"},
		},
		{
			input:    "abacabadabacaba",
			expected: map[rune]string{'d': "000", 'b': "01", 'c': "001", 'a': "1"},
		},
	}

	for _, test := range tests {
		c := NewFrequencyTable()
		c.Create(test.input)
		bt := NewBinaryTree(&c.Table)
		ct := bt.GetPrefixCodeTable()

		assertMapEqualString(t, test.expected, *ct)
	}
}

func assertMapEqualString(t *testing.T, expected map[rune]string, actual map[rune]string) {
	if len(expected) != len(actual) {
		t.Errorf("Expected map length %d, got %d", len(expected), len(actual))
		t.FailNow()
	}

	for k, v := range expected {
		if actual[k] != v {
			t.Errorf("Expected value %s for key %c, got %s", v, k, actual[k])
			t.FailNow()
		}
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
