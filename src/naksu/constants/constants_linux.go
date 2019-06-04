package constants

// ExtNicsToIgnore is a list of devices to ignore. Each entry is a regular expression
// and the entry is dropped if any of these regexs match. See network.IgnoreExtInterface()
// for more.
var ExtNicsToIgnore = []string{
	"^lo$",
	"^vboxnet\\d",
}

