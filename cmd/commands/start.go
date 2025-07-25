package commands

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/inngest/inngest/cmd/commands/internal/localconfig"
	"github.com/inngest/inngest/pkg/authn"
	"github.com/inngest/inngest/pkg/config"
	"github.com/inngest/inngest/pkg/devserver"
	"github.com/inngest/inngest/pkg/headers"
	itrace "github.com/inngest/inngest/pkg/telemetry/trace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagGroup struct {
	name string
	fs   *pflag.FlagSet
}

func NewCmdStart(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "[Beta] Run Inngest as a single-node service.",
		Example: "inngest start",
		Run:     doStart,
	}

	groups := []FlagGroup{}

	baseFlags := pflag.NewFlagSet("base", pflag.ExitOnError)
	baseFlags.String("config", "", "Path to an Inngest configuration file")
	baseFlags.BoolP("help", "h", false, "Output this help information")
	baseFlags.String("host", "", "Inngest server hostname")
	baseFlags.StringP("port", "p", "8288", "Inngest server port")
	baseFlags.StringSliceP("sdk-url", "u", []string{}, "App serve URLs to sync (ex. http://localhost:3000/api/inngest)")
	baseFlags.String("signing-key", "", "Signing key used to sign and validate data between the server and apps.")
	baseFlags.StringSlice("event-key", []string{}, "Event key(s) that will be used by apps to send events to the server.")
	cmd.Flags().AddFlagSet(baseFlags)
	groups = append(groups, FlagGroup{name: "Flags:", fs: baseFlags})

	persistenceFlags := pflag.NewFlagSet("persistence", pflag.ExitOnError)
	persistenceFlags.String("sqlite-dir", "", "Directory for where to write SQLite database.")
	persistenceFlags.String("redis-uri", "", "Redis server URI for external queue and run state. Defaults to self-contained, in-memory Redis server with periodic snapshot backups.")
	persistenceFlags.String("postgres-uri", "", "PostgreSQL database URI for configuration and history persistence. Defaults to SQLite database.")
	cmd.Flags().AddFlagSet(persistenceFlags)
	groups = append(groups, FlagGroup{name: "Persistence Flags:", fs: persistenceFlags})

	advancedFlags := pflag.NewFlagSet("advanced", pflag.ExitOnError)
	advancedFlags.Int("poll-interval", 0, "Interval in seconds between polling for updates to apps")
	advancedFlags.Int("retry-interval", 0, "Retry interval in seconds for linear backoff when retrying functions - must be 1 or above")
	advancedFlags.Int("queue-workers", devserver.DefaultQueueWorkers, "Number of executor workers to execute steps from the queue")
	advancedFlags.Int("tick", devserver.DefaultTick, "The interval (in milliseconds) at which the executor polls the queue")
	advancedFlags.Int("connect-gateway-port", devserver.DefaultConnectGatewayPort, "Port to expose connect gateway endpoint")
	advancedFlags.Bool("no-ui", false, "Disable the web UI and GraphQL API endpoint")
	cmd.Flags().AddFlagSet(advancedFlags)
	groups = append(groups, FlagGroup{name: "Advanced Flags:", fs: advancedFlags})

	// Also add global flags
	groups = append(groups, FlagGroup{name: "Global Flags:", fs: rootCmd.PersistentFlags()})

	cmd.SetUsageFunc(func(c *cobra.Command) error {
		fmt.Printf("%s\n  %s\n\n%s\n%s\n\n",
			"Usage:",
			"inngest start [flags]",
			"Examples:",
			"inngest start",
		)

		for _, group := range groups {
			usage := group.fs.FlagUsages()

			help := ""
			if group.name != "" {
				help = help + group.name + "\n"
			}
			help = help + usage
			fmt.Println(help)
		}

		return nil
	})

	return cmd
}

func doStart(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	// TODO Likely need a `Start()`
	conf, err := config.Dev(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err = localconfig.InitStartConfig(ctx, cmd); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	port, err := strconv.Atoi(viper.GetString("port"))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	conf.EventAPI.Port = port
	conf.CoreAPI.Port = port

	host := viper.GetString("host")
	if host != "" {
		conf.EventAPI.Addr = host
		conf.CoreAPI.Addr = host
	}

	traceEndpoint := fmt.Sprintf("localhost:%d", port)
	if err := itrace.NewUserTracer(ctx, itrace.TracerOpts{
		ServiceName:   "tracing",
		TraceEndpoint: traceEndpoint,
		TraceURLPath:  "/dev/traces",
		Type:          itrace.TracerTypeOTLPHTTP,
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = itrace.CloseUserTracer(ctx)
	}()

	if err := itrace.NewSystemTracer(ctx, itrace.TracerOpts{
		ServiceName:   "tracing-system",
		TraceEndpoint: traceEndpoint,
		TraceURLPath:  "/dev/traces/system",
		Type:          itrace.TracerTypeOTLPHTTP,
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = itrace.CloseSystemTracer(ctx)
	}()

	tick := viper.GetInt("tick")
	if tick < 1 {
		tick = devserver.DefaultTick
	}

	signingKey := viper.GetString("signing-key")
	if signingKey == "" {
		fmt.Println("Error: signing-key is required")
		os.Exit(1)
	}
	_, err = authn.HashedSigningKey(signingKey)
	if err != nil {
		fmt.Printf("Error: signing-key must be a valid hexadecimal string\n")
		os.Exit(1)
	}

	eventKeys := viper.GetStringSlice("event-key")
	if len(eventKeys) == 0 {
		fmt.Println("Error: at least one event-key is required")
		os.Exit(1)
	}

	conf.ServerKind = headers.ServerKindCloud

	opts := devserver.StartOpts{
		Config:             *conf,
		ConnectGatewayHost: conf.CoreAPI.Addr,
		ConnectGatewayPort: viper.GetInt("connect-gateway-port"),
		EventKeys:          eventKeys,
		InMemory:           false,
		NoUI:               viper.GetBool("no-ui"),
		PollInterval:       viper.GetInt("poll-interval"),
		PostgresURI:        viper.GetString("postgres-uri"),
		QueueWorkers:       viper.GetInt("queue-workers"),
		RedisURI:           viper.GetString("redis-uri"),
		RequireKeys:        true,
		RetryInterval:      viper.GetInt("retry-interval"),
		SigningKey:         &signingKey,
		SQLiteDir:          viper.GetString("sqlite-dir"),
		Tick:               time.Duration(tick) * time.Millisecond,
		URLs:               viper.GetStringSlice("sdk-url"),
	}

	err = devserver.New(ctx, opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
