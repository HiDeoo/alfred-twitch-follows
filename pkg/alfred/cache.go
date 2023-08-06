package alfred

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
)

func GetCache(cache any) {
	if _, err := os.Stat(getCachePath()); err == nil {
		data, err := os.ReadFile(getCachePath())

		if err == nil {
			dec := json.NewDecoder(bytes.NewReader(data))
			dec.Decode(&cache)
		}
	}
}

func SetCache(cache any) error {
	if getCacheDirPath() == "" {
		return nil
	}

	err := ensureCacheDir()

	if err != nil {
		return err
	}

	data, err := json.Marshal(cache)

	if err != nil {
		return err
	}

	return os.WriteFile(getCachePath(), data, 0644)
}

func ensureCacheDir() error {
	if _, err := os.Stat(getCacheDirPath()); os.IsNotExist(err) {
		return os.MkdirAll(getCacheDirPath(), 0755)
	}

	return nil
}

func getCacheDirPath() string {
	return path.Join(os.Getenv("alfred_workflow_cache"))
}

func getCachePath() string {
	return path.Join(getCacheDirPath(), "cache.json")
}
