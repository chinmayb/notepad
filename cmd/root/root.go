package root

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultPort = "8080"
	defaultAddr = "0.0.0.0"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chaosctl",
	Short: "chaosctl ",
	Long:  ``,
}

var cfgFile string

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, _ := homedir.Dir()
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func NewRootCommand() *cobra.Command {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("addr", "a", defaultAddr, "addr")
	rootCmd.PersistentFlags().StringP("port", "p", defaultPort, "port")

	viper.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	rootCmd = &cobra.Command{
		Use:   "app [sub]",
		Short: "",
		// PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 	fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
		// },
		// PreRun: func(cmd *cobra.Command, args []string) {
		// 	fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
		// },

		// PostRun: func(cmd *cobra.Command, args []string) {
		// 	fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
		// },
		// PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// 	fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
		// },
		Run: func(cmd *cobra.Command, args []string) {

			if err := Register(viper.GetString("addr"), viper.GetString("port")); err != nil {
				fmt.Printf("error %v", err)
				os.Exit(1)
			}
		},
	}

	// rootCMd.SetArgs
	rootCmd.AddCommand()
	return rootCmd
}
