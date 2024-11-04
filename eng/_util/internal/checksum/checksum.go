package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func WriteSHA256ChecksumFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	checksum := sha256.New()
	if _, err = io.Copy(checksum, file); err != nil {
		return err
	}
	// Write the checksum in a format that "sha256sum -c" can work with. Use the base path of the
	// tarball (not full path, not relative path) because then "sha256sum -c" automatically works
	// when the file and the checksum file are downloaded to the same directory.
	content := fmt.Sprintf("%v  %v\n", hex.EncodeToString(checksum.Sum(nil)), filepath.Base(path))
	outputPath := path + ".sha256"
	if err := os.WriteFile(outputPath, []byte(content), 0o666); err != nil {
		return err
	}
	fmt.Printf("Wrote checksum file %q with content: %v", outputPath, content)
	return nil
}
