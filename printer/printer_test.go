package printer_test

import (
	"bytes"
	eb "github.com/dtomasi/go-event-bus/v3"
	"github.com/dtomasi/go-event-bus/v3/printer"
	"github.com/goombaio/namegenerator"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
	"text/tabwriter"
	"time"
)

var nameGenerator namegenerator.Generator //nolint:gochecknoglobals

//nolint:gochecknoinits
func init() {
	rand.Seed(time.Now().UnixNano())
	seed := time.Now().UTC().UnixNano()
	nameGenerator = namegenerator.NewNameGenerator(seed)
}

func linesStringCount(s string) int {
	n := strings.Count(s, "\n")
	if len(s) > 0 && !strings.HasSuffix(s, "\n") {
		n++
	}

	return n
}

func TestPrintStats(t *testing.T) {
	// Generate some random data
	var dataSlice []*eb.TopicStats

	for i := 0; i < 50; i++ {
		ts := &eb.TopicStats{
			Name:            nameGenerator.Generate(),
			PublishedCount:  eb.NewSafeCounter(),
			SubscriberCount: eb.NewSafeCounter(),
		}
		ts.PublishedCount.IncBy(uint(rand.Intn(10)))  //nolint:gosec
		ts.SubscriberCount.IncBy(uint(rand.Intn(10))) //nolint:gosec

		dataSlice = append(dataSlice, ts)
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	printer.PrintStatsTo(w, dataSlice)

	lineCount := linesStringCount(buf.String())

	// Line count should be loop count (50) + 2 lines header
	assert.Equal(t, 52, lineCount)
}
