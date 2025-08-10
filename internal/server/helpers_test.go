package server

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name:           "valid on input",
			input:          "some-plugin-name-1-23-45",
			expectedOutput: "some-plugin-name-1-23-45",
		},
		{
			name:           "lowercase",
			input:          "somepluginname",
			expectedOutput: "somepluginname",
		},
		{
			name:           "uppercase",
			input:          "SOMEPLUGINNAME",
			expectedOutput: "somepluginname",
		},
		{
			name:           "numbers",
			input:          "12345",
			expectedOutput: "12345",
		},
		{
			name:           "special symbols",
			input:          "!@#$%^&*()_+=*/|\\,.",
			expectedOutput: "",
		},
		{
			name:           "separated by spaces",
			input:          "some plugin  name",
			expectedOutput: "some-plugin-name",
		},
		{
			name:           "mixed",
			input:          "   Some plugin   NaMe !@#$ %^ &*(1.23.45)_+=, /|\\  ",
			expectedOutput: "some-plugin-name-1-23-45",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			slug := slugify(test.input)
			if slug != test.expectedOutput {
				t.Errorf("slugify() slug = %v, want %v", slug, test.expectedOutput)
			}
		})
	}
}
