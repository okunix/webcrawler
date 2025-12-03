package main

import "testing"

func TestExtractRawLinks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "multiple valid hrefs",
			input:    `<a href="https://example.com">Link1</a><a href="/path">Link2</a>`,
			expected: []string{"https://example.com", "/path"},
		},
		{
			name:     "no links",
			input:    `Just text without links`,
			expected: []string{},
		},
		{
			name:     "malformed href",
			input:    `<a href=invalid>Bad</a>`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractRawLinks(tt.input)
			t.Logf("%+v\n", got)
			if len(got) != len(tt.expected) {
				t.Errorf("expected %d links, got %d", len(tt.expected), len(got))
				return
			}
			for i, expected := range tt.expected {
				if got[i] != expected {
					t.Errorf("expected %q at index %d, got %q", expected, i, got[i])
				}
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		links    []string
		expected []string
	}{
		{
			name:     "absolute urls unchanged",
			baseURL:  "https://example.com",
			links:    []string{"https://google.com", "http://other.com"},
			expected: []string{"https://google.com", "http://other.com"},
		},
		{
			name:     "relative paths normalized",
			baseURL:  "https://example.com",
			links:    []string{"/api/users", "/path?query=1"},
			expected: []string{"https://example.com/api/users", "https://example.com/path?query=1"},
		},
		{
			name:     "mixed and invalid filtered",
			baseURL:  "https://example.com",
			links:    []string{"https://valid.com", "invalid", "/good", "mailto:email"},
			expected: []string{"https://valid.com", "https://example.com/good"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.baseURL, tt.links)
			t.Logf("%+v\n", got)
			if len(got) != len(tt.expected) {
				t.Errorf("expected %d links, got %d", len(tt.expected), len(got))
				return
			}
			for i, expected := range tt.expected {
				if got[i] != expected {
					t.Errorf("expected %q at index %d, got %q", expected, i, got[i])
				}
			}
		})
	}
}

func TestBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "valid absolute url",
			input: "https://example.com/path?query=1",
			want:  "https://example.com",
		},
		{
			name:  "valid http url",
			input: "http://test.com/api",
			want:  "http://test.com",
		},
		{
			name:    "relative path",
			input:   "/path/to/resource",
			wantErr: true,
		},
		{
			name:    "invalid url",
			input:   "invalid://url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BaseURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("BaseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("BaseURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractLinks(t *testing.T) {
	html := `<html><a href="https://external.com">External</a><a href="/internal">Internal</a></html>`
	base := "https://mysite.com"

	expected := []string{
		"https://external.com",
		"https://mysite.com/internal",
	}

	got := ExtractLinks(html, base)

	if len(got) != len(expected) {
		t.Errorf("expected %d links, got %d", len(expected), len(got))
		return
	}

	for i, exp := range expected {
		if got[i] != exp {
			t.Errorf("expected %q at %d, got %q", exp, i, got[i])
		}
	}
}
