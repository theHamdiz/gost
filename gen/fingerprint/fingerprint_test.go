package fingerprint

import (
	"fmt"
	"testing"
)

func TestFingerprintingIsDifferentForTheSameAppName(t *testing.T) {
	input := "cody"

	f1, err := Fingerprint(input)

	if err != nil {
		t.Fatalf(">> Fingerprinting issue: %+v", err)
	}

	f2, err := Fingerprint(input)
	if err != nil {
		t.Fatalf(">> Fingerprinting issue: %+v", err)
	}

	t.Logf("First Fingerprint: %s", f1)
	t.Logf("Second Fingerprint: %s", f2)

	fmt.Printf("First Fingerprint: %s\n", f1)
	fmt.Printf("Second Fingerprint: %s\n", f2)

	if f2 == f1 {
		t.Errorf("> Expected %s not to be equal to %s", f1, f2)
	}
}
