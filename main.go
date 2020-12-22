package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

var (
	verbose = false
	debug   = false

	emailConfig *viper.Viper
)

func init() {
	flag.IntP("port", "p", 8080, "port to start the server at")
	flag.String("itemDir", "items/", "Directory containing item files")
	flag.String("siteDir", "", "Serve directory from the HTTP server")
	flag.Bool("interactive", false, "Interactively enter mail credentials")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Print verbose information")
	flag.BoolVarP(&debug, "debug", "d", false, "Run program in debug mode")
}

func main() {
	flag.Parse()
	viper.BindPFlags(flag.CommandLine)

	err := setupConfig()
	if err != nil {
		fmt.Printf("Failed to read config: %s\n", err)
		return
	}

	if viper.GetBool("Interactive") {
		err = readMailFromStdIn()
		if err != nil {
			fmt.Printf("Failed to read mail credentials: %s\n", err)
		}
	}

	fmt.Println(startServer())
}

func setupConfig() error {
	viper.SetDefault("Port", 8080)

	viper.SetDefault("itemDir", "items/")
	viper.SetDefault("Interactive", false)

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/woodnthings")
	viper.AddConfigPath(".")

	emailConfig = viper.New()
	emailConfig.SetConfigName("emailAuth")
	emailConfig.AddConfigPath("$HOME/.config")

	err := emailConfig.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to read email config, defaulting to no email\n")
	}

	return viper.ReadInConfig()
}

func readMailFromStdIn() error {
	email := ""

	fmt.Printf("Email: ")
	fmt.Fscanf(os.Stdin, "%s", &email)
	fmt.Printf("Password: ")
	pass, err := terminal.ReadPassword(0)
	fmt.Printf("\n")

	if err != nil {
		return err
	}

	emailConfig.Set("Email", strings.TrimSpace(email))
	emailConfig.Set("Password", strings.TrimSpace(string(pass)))
	return nil
}

func verbosePrint(str string, args ...interface{}) {
	if verbose {
		fmt.Printf(str+"\n", args...)
	}
}
