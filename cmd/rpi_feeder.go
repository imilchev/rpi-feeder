package main

import (
	"github.com/imilchev/rpi-feeder/pkg/feeder"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "rpi-feeder",
		Short:        "Raspberry Pi automated feeder.",
		SilenceUsage: true,
	}
	cmd.AddCommand(newFeederCmd())

	return cmd
}

func newFeederCmd() *cobra.Command {
	var debug bool
	cmd := &cobra.Command{
		Use:          "start [configFilePath]]",
		Short:        "Starts the Raspberry Pi automated feeder.",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := initLogger(debug); err != nil {
				panic(err)
			}
			defer zap.S().Sync() //nolint

			fm, err := feeder.NewFeederManager(args[0])
			if err != nil {
				return err
			}

			return fm.Start()
		},
	}

	cmd.Flags().BoolVar(
		&debug,
		"debug",
		false,
		"Enable debug logging.")

	return cmd
}

func initLogger(enableDebug bool) error {
	var cfg zap.Config
	if enableDebug {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "console"
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	logger, err := cfg.Build()
	zap.ReplaceGlobals(logger)

	return err
}
