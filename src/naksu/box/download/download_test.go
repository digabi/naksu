package download

import (
	"os"
	"testing"
)

func TestCleanSHA256ChecksumString(t *testing.T) {
	testCases := []struct {
		sourceChecksum   string
		expectedChecksum string
	}{
		{"aff72a2cd83323e21c48d8686f3bcb7469b5131eb678a1bdcdbef27ff4f05b94", "aff72a2cd83323e21c48d8686f3bcb7469b5131eb678a1bdcdbef27ff4f05b94"},
		{"aff72a2cd83323e21c48d8686f3bcb7469b5131eb678a1bdcdbef27ff4f05b94 ", "aff72a2cd83323e21c48d8686f3bcb7469b5131eb678a1bdcdbef27ff4f05b94"},
		{"aff72a2cd83323e21c48d8686f3bcb7469b5131eb678a1bdcdbef27ff4f05b94 foo", "aff72a2cd83323e21c48d8686f3bcb7469b5131eb678a1bdcdbef27ff4f05b94"},
		{`aff72a2cd83323e21c48d8686f3bcb7469b5131eb678a1bdcdbef27ff4f05b94 foo\nbar`, "aff72a2cd83323e21c48d8686f3bcb7469b5131eb678a1bdcdbef27ff4f05b94"},
		{"  aff72a2cd83323e21c48d8686f3bcb7469b5131eb678a1bdcdbef27ff4f05b94", ""},
		{"ff72a2cd83323e21c48d8686f3bcb7469b5131eb678a1bdcdbef27ff4f05b94", ""},
	}

	for _, testCase := range testCases {
		if CleanSHA256ChecksumString(testCase.sourceChecksum) != testCase.expectedChecksum {
			t.Errorf("CleanSHA256ChecksumString fails, [%s] should be [%s]", testCase.sourceChecksum, testCase.expectedChecksum)
		}
	}
}

func TestGetSHA256ChecksumFromFile(t *testing.T) {
	testCases := []struct {
		fileContent      string
		expectedChecksum string
	}{
		{"foobar", "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2"},
		{"foobar\n\n", "de9e0da624af6ccda42959e1f5a183153467831f5e81faadc641aaebf662e6cd"},
	}

	for _, testCase := range testCases {
		testOneGetSHA256ChecksumFromFile(t, testCase.fileContent, testCase.expectedChecksum)
	}
}

func testOneGetSHA256ChecksumFromFile(t *testing.T, fileContent string, expectedChecksum string) {
	tempFile, err := os.CreateTemp(os.TempDir(), "naksu-test-")
	if err != nil {
		t.Errorf("Cannot create temporary file: %v", err)
	}

	testContent := []byte(fileContent)
	if _, err = tempFile.Write(testContent); err != nil {
		t.Errorf("Failed to write to temporary file %s: %v", tempFile.Name(), err)
	}

	if err := tempFile.Close(); err != nil {
		t.Errorf("Failed to close temporary file %s: %v", tempFile.Name(), err)
	}

	nilProgressCallbackFn := func(message string, value int) {
		// fmt.Println(message)
	}

	calculatedChecksum, err := GetSHA256ChecksumFromFile(tempFile.Name(), nilProgressCallbackFn)
	if err != nil {
		t.Errorf("Error while calculating checksum from file %s: %v", tempFile.Name(), err)
	}

	if calculatedChecksum != expectedChecksum {
		// Spare the temporary file for debugging on error
		t.Errorf("GetSHA256ChecksumFromFile fails for file %s, calculated: %s, expected: %s", tempFile.Name(), calculatedChecksum, expectedChecksum)
	} else {
		os.Remove(tempFile.Name())
	}
}
