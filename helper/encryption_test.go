package helper

import (
	"fmt"
	"mpt_data/helper/config"
	"os"
	"testing"
)

func TestEncryption(t *testing.T) {
	if err := os.Chdir(".."); err != nil {
		fmt.Println(err)
	}
	config.LoadConfig()
	t.Run(
		"Test",
		func(t *testing.T) {
			orig := "Test"
			result, err := EncryptData(orig)
			if err != nil {
				t.Errorf("got error: %v", err)
			}
			decrypt, err := DecryptData(result)
			if err != nil {
				t.Errorf("got error: %v", err)
			}
			if orig != string(decrypt) {
				t.Errorf("expected %s, got %s", orig, string(decrypt))
			}
		},
	)
}
