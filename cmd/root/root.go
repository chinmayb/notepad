package root

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	defaultPort = "8080"
)

func NewRootCommand() cobra.Command {
	rootCMd := cobra.Command{
		Use:   "app [sub]",
		Short: "My app command",
		// PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 	fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
		// },
		// PreRun: func(cmd *cobra.Command, args []string) {
		// 	fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
		// },
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Printf("Inside rootCmd Run with args: %v\n", args)
		// },
		// PostRun: func(cmd *cobra.Command, args []string) {
		// 	fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
		// },
		// PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// 	fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
		// },
		Run: func(cmd *cobra.Command, args []string) {
			fAddr := flag.CommandLine.Lookup("addr")
			fPort := flag.CommandLine.Lookup("port")
			fPort.Value.Set(defaultPort)
			if err := Register(fAddr.Value.String(), fPort.Value.String()); err != nil {
				fmt.Printf("error %v", err)
				os.Exit(1)
			}
		},
	}

	// rootCMd.SetArgs
	rootCMd.AddCommand()
	return rootCMd
}
