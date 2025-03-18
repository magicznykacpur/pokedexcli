package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "big boy CHADINGTON",
			expected: []string{"big", "boy", "chadington"},
		},
		{
			input:    " can i HAZ cheeseburger",
			expected: []string{"can", "i", "haz", "cheeseburger"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("actual doesn't match expected: %d != %d", len(actual), len(c.expected))
		}

		for i := range actual {
			word := actual[i]
			expected := c.expected[i]

			if word != expected {
				t.Errorf("actual doesn't match expected: %s != %s", word, expected)
			}
		}
	}
}
