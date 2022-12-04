package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_runCommand(t *testing.T) {

	tests := []struct {
		name      string
		args      []string
		fileLines []string // if you set fileLines then they will be saved to a file and teh filename will be appended ot the args
		sendLines []string
		wantLines []string
		wantErr   bool
	}{
		{
			name:    "No arguments",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "Test for help flag",
			args:    []string{"--help"},
			wantErr: true,
		},
		{
			name:      "Most basic regex test",
			args:      []string{"hello"},
			sendLines: []string{"hello"},
			wantLines: []string{"hello"},
			wantErr:   false,
		},
		{
			name:      "Test multi line input",
			args:      []string{"hello"},
			sendLines: []string{"hello", "world"},
			wantLines: []string{"hello"},
			wantErr:   false,
		},
		{
			name:      "Test -e flag",
			args:      []string{"-e", "hello"},
			sendLines: []string{"hello", "world"},
			wantLines: []string{"hello"},
			wantErr:   false,
		},
		{
			name:      "Test multiple -e flag",
			args:      []string{"-e", "world", "-e", "hello"},
			sendLines: []string{"hello", "bob", "world"},
			wantLines: []string{"hello", "world"},
			wantErr:   false,
		},
		{
			name:      "Test templating",
			args:      []string{"-e", "email"},
			sendLines: []string{"hello", "bob", "example@example.com"},
			wantLines: []string{"example@example.com"},
			wantErr:   false,
		},
		{
			name:      "Test templating",
			args:      []string{"-e", "(my email is {{email}})"},
			sendLines: []string{"hello", "bob", "Hello, (my email is example@example.com)"},
			wantLines: []string{"Hello, (my email is example@example.com)"},
			wantErr:   false,
		},
		{
			name:      "Test filename",
			args:      []string{"-e", "email"},
			fileLines: []string{"hello world", "my email is example@example.com"},
			sendLines: []string{},
			wantLines: []string{"my email is example@example.com"},
			wantErr:   false,
		},
		{
			name:      "Case insensitive",
			args:      []string{"-i", "-e", "hello"},
			fileLines: []string{"HELLO WORLD", "my email is example@example.com"},
			sendLines: []string{},
			wantLines: []string{"HELLO WORLD"},
			wantErr:   false,
		},
		{
			name:      "Non-templated regex",
			args:      []string{"-E", "-e", "he[lo]{3}"},
			fileLines: []string{"hello world", "my email is example@example.com"},
			sendLines: []string{},
			wantLines: []string{"hello world"},
			wantErr:   false,
		},
		{
			name:      "Case insensitive Non-templated regex",
			args:      []string{"-E", "-i", "-e", "he[lo]{3}"},
			fileLines: []string{"HELLO WORLD", "my email is example@example.com"},
			sendLines: []string{},
			wantLines: []string{"HELLO WORLD"},
			wantErr:   false,
		},
		{
			name:      "Word boundary templated regex",
			args:      []string{"-w", "-e", "mail"},
			fileLines: []string{"HELLO WORLD", "i like mail", "my email is example@example.com"},
			sendLines: []string{},
			wantLines: []string{"i like mail"},
			wantErr:   false,
		},
		{
			name:      "Word boundary non-templated regex",
			args:      []string{"-E", "-w", "-e", "mai[a-z]"},
			fileLines: []string{"HELLO WORLD", "i like mail", "my email is example@example.com"},
			sendLines: []string{},
			wantLines: []string{"i like mail"},
			wantErr:   false,
		},
		{
			name:      "Combine short flags",
			args:      []string{"-Ewe", "mai[a-z]"},
			fileLines: []string{"HELLO WORLD", "i like mail", "my email is example@example.com"},
			sendLines: []string{},
			wantLines: []string{"i like mail"},
			wantErr:   false,
		},
		{
			name:      "Double dash switches to positional args",
			args:      []string{"-Ew", "--", "-mai[a-z]"},
			fileLines: []string{"HELLO WORLD", "i-like-mail", "my email is example@example.com"},
			sendLines: []string{},
			wantLines: []string{"i-like-mail"},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.fileLines) > 0 {
				tempFile, err := ioutil.TempFile(os.TempDir(), "testing-tgrep")
				if err != nil {
					panic(err)
				}
				defer os.Remove(tempFile.Name())

				tt.args = append(tt.args, tempFile.Name())

				for _, l := range tt.fileLines {
					tempFile.WriteString(l + "\n")
				}
				tempFile.Close()
			}
			output := &bytes.Buffer{}
			input := bytes.NewBufferString(strings.Join(tt.sendLines, "\n"))
			if err := runCommand(tt.args, input, output); err != nil {
				if tt.wantErr {
					return
				}
				t.Errorf("runCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput := output.String(); gotOutput != strings.Join(tt.wantLines, "\n")+"\n" {
				t.Errorf("runCommand() = %v, want %v", strings.Split(strings.TrimSuffix(gotOutput, "\n"), "\n"), tt.wantLines)
			}
		})
	}
}
