package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"github.com/porter-dev/switchboard/pkg/drivers/helm"
	"github.com/porter-dev/switchboard/pkg/drivers/kubernetes"
	"github.com/porter-dev/switchboard/pkg/drivers/terraform"
	"github.com/porter-dev/switchboard/pkg/parser"
	"github.com/porter-dev/switchboard/pkg/types"
	"github.com/porter-dev/switchboard/pkg/worker"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var Version string = "dev"

var rootCmd = &cobra.Command{
	Use: "switchboard",
}

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}

var applyCmd = &cobra.Command{
	Use:  "apply [file]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := zerolog.New(zerolog.NewConsoleWriter())

		err := apply(args, &logger)

		if err != nil {
			logger.Err(err).Send()
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd, versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		color.New(color.FgRed).Println(err)
		os.Exit(1)
	}
}

func apply(args []string, logger *zerolog.Logger) error {
	filepath := args[0]

	fileBytes, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}

	resGroup, err := parser.ParseRawBytes(fileBytes)

	if err != nil {
		return err
	}

	basePath, err := os.Getwd()

	if err != nil {
		return err
	}

	worker := worker.NewWorker()
	worker.RegisterDriver("helm", helm.NewHelmDriver)
	worker.RegisterDriver("kubernetes", kubernetes.NewKubernetesDriver)
	worker.RegisterDriver("terraform", terraform.NewTerraformDriver)
	worker.SetDefaultDriver("helm")

	return worker.Apply(resGroup, &types.ApplyOpts{
		BasePath: basePath,
	})
}
