package mebroutines_test

import "testing"
import "naksu/mebroutines"

func TestIfIntlCharsInPath (t *testing.T) {
  tables := []struct {
    path string
    result bool
  }{
    {"C:\\Users\\john.doe\\ktp-jako", false},
    {"C:\\Users\\raimo.keski-vääntö\\ktp-jako", true},
    {"C:\\Users\\john doe\\ktp-jako", false},

    {"/home/someuser/ktp-jako", false},
    {"/home/öylätti/ktp-jako", true},
    {"~/ktp-jako", true},
    {"~root/loremipsum", true},

    {"random whatever string", false},
    {"random whatever string with öljyrätti", true},
    {"wtf!", true},
    {"what?", true},
    {"/home/ktp-user/*", true},
  }

  for _, table := range tables {
    isIntl := mebroutines.IfIntlCharsInPath(table.path)
    if isIntl != table.result {
      t.Errorf("IfIntlCharsInPath gives '%t' instead of '%t' for path '%s'", isIntl, table.result, table.path)
    }
  }
}
