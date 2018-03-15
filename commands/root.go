// Package commands implements the Cobra root command for JiraCLI
package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// JiraCommand is the root command for JiraCLI
var JiraCommand = &cobra.Command{
	Use:   "jiracli",
	Short: "JiraShort",
	Long: `Jira Command Line Interface

Simple application for getting data from Jira from the command line.
`,

	Run: func(cc *cobra.Command, args []string) {
		readConfig()
		cc.HelpFunc()(cc, args)
	},
}

// init configures root cobra command and viper flags
func init() {

	// Flags for all commands
	JiraCommand.PersistentFlags().StringP("config", "c", "", "/path/to/config.toml")
	JiraCommand.PersistentFlags().StringP("path", "P", "/rest/api/2/", "API base path")
	JiraCommand.PersistentFlags().StringP("host", "H", "https://pondfive.atlassian.net", "Jira host base URL")
	JiraCommand.PersistentFlags().StringP("user", "u", "", "Jira user/email")
	JiraCommand.PersistentFlags().StringP("pass", "p", "", "Jira password")
	JiraCommand.PersistentFlags().IntP("limit", "l", 25, "Limit number of items in listings")
	JiraCommand.PersistentFlags().IntP("skip", "s", 0, "Skip first number of items in listings")
	JiraCommand.PersistentFlags().BoolP("verbose", "v", false, "Verbose text output")
	JiraCommand.PersistentFlags().BoolP("dump-response", "d", false, "Dump API response data")

	// Bind flags to toml compatible names
	viper.BindPFlag("config", JiraCommand.PersistentFlags().Lookup("config"))
	viper.BindPFlag("api.host", JiraCommand.PersistentFlags().Lookup("host"))
	viper.BindPFlag("api.path", JiraCommand.PersistentFlags().Lookup("path"))
	viper.BindPFlag("api.user", JiraCommand.PersistentFlags().Lookup("user"))
	viper.BindPFlag("api.pass", JiraCommand.PersistentFlags().Lookup("pass"))
	viper.BindPFlag("app.limit", JiraCommand.PersistentFlags().Lookup("limit"))
	viper.BindPFlag("app.skip", JiraCommand.PersistentFlags().Lookup("skip"))
	viper.BindPFlag("app.verbose", JiraCommand.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("dbg.dump-response", JiraCommand.PersistentFlags().Lookup("dump-response"))
}

// readConfig reads config from command line or fallback on default config
func readConfig() {

	// Read config file from --config flag, if present
	config := viper.GetString("config")
	if config != "" {
		viper.SetConfigFile(config)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read config file from command line: %v\n", err)
			os.Exit(1)
		}
		return
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.jiracli")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file, continuing without : %v", err)
	}
}

// resolveUserAndPass gets the user to input user and password if missing from settings
func resolveUserAndPass() {

	if viper.GetString("api.user") == "" || viper.GetString("api.pass") == "" {
		fmt.Println(`JiraCLI need credentials!

You can avoid entering password by creating a config.toml file and put it in "$HOME/.jiracli/".

Example config.toml:

	# Config file for JiraCLI
	[api]
	user = "<your_jira_email>"
	pass = "<your_jira_password>"

`)
	}

	r := bufio.NewReader(os.Stdin)

	if viper.GetString("api.user") == "" {
		fmt.Printf("Please enter Jira username > ")
		user, _ := r.ReadString('\n')
		viper.Set("api.user", user)
	}
	if viper.GetString("api.pass") == "" {

		fmt.Printf("Please enter Jira password (WARNING: chars will be echod!) > ")
		pass, _ := r.ReadString('\n')
		viper.Set("api.pass", pass)
	}

}

// verbose prints formatted string, if app.verbose setting is true
func verbose(f string, args ...interface{}) {
	if viper.GetBool("app.verbose") {
		log.Printf(f, args...)
	}
}
