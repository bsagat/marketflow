package app

import (
	"marketflow/pkg/envzilla"
)

// SetConfig loads environment variables from the configuration file (.env)
// and sets them in the application’s runtime environment.
func SetConfig() error {
	err := envzilla.Loader("config/.env")
	if err != nil {
		return err
	}
	return nil
}

// Setup function sets connection to the adapters
func Setup() {

}
