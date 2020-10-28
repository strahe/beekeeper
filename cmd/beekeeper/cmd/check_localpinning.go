package cmd

import (
	"github.com/ethersphere/beekeeper/pkg/bee"
	"github.com/ethersphere/beekeeper/pkg/check/localpinning"
	"github.com/ethersphere/beekeeper/pkg/random"

	"github.com/spf13/cobra"
)

func (c *command) initCheckLocalPinning() []*cobra.Command {
	const (
		optionNameDBCapacity = "db-capacity"
		optionNameFileName   = "file-name"
		optionNameDivisor    = "capacity-divisor"
		optionNameSeed       = "seed"
	)

	cmdChunk := &cobra.Command{
		Use:   "pin-chunk",
		Short: "Checks that a node on the cluster pins one chunk correctly.",
		Long:  "Checks that a node on the cluster pins one chunk correctly.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			cluster, err := bee.NewCluster(bee.ClusterOptions{
				APIScheme:               c.config.GetString(optionNameAPIScheme),
				APIHostnamePattern:      c.config.GetString(optionNameAPIHostnamePattern),
				APIDomain:               c.config.GetString(optionNameAPIDomain),
				APIInsecureTLS:          insecureTLSAPI,
				DebugAPIScheme:          c.config.GetString(optionNameDebugAPIScheme),
				DebugAPIHostnamePattern: c.config.GetString(optionNameDebugAPIHostnamePattern),
				DebugAPIDomain:          c.config.GetString(optionNameDebugAPIDomain),
				DebugAPIInsecureTLS:     insecureTLSDebugAPI,
				DisableNamespace:        disableNamespace,
				Namespace:               c.config.GetString(optionNameNamespace),
				Size:                    c.config.GetInt(optionNameNodeCount),
			})
			if err != nil {
				return err
			}

			var seed int64
			if cmd.Flags().Changed("seed") {
				seed = c.config.GetInt64(optionNameSeed)
			} else {
				seed = random.Int64()
			}

			return localpinning.CheckChunkFound(cluster, localpinning.Options{
				Seed: seed,
			})
		},
		PreRunE: c.checkPreRunE,
	}

	cmdChunk.Flags().Int(optionNameDBCapacity, 1000, "DB capacity in chunks")
	cmdChunk.Flags().Int(optionNameDivisor, 3, "divide store size by which value when uploading bytes")
	cmdChunk.Flags().Int64P(optionNameSeed, "s", 0, "seed for generating files; if not set, will be random")

	cmdBytes := &cobra.Command{
		Use:   "pin-bytes",
		Short: "Checks that a node on the cluster pins bytes correctly.",
		Long:  "Checks that a node on the cluster pins bytes correctly.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			cluster, err := bee.NewCluster(bee.ClusterOptions{
				APIScheme:               c.config.GetString(optionNameAPIScheme),
				APIHostnamePattern:      c.config.GetString(optionNameAPIHostnamePattern),
				APIDomain:               c.config.GetString(optionNameAPIDomain),
				APIInsecureTLS:          insecureTLSAPI,
				DebugAPIScheme:          c.config.GetString(optionNameDebugAPIScheme),
				DebugAPIHostnamePattern: c.config.GetString(optionNameDebugAPIHostnamePattern),
				DebugAPIDomain:          c.config.GetString(optionNameDebugAPIDomain),
				DebugAPIInsecureTLS:     insecureTLSDebugAPI,
				DisableNamespace:        disableNamespace,
				Namespace:               c.config.GetString(optionNameNamespace),
				Size:                    c.config.GetInt(optionNameNodeCount),
			})
			if err != nil {
				return err
			}

			var seed int64
			if cmd.Flags().Changed("seed") {
				seed = c.config.GetInt64(optionNameSeed)
			} else {
				seed = random.Int64()
			}

			return localpinning.CheckChunkFound(cluster, localpinning.Options{
				Seed: seed,
			})
		},
		PreRunE: c.checkPreRunE,
	}

	cmdBytes.Flags().Int(optionNameDBCapacity, 1000, "DB capacity in chunks")
	cmdBytes.Flags().Int(optionNameDivisor, 3, "divide store size by which value when uploading bytes")
	cmdBytes.Flags().Int64P(optionNameSeed, "s", 0, "seed for generating files; if not set, will be random")

	return []*cobra.Command{cmdChunk, cmdBytes}
}
