package host_test

import (
  "testing"
  "naksu/host"
)

func TestGetWinProcessorAvailabilityLegend(t *testing.T) {
  tables := []struct {
    availabilityCode uint16
    expectedAvailabilityLegend string
  }{
    {1, "Other"},
    {21, "Quiesced"},
    {0, "N/A"},
    {40, "N/A"},
  }

  for _, table := range tables {
    observedLegend := host.GetWinProcessorAvailabilityLegend(table.availabilityCode)
    if observedLegend != table.expectedAvailabilityLegend {
      t.Errorf("getWinProcessorAvailabilityLegend() returns '%s' instead of expected '%s'", observedLegend, table.expectedAvailabilityLegend)
    }
  }
}
