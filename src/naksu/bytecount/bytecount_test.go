package bytecount_test

import "testing"
import "naksu/bytecount"
import "math"

func TestBytesToHumanSI (t *testing.T) {
  tables := []struct {
    bytes int64
    human string
  }{
    {999, "999 B"},
    {1000, "1.0 kB"},
    {1023, "1.0 kB"},
    {1024, "1.0 kB"},
    {987654321, "987.7 MB"},
    {math.MaxInt64, "9.2 EB"},
  }

  for _, table := range tables {
    test_human := bytecount.BytesToHumanSI(table.bytes)
    if test_human != table.human {
      t.Errorf("BytesToHumanSI gives '%s' instead of '%s' for integer '%d'", test_human, table.human, table.bytes)
    }
  }
}

func TestBytesToHumanIEC (t *testing.T) {
  tables := []struct {
    bytes int64
    human string
  }{
    {999, "999 B"},
    {1000, "1000 B"},
    {1023, "1023 B"},
    {1024, "1.0 KiB"},
    {987654321, "941.9 MiB"},
    {math.MaxInt64, "8.0 EiB"},
  }

  for _, table := range tables {
    test_human := bytecount.BytesToHumanIEC(table.bytes)
    if test_human != table.human {
      t.Errorf("BytesToHumanIEC gives '%s' instead of '%s' for integer '%d'", test_human, table.human, table.bytes)
    }
  }
}
