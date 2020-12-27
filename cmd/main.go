package main

import (
	"fmt"
	changelog "github.com/masseelch/bbg-changelog-helper/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	flagNameSource = "source"
	flagNameTarget = "target"
)

var (
	root = &cobra.Command{
		Use:   "bbg-changelog-helper", // todo
		Short: "Command Line Tool to Assist in Maintaining a Full Change Log",
		Run: func(cmd *cobra.Command, args []string) {
			r, err := changelog.OpenRepository()
			if err != nil {
				log.Fatal(err)
			}

			// Retrieve the starting commit.
			sc, err := changelog.RetrieveCommit(r, viper.GetString(flagNameSource))
			if err != nil {
				log.Fatal(err)
			}

			// Retrieve the target commit.
			tc, err := changelog.RetrieveCommit(r, viper.GetString(flagNameTarget))
			if err != nil {
				log.Fatal(err)
			}

			// Changes needed to patch source to target.
			p, err := sc.Patch(tc)
			if err != nil {
				log.Fatalf("Could not compute diff: %s", err)
			}

			fmt.Println(p)
		},
	}
)

func init() {
	root.Flags().String(flagNameSource, "", "git commit hash or git tag to use as source") // todo - remove debug
	root.Flags().String(flagNameTarget, "", "git commit hash or git tag to use as target") // todo - remove debug

	_ = viper.BindPFlags(root.Flags())

	// _ = root.MarkFlagRequired(flagNameSource)
	// _ = root.MarkFlagRequired(flagNameTarget)
}

func main() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
