package cmd

import (
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"

	"github.com/twpayne/chezmoi/internal/chezmoi"
)

var verifyCmd = &cobra.Command{
	Use:     "verify [targets...]",
	Short:   "Exit with success if the destination state matches the target state, fail otherwise",
	Long:    mustGetLongHelp("verify"),
	Example: getExample("verify"),
	PreRunE: config.ensureNoError,
	RunE:    config.runVerifyCmd,
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	markRemainingZshCompPositionalArgumentsAsFiles(verifyCmd, 1)
}

func (c *Config) runVerifyCmd(cmd *cobra.Command, args []string) error {
	mutator := chezmoi.NewAnyMutator(chezmoi.NullMutator{})
	c.mutator = mutator

	persistentState, err := c.getPersistentState(&bolt.Options{
		ReadOnly: true,
	})
	if err != nil {
		return err
	}
	defer persistentState.Close()

	if err := c.applyArgs(args, persistentState); err != nil {
		return err
	}
	if mutator.Mutated() {
		return errExitFailure
	}
	return nil
}
