package load

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/anchore/quill/internal/log"
)

func BytesFromFileOrEnv(path string) ([]byte, error) {
	if strings.HasPrefix(path, "env:") {
		// comes from an env var...
		fields := strings.Split(path, "env:")
		if len(fields) < 2 {
			return nil, fmt.Errorf("key path has 'env:' prefix, but cannot parse env variable: %q", path)
		}
		envVar := fields[1]

		log.WithFields("var", envVar).Trace("loading bytes from environment")

		value := os.Getenv(envVar)
		if value == "" {
			return nil, fmt.Errorf("no key found in environment variable %q", envVar)
		}

		keyBytes, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return nil, err
		}
		return keyBytes, nil
	}

	// comes from the config...

	if _, err := os.Stat(path); err != nil {
		log.Trace("using bytes from config")

		decodedKey, err := base64.StdEncoding.DecodeString(path)
		if err != nil {
			return nil, fmt.Errorf("unable to base64 decode key: %w", err)
		}

		return decodedKey, nil
	}

	// comes from a file...

	log.WithFields("path", path).Trace("loading bytes from file")

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	return io.ReadAll(f)
}
