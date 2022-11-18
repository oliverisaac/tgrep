package templating

import (
	"testing"
)

func TestParse(t *testing.T) {
	intRegex := "-?[0-9]+"
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "Input strings without any templates should return the same value",
			input:   "hello world",
			want:    "hello world",
			wantErr: false,
		},
		{
			name:    "If you pass in just a template name, it should return that template",
			input:   "int",
			want:    intRegex,
			wantErr: false,
		},
		{
			name:    "Templated regex should be not escaped",
			input:   "{{[0-9]}}",
			want:    "[0-9]",
			wantErr: false,
		},
		{
			name:    "Magic template words inside templates should expose the underlying regex",
			input:   "{{int}}",
			want:    intRegex,
			wantErr: false,
		},
		{
			name:    "If a string is outside the template it should not be touched",
			input:   "hello{{int}}world",
			want:    "hello" + intRegex + "world",
			wantErr: false,
		},
		{
			name:    "If a string outside the template contains regex characters, it should be escaped",
			input:   "[",
			want:    `\[`,
			wantErr: false,
		},
		{
			name:    "If a tempalte is not closed, that's a paddlin'",
			input:   "{{int",
			want:    "",
			wantErr: true,
		},
		{
			name:    "If double braces are inside a template, that should be fine",
			input:   "{{{{}}",
			want:    "{{",
			wantErr: false,
		},
		{
			name:    "If double braces are escaped, they should remain escaped",
			input:   `\{{name}}`,
			want:    `\{\{name\}\}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_matchesRegex(t *testing.T) {
	tests := []struct {
		name          string
		searchRegex   string
		stringToCheck string
		want          bool
		wantErr       bool
	}{
		{
			name:          "Validate that raw string matches raw string",
			searchRegex:   "hello world",
			stringToCheck: "hello world",
			want:          true,
			wantErr:       false,
		},
		{
			name:          "Validate that int matches 0",
			searchRegex:   "int",
			stringToCheck: "0",
			want:          true,
			wantErr:       false,
		},
		{
			name:          "Validate that int does not match a string",
			searchRegex:   "int",
			stringToCheck: "hello",
			want:          false,
			wantErr:       false,
		},
		{
			name:          "Validate that number handles negative numbers",
			searchRegex:   "number",
			stringToCheck: "-1",
			want:          true,
			wantErr:       false,
		},
		{
			name:          "Regex inside braces should work",
			searchRegex:   "{{[0-9]+}}",
			stringToCheck: "123",
			want:          true,
			wantErr:       false,
		},
		{
			name:          "Special characters in regex should be escaped",
			searchRegex:   "{{^}}(hello {{[0-9]+$}}",
			stringToCheck: "(hello 123",
			want:          true,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchesRegex(tt.stringToCheck, tt.searchRegex)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchesRegex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("matchesRegex() = %v, want %v", got, tt.want)
			}
		})
	}
}
