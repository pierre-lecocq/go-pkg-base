package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnvFileIfSet(path string) error {
	if path == "" {
		return nil
	}

	st, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("-env: %w", err)
	}

	if st.IsDir() {
		return fmt.Errorf("-env: %q is a directory", path)
	}

	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("-env: load %q: %w", path, err)
	}

	return nil
}

func ValidatePresenceOf(names ...string) error {
	var missing = []string{}

	for _, name := range names {
		_, found := os.LookupEnv(name)
		if !found {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}

	return nil
}
