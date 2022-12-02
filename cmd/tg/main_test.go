package main

import (
	"bytes"
	"testing"

	"github.com/oliverisaac/tgrep/pkg/templating"
)

func Test_runCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "No templating",
			args:       []string{"hello"},
			wantOutput: "hello\n",
			wantErr:    false,
		},
		{
			name:       "Multiple arguments",
			args:       []string{"hello", "world"},
			wantOutput: "hello\nworld\n",
			wantErr:    false,
		},
		{
			name:       "No templating",
			args:       []string{},
			wantOutput: "",
			wantErr:    true,
		},
		{
			name:       "Some templating",
			args:       []string{"Hello {{world|bob}}"},
			wantOutput: "Hello world|bob\n",
			wantErr:    false,
		},
		{
			name:       "With quick templates",
			args:       []string{"email"},
			wantOutput: templating.GetTemplates()["email"] + "\n",
			wantErr:    false,
		},
		{
			name:       "If ask for help",
			args:       []string{"email", "--help"},
			wantOutput: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			output := &bytes.Buffer{}
			err = runCommand(tt.args, output)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Errorf("runCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotOutput := output.String(); gotOutput != tt.wantOutput {
				t.Errorf("runCommand() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
