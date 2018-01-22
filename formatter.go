package main

import (
	"fmt"
	"io"
)

type formatter interface {
	output(out io.Writer, stats []*workloadStat)
}

type tsvFormatter struct{}

func newTsvFormatter() *tsvFormatter {
	fmtr := &tsvFormatter{}
	return fmtr
}

func (f *tsvFormatter) output(out io.Writer, stats []*workloadStat) {
	fmt.Fprintf(out, "user\tdone\treviewed\n")
	for _, w := range stats {
		fmt.Fprintf(out, "%s\t%d\t%d\n", w.user, w.sentPullRequests, w.reviewedPullRequests)
	}
}
