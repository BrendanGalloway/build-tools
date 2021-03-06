/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package junit

import (
	"fmt"
	"github.com/google/go-github/v32/github"
	"github.com/joshdk/go-junit"
)

var (
	AnnotationWarning = "warning"
	AnnotationNotice  = "notice"
	AnnotationFailure = "failure"
)

const mdTableHeader = `| Result | Class | Message |
|--------|-------|--------|
`

func (results TestResults) GenerateMarkdown() string {
	mdResult := ""
	mdResult += fmt.Sprintf("<details><summary>%d test suites - Totals:  %d tests, %d failed, %d skipped, %d passed</summary>\n\n",len(results.Suites), results.Total, results.Failed, results.Skipped,  results.Passed)

	for _, suite := range results.Suites {
		//44 tests, 21 passed, 0 failed, 23 skipped
		mdResult += fmt.Sprintf("<details><summary>%s:  %d tests, %d failed, %d skipped, %d passed</summary>\n\n",suite.Name, suite.Totals.Tests, suite.Totals.Failed, suite.Totals.Skipped, suite.Totals.Passed)
		mdResult += mdTableHeader
		for _, test := range suite.Tests {
			switch test.Status {
			case junit.StatusFailed:
				mdResult += fmt.Sprintf("| :x: | **%s** | `%s` |\n", test.Classname, test.Name)
				// no default:
				// any other status will be ignored
			}
		}
		for _, test := range suite.Tests {
			switch test.Status {
			case junit.StatusSkipped:
				mdResult += fmt.Sprintf("| :white_circle: | **%s** | `%s` |\n", test.Classname, test.Name)
				// no default:
				// any other status will be ignored
			}
		}
		mdResult += "\n</details>\n"
	}
	mdResult += "\n</details>\n"
	return mdResult
}

func toPtr(s string) *string {
	return &s
}

const MAX_ANNOTATIONS = 50

func (results TestResults) GetGithubAnnotations() []*github.CheckRunAnnotation {
	list := []*github.CheckRunAnnotation{}
	count := 0
	for _, suite := range results.Suites {
		count++
		if count > MAX_ANNOTATIONS {
			return list
		}
		for _, test := range suite.Tests {
			msg := fmt.Sprintf("stdout:%s stderr:%s", test.SystemOut, test.SystemErr)
			zero := 0
			annotation := &github.CheckRunAnnotation{
				Title:     &test.Classname,
				StartLine: &zero,
				EndLine:   &zero,
				Path:      &test.Name,
				Message:   &msg,
			}

			switch test.Status {
			case junit.StatusFailed:
				annotation.AnnotationLevel = &AnnotationFailure
			default:
				continue
			}

			list = append(list, annotation)
		}
	}

	return list
}
