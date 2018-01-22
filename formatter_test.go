package main

import (
	"testing"
)

func TestBuildFormatterFor(t *testing.T) {
	var (
		fmtr formatter
		err  error
	)

	fmtr, err = buildFormatterFor("tsv", "")
	if err != nil {
		t.Errorf("formatter for tsv can be built")
	}
	if _, ok := fmtr.(*tsvFormatter); !ok {
		t.Errorf("formatter should be tsvFormatter")
	}

	givenPrefix := "metric_prefix"
	fmtr, err = buildFormatterFor("sensu", givenPrefix)
	if err != nil {
		t.Errorf("formatter for sensu can be built")
	}
	if sensuFmtr, ok := fmtr.(*sensuFormatter); !ok {
		t.Errorf("formatter should be sensuFormatter")
	} else {
		if sensuFmtr.prefix != givenPrefix {
			t.Errorf("sensuFormatter has a prefix")
		}
	}

	fmtr, err = buildFormatterFor("unknown", "")
	if err == nil {
		t.Errorf("should die if unknown format given")
	}
}
