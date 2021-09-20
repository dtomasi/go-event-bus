package printer

import (
	"github.com/cheynewallace/tabby"
	eb "github.com/dtomasi/go-event-bus/v3"
	"os"
	"text/tabwriter"
)

func PrintStats(topicStats []*eb.TopicStats) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0) //nolint:gomnd
	PrintStatsTo(w, topicStats)
}

func PrintStatsTo(writer *tabwriter.Writer, topicStats []*eb.TopicStats) {
	t := tabby.NewCustom(writer)
	t.AddHeader("Topic", "Subscriber Count", "Published Count")

	for _, ts := range topicStats {
		t.AddLine(ts.Name, ts.SubscriberCount.Value(), ts.PublishedCount.Value())
	}

	t.Print()
}
