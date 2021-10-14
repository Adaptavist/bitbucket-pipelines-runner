package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var dryRun bool
var cfgFile string
var variablesFlag []string
var secureVarsFlag []string

var rootCmd = &cobra.Command{
	Use:   "bpr",
	Short: "Cli for BitBucket pipelines API",
	Long:  `Trigger BitBucket pipelines remotely using this CLI so you can bring processes segmented tools together`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bpr.yaml)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry", false, "Dry run")
	rootCmd.PersistentFlags().StringSliceVar(&variablesFlag, "var", []string{}, "-var key=value")
	rootCmd.PersistentFlags().StringSliceVar(&secureVarsFlag, "secret", []string{}, "-s-var key=value")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigName(".bpr")
		viper.SetConfigType("env")
	}

	viper.SetEnvPrefix("BPR")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
