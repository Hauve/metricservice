package agent

import "fmt"

func (ag *MyAgent) sendMetrics() error {
	for _, m := range ag.storage.GetMetrics() {
		if err := ag.sender.Send(*m); err != nil {
			return fmt.Errorf("cannot send metric: %w", err)
		}
	}
	ag.storage.SetCounter("PollCount", 0)
	return nil
}
