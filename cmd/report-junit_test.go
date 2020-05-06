package cmd

import (
	"github.com/flanksource/build-tools/cmd/test"
	"github.com/flanksource/build-tools/pkg/gh/pr"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestGetReportJUnitCmd(t *testing.T) {
	cmd := GetReportJUnitCommand()
	assert.NotNil(t,cmd, "We must have a report-junit command")
}

func TestReportJunit_HasTokenField(t *testing.T) {
	cmd := GetReportJUnitCommand()
	token := "SOME_TOKEN_HERE"
	args := []string{"--auth-token",token}
	test.ParsesStringFlag(t,cmd,"auth-token",token, args,"We need to be able to parse an auth token to access Github")
}

func TestReportJunit_HasSilentSuccessField(t *testing.T) {
	cmd := GetReportJUnitCommand()
	args := []string{"--silent-success","true"}
	test.ParsesBoolFlag(t,cmd,"silent-success",true, args,"We want to pass a flag suppress a post on success only test reports")
}

func TestReportJUnit_HasReportJUnitSubcommand(t *testing.T) {
	cmd := GetReportJUnitCommand()
	test.HasSubcommand(t,cmd,"report-junit","We want a report-junit subcommand")
}

func Test_parseReportJUnitFlagsAndArguments(t *testing.T) {

	tests := []struct {
		name string
		args []string
		shouldFail bool
		wantedPR *pr.PR
		wantedFiles []string
	}{
		{
			name: "happy path",
			args: []string {"flanksource/platform-cli", "1","junit.xml", "--auth-token", "SOME_TOKEN"},
			shouldFail: false,
			wantedPR: &pr.PR{
				APIToken: "SOME_TOKEN",
				Owner:    "flanksource",
				Repo:     "platform-cli",
				Num:      1,
			},
			wantedFiles: []string{"junit.xml"},
		},
		{
			name: "too few arguments fails",
			args: []string {"--auth-token", "SOME_TOKEN"},
			shouldFail: true,
			wantedPR: nil,
			wantedFiles: []string{},
		},
		{
			name: "github token is mandatory",
			args: []string {"flanksource/platform-cli", "1","junit.xml"},
			shouldFail: true,
			wantedPR: nil,
			wantedFiles: []string{},
		},
		{
			name: "multiple files are parsed",
			args: []string {"flanksource/platform-cli", "1","junit1.xml", "junit2.xml", "--auth-token", "SOME_TOKEN"},
			shouldFail: false,
			wantedPR: &pr.PR{
				APIToken: "SOME_TOKEN",
				Owner:    "flanksource",
				Repo:     "platform-cli",
				Num:      1,
			},
			wantedFiles: []string{"junit1.xml", "junit2.xml" },
		},
		{
			name: "Invalid PR number fails",
			args: []string {"flanksource/platform-cli", "NaN","junit1.xml", "--auth-token", "SOME_TOKEN"},
			shouldFail: true,
			wantedPR: nil,
			wantedFiles: []string{},
		},
		{
			name: "Invalid owner/repo fails",
			args: []string {"flanksource","platform-cli", "1","junit1.xml", "--auth-token", "SOME_TOKEN"},
			shouldFail: true,
			wantedPR: nil,
			wantedFiles: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd:= GetReportJUnitCommand()
			cmd.ParseFlags(tt.args)
			err := validateReportJunitArguments(cmd, cmd.Flags().Args())

			if !tt.shouldFail {
				assert.NoError(t,err,"Validation failed with error: %v", err)
			} else if tt.shouldFail {
				t.Logf("Testcase %v should fail, but and did with error %v", tt.name, err)
				assert.Error(t,err,"Validation failed with error: %v", err)
				return // not testing further
			}
			gotPR, gotFiles, _, err := parseReportJunitFlagsAndArguments(cmd)
			assert.NoError(t,err,"Parsing failed with error: %v", err)
			assert.ElementsMatch(t,tt.wantedFiles,gotFiles,"We wanted different files.")
			t.Logf("%v -> %v, %v",tt.args, gotPR, gotFiles)
			assert.Equal(t,*(tt.wantedPR),gotPR,"We needed a PR to post a comment to.")
			//assert.Equal(t,tt.wantedPR,got,"We needed a PR to post a comment to.")
		})
	}
}


func Test_runReportJUnitCmd(t *testing.T) {
	// This is not an e2e test so we mostly verify code coverage
	// for paths causing errors

	tests := []struct {
		name    string
		args []string
		wantErr bool
	}{
		{
			name: "success path with no github API call - silent success",
			args: []string {"flanksource/platform-cli", "1","../fixtures/junit/junit-single-success.xml", "--auth-token", "SOME_TOKEN","--silent-success"},
			wantErr: false,
		},
		{
			name: "command arg parsing errors result in error",
			args: []string {"flanksource\\platform-cli", "1"},
			wantErr: true,
		},
		{
			name: "A JUnit report parsing error result in error",
			args: []string {"flanksource/platform-cli", "1","../fixtures/junit/malformed-xml.xml", "--auth-token", "SOME_TOKEN","--silent-success"},
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd:= GetReportJUnitCommand()
			cmd.ParseFlags(tt.args)
			err := runReportJUnitCmd(cmd, tt.args)
			if !tt.wantErr {
				assert.NoError(t,err,"Validation failed with error: %v", err)
			} else if tt.wantErr {
				t.Logf("Testcase %v should fail and did with error %v", tt.name, err)
				assert.Error(t,err)
			}
		})
	}
}