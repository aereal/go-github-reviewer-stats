package main

import (
	"fmt"
	"io"
	"time"
)

type formatter interface {
	output(out io.Writer, stats []*workloadStat)
}

func buildFormatterFor(format string) (formatter, error) {
	switch format {
	case OUTPUT_TSV:
		fmtr := newTsvFormatter()
		return fmtr, nil
	case OUTPUT_SENSU:
		fmtr := newSensuFormatter("pull_requests")
		return fmtr, nil
	default:
		return nil, fmt.Errorf("Unknown format: %s", format)
	}
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

type sensuFormatter struct {
	prefix string
}

func newSensuFormatter(prefix string) *sensuFormatter {
	fmtr := &sensuFormatter{prefix: prefix}
	return fmtr
}

func (f *sensuFormatter) output(out io.Writer, stats []*workloadStat) {
	ts := time.Now()
	for _, w := range stats {
		fmt.Fprintf(out, "%s.%s.%s\t%d\t%d\n", f.prefix, "sent", w.user, w.sentPullRequests, ts.Unix())
		fmt.Fprintf(out, "%s.%s.%s\t%d\t%d\n", f.prefix, "reviewed", w.user, w.reviewedPullRequests, ts.Unix())
	}
}
