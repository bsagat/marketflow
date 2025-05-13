package app

import (
	"fmt"
	"marketflow/pkg/envzilla"
	"os"
)

// SetConfig loads environment variables from the configuration file (.env)
// and sets them in the application’s runtime environment.
func SetConfig() error {
	err := envzilla.Loader("config/.env")
	if err != nil {
		return err
	}

	fmt.Println(os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	return nil
}

// Setup function sets connection to the adapters
func Setup() {

}
