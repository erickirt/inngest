package localconfig

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/inngest/inngest/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitDevConfig(ctx context.Context, cmd *cobra.Command) error {
	if err := mapDevFlags(cmd); err != nil {
		return err
	}

	loadConfigFile(ctx, cmd)

	return nil
}

func InitStartConfig(ctx context.Context, cmd *cobra.Command) error {
	if err := mapStartFlags(cmd); err != nil {
		return err
	}

	loadConfigFile(ctx, cmd)

	return nil
}

func loadConfigFile(ctx context.Context, cmd *cobra.Command) {
	l := logger.StdlibLogger(ctx)

	// Automatially bind environment variables
	viper.SetEnvPrefix("INNGEST")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	configPath, _ := cmd.Flags().GetString("config")
	if configPath != "" {
		// User specified the config file so we'll use that
		viper.SetConfigFile(configPath)
	} else {
		// Don't need to specify the extension since Viper will try to load
		// various extensions (inngest.json, inngest.yaml, etc.)
		viper.SetConfigName("inngest")

		if cwd, err := os.Getwd(); err != nil {
			l.Warn("error getting current directory", "error", err)
		} else {
			// Walk up the directory tree looking for a config file
			dir := cwd
			for {
				viper.AddConfigPath(dir)

				parent := filepath.Dir(dir)
				if parent == dir {
					break
				}

				dir = parent
			}
		}

		if homeDir, err := os.UserHomeDir(); err != nil {
			l.Warn("error getting home directory", "error", err)
		} else {
			// Fallback to ~/.config/inngest
			viper.AddConfigPath(filepath.Join(homeDir, ".config", "inngest"))
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if configPath != "" {
			// User explicitly specified a config file but we couldn't read it
			log.Fatalf("Error reading config file: %v", err)
		}
	} else {
		l.Info("using config", "file", viper.ConfigFileUsed())
	}
}

// mapDevFlags binds the command line flags to the viper configuration
func mapDevFlags(cmd *cobra.Command) error {
	var err error
	err = errors.Join(err, viper.BindPFlag("host", cmd.Flags().Lookup("host")))
	err = errors.Join(err, viper.BindPFlag("no-discovery", cmd.Flags().Lookup("no-discovery")))
	err = errors.Join(err, viper.BindPFlag("no-poll", cmd.Flags().Lookup("no-poll")))
	err = errors.Join(err, viper.BindPFlag("poll-interval", cmd.Flags().Lookup("poll-interval")))
	err = errors.Join(err, viper.BindPFlag("port", cmd.Flags().Lookup("port")))
	err = errors.Join(err, viper.BindPFlag("retry-interval", cmd.Flags().Lookup("retry-interval")))
	err = errors.Join(err, viper.BindPFlag("queue-workers", cmd.Flags().Lookup("queue-workers")))
	err = errors.Join(err, viper.BindPFlag("tick", cmd.Flags().Lookup("tick")))
	err = errors.Join(err, viper.BindPFlag("sdk-url", cmd.Flags().Lookup("sdk-url")))
	err = errors.Join(err, viper.BindPFlag("connect-gateway-port", cmd.Flags().Lookup("connect-gateway-port")))
	err = errors.Join(err, viper.BindPFlag("in-memory", cmd.Flags().Lookup("in-memory")))

	return err
}

// mapStartFlags binds the command line flags to the viper configuration
func mapStartFlags(cmd *cobra.Command) error {
	var err error
	err = errors.Join(err, viper.BindPFlag("host", cmd.Flags().Lookup("host")))
	err = errors.Join(err, viper.BindPFlag("port", cmd.Flags().Lookup("port")))
	err = errors.Join(err, viper.BindPFlag("signing-key", cmd.Flags().Lookup("signing-key")))
	err = errors.Join(err, viper.BindPFlag("event-key", cmd.Flags().Lookup("event-key")))
	err = errors.Join(err, viper.BindPFlag("redis-uri", cmd.Flags().Lookup("redis-uri")))
	err = errors.Join(err, viper.BindPFlag("postgres-uri", cmd.Flags().Lookup("postgres-uri")))
	err = errors.Join(err, viper.BindPFlag("poll-interval", cmd.Flags().Lookup("poll-interval")))
	err = errors.Join(err, viper.BindPFlag("retry-interval", cmd.Flags().Lookup("retry-interval")))
	err = errors.Join(err, viper.BindPFlag("queue-workers", cmd.Flags().Lookup("queue-workers")))
	err = errors.Join(err, viper.BindPFlag("sdk-url", cmd.Flags().Lookup("sdk-url")))
	err = errors.Join(err, viper.BindPFlag("sqlite-dir", cmd.Flags().Lookup("sqlite-dir")))
	err = errors.Join(err, viper.BindPFlag("tick", cmd.Flags().Lookup("tick")))
	err = errors.Join(err, viper.BindPFlag("connect-gateway-port", cmd.Flags().Lookup("connect-gateway-port")))
	err = errors.Join(err, viper.BindPFlag("no-ui", cmd.Flags().Lookup("no-ui")))

	return err
}
