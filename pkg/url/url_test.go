package url

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple text",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "with special characters",
			input:    "Hello, World!",
			expected: "hello-world",
		},
		{
			name:     "already lowercase",
			input:    "hello world",
			expected: "hello-world",
		},
		{
			name:     "with spaces",
			input:    "  leading and trailing  ",
			expected: "leading-and-trailing",
		},
		{
			name:     "with numbers",
			input:    "Product 123",
			expected: "product-123",
		},
		{
			name:     "multiple separators",
			input:    "multiple---separators___here",
			expected: "multiple-separators-here",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "already a slug",
			input:    "my-slug-here",
			expected: "my-slug-here",
		},
		{
			name:     "mixed case with accents",
			input:    "Café Au Lait",
			expected: "caf-au-lait",
		},
		{
			name:     "unicode characters",
			input:    "日本語テスト",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Slugify(tt.input)
			if result != tt.expected {
				t.Errorf("Slugify(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
