package eventbus

type TopicStats struct {
	Name            string
	PublishedCount  *SafeCounter
	SubscriberCount *SafeCounter
}

type topicStatsMap map[string]*TopicStats

type Stats struct {
	data topicStatsMap
}

func newStats() *Stats {
	return &Stats{
		data: map[string]*TopicStats{},
	}
}

func (s *Stats) getOrCreateTopicStats(topicName string) *TopicStats {
	_, ok := s.data[topicName]
	if !ok {
		s.data[topicName] = &TopicStats{
			Name:            topicName,
			PublishedCount:  NewSafeCounter(),
			SubscriberCount: NewSafeCounter(),
		}
	}

	return s.data[topicName]
}

func (s *Stats) incSubscriberCountByTopic(topicName string) {
	s.getOrCreateTopicStats(topicName).SubscriberCount.Inc()
}

func (s *Stats) GetSubscriberCountByTopic(topicName string) int {
	return s.getOrCreateTopicStats(topicName).SubscriberCount.Value()
}

func (s *Stats) incPublishedCountByTopic(topicName string) {
	s.getOrCreateTopicStats(topicName).PublishedCount.Inc()
}

func (s *Stats) GetPublishedCountByTopic(topicName string) int {
	return s.getOrCreateTopicStats(topicName).PublishedCount.Value()
}

func (s *Stats) GetTopicStats() []*TopicStats {
	var tStatsSlice []*TopicStats
	for _, tStats := range s.data {
		tStatsSlice = append(tStatsSlice, tStats)
	}

	return tStatsSlice
}

func (s *Stats) GetTopicStatsByName(topicName string) *TopicStats {
	return s.getOrCreateTopicStats(topicName)
}
