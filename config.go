package plugin

import (
	"github.com/spf13/viper"
	"strings"
)

func init() {
	viper.AutomaticEnv()          // load matched environment variables
	viper.SetEnvPrefix("CRAWLAB") // environment variable prefix as CRAWLAB
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
}
