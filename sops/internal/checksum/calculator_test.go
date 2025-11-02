package checksum_test

import (
	"testing"

	"github.com/carlpett/terraform-provider-sops/sops/internal/checksum"
)

func TestCalculateMD5(t *testing.T) {
	input := "Some text to calculate md5"
	expectedMd5 := "f6e1dd4fbcedfe10c7ffb832b961d1d0"

	result := checksum.CalculateMD5(input)

	if result != expectedMd5 {
		t.Errorf("Expected calculated hash %s, got %s", expectedMd5, result)
	}
}

func TestCalculateMD5OnLongString(t *testing.T) {
	input := `Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod 
tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam 
et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum
dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor 
invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua.`
	expectedMd5 := "a460956c1299344a6a7cd36111e14e55"

	result := checksum.CalculateMD5(input)

	if result != expectedMd5 {
		t.Errorf("Expected calculated hash %s, got %s", expectedMd5, result)
	}
}
