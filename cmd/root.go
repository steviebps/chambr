package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/steviebps/rein/internal/logger"
	rein "github.com/steviebps/rein/pkg"
	utils "github.com/steviebps/rein/utils"
)

var home string
var cfgFile string
var chamber string
var globalChamber = rein.Chamber{Toggles: map[string]*rein.Toggle{}, Children: []*rein.Chamber{}}

// Version the version of rein
var Version = "development"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               "rein",
	Short:             "Local and remote configuration management",
	Long:              `CLI for managing application configuration of local and remote JSON files`,
	DisableAutoGenTag: true,
	Version:           Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.ErrorString(fmt.Sprintf("Error while starting rein: %v", err))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	var err error
	home, err = homedir.Dir()
	if err != nil {
		logger.ErrorString(err.Error())
		os.Exit(1)
	}

	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "rein configuration file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" && utils.Exists(cfgFile) {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/.rein")
		viper.SetConfigName("rein")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {

		configFileUsed := viper.ConfigFileUsed()
		if configFileUsed == "" {
			logger.ErrorString(err.Error())
			os.Exit(1)
		}

		logger.ErrorString(fmt.Sprintf("Error reading config file: %v\n", configFileUsed))
	}
}

func configPreRun(cmd *cobra.Command, args []string) {
	var jsonFile io.ReadCloser
	var err error
	chamberFile := viper.GetString("chamber")

	validURL, url := utils.IsURL(chamberFile)
	if validURL {
		res, err := http.Get(url.String())

		if err != nil {
			logger.ErrorString(fmt.Sprintf("Error trying to GET this resource \"%v\": %v\n", chamberFile, err))
			os.Exit(1)
		}
		jsonFile = res.Body
		defer jsonFile.Close()
	} else {
		if !utils.Exists(chamberFile) {
			logger.ErrorString(fmt.Sprintf("Could not find file \"%v\"\n", chamberFile))
			os.Exit(1)
		}

		jsonFile, err = os.Open(chamberFile)
		if err != nil {
			logger.ErrorString(fmt.Sprintf("Could not open file \"%v\": %v\n", chamberFile, err))
			os.Exit(1)
		}
		defer jsonFile.Close()
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logger.ErrorString(fmt.Sprintf("Error reading file \"%v\": %v\n", chamberFile, err))
		os.Exit(1)
	}

	if err := json.Unmarshal(byteValue, &globalChamber); err != nil {
		logger.ErrorString(fmt.Sprintf("Error reading \"%v\": %v\n", chamberFile, err))
		os.Exit(1)
	}
}
