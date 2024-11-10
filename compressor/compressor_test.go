package compressor

import (
	"bytes"
	"testing"
)

func TestFreqTable(t *testing.T) {
	tests := []struct {
		input    []byte
		expected map[byte]int
	}{
		{
			input:    []byte{},
			expected: map[byte]int{},
		},
		{
			input:    []byte("a"),
			expected: map[byte]int{'a': 1},
		},
	}

	for _, test := range tests {
		c := NewFrequencyTable()
		reader := bytes.NewReader(test.input)
		err := c.Create(reader)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			continue
		}
		actual := c.Get()

		assertMapEqual(t, test.expected, actual)
	}
}

func TestCodeTable(t *testing.T) {
	tests := []struct {
		input    []byte
		expected map[byte]string
	}{
		{
			input:    []byte{},
			expected: map[byte]string{},
		},
		{
			input:    []byte("a"),
			expected: map[byte]string{'a': "0"},
		},
		{
			input:    []byte("aa"),
			expected: map[byte]string{'a': "0"},
		},
		{
			input:    []byte("ab"),
			expected: map[byte]string{'a': "0", 'b': "1"},
		},
		{
			input:    []byte("abc"),
			expected: map[byte]string{'a': "10", 'b': "11", 'c': "0"},
		},
		{
			input:    []byte("abacabadabacaba"),
			expected: map[byte]string{'d': "000", 'b': "01", 'c': "001", 'a': "1"},
		},
	}

	for _, test := range tests {
		c := NewFrequencyTable()
		reader := bytes.NewReader(test.input) // Create a Reader from the string
		err := c.Create(reader)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			continue
		}
		bt := NewBinaryTree(&c.Table)
		bt.GetPrefixCodeTable()
		ct := bt.GetCodeTable()

		assertMapEqualString(t, test.input, test.expected, *ct)
	}
}

func assertMapEqual(t *testing.T, expected, actual map[byte]int) {
	if len(expected) != len(actual) {
		t.Errorf("Expected map length %d, got %d", len(expected), len(actual))
		t.FailNow()
	}

	for k, v := range expected {
		if actual[k] != v {
			t.Errorf("Expected value %d for key %02x, got %d", v, k, actual[k])
			t.FailNow()
		}
	}
}

func assertMapEqualString(t *testing.T, input []byte, expected map[byte]string, actual map[byte]string) {
	if len(expected) != len(actual) {
		t.Errorf("Expected map length %d, got %d for input %v", len(expected), len(actual), input)
		t.FailNow()
	}

	for k, v := range expected {
		if actual[k] != v {
			t.Errorf("Expected value %s for key %02x, got %s for input %v", v, k, actual[k], input)
			t.FailNow()
		}
	}
}
