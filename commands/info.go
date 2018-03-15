// Package commands info command
package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	infoCommand.PersistentFlags().Bool("include-password", false, "Display password in output")
	viper.BindPFlag("include-password", infoCommand.PersistentFlags().Lookup("include-password"))

	JiraCommand.AddCommand(infoCommand)
}

var infoCommand = &cobra.Command{
	Use:   "info",
	Short: "Print configuration settings",
	Long:  `JiraCLI Info Command`,

	Run: func(cc *cobra.Command, args []string) {
		readConfig()
		fmt.Printf("JiraCLI Information\n\n")
		printSettings(viper.GetBool("include-password"))
		fmt.Println("")
	},
}

// printSettings prints the persistent flags bound to viper
func printSettings(showPass bool) {
	fmt.Println("Settings:")

	fmt.Printf("  config: %q\n", viper.GetString("api.config"))
	fmt.Printf("  host: %q\n", viper.GetString("api.host"))
	fmt.Printf("  path: %q\n", viper.GetString("api.path"))
	fmt.Printf("  user: %q\n", viper.GetString("api.user"))
	if showPass {
		fmt.Printf("  pass: %q\n", viper.GetString("api.pass"))
	}
	fmt.Printf("  limit: %d\n", viper.GetInt("app.limit"))
	fmt.Printf("  skip: %d\n", viper.GetInt("app.skip"))
	fmt.Printf("  verbose: %v\n", viper.GetBool("app.verbose"))
	fmt.Printf("  dump-response: %v\n", viper.GetBool("dbg.dump-response"))

}
