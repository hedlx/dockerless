package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hedlx/doless/cli/ops"
	"github.com/spf13/cobra"
)

var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Lambda API methods",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var lambdaID string
var lambdaName string
var lambdaRuntime string
var lambdaEndpoint string
var lambdaPath string

var lambdaCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create",
	Run: func(cmd *cobra.Command, args []string) {
		lambdaM := ops.CreateLambdaM{
			Name:     lambdaName,
			Runtime:  lambdaRuntime,
			Endpoint: lambdaEndpoint,
		}
		lambda, err := ops.CreateLambda(cmd.Context(), lambdaM, lambdaPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		j, _ := json.MarshalIndent(lambda, "", "  ")
		fmt.Println(string(j))
	},
}

var lambdaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List",
	Run: func(cmd *cobra.Command, args []string) {
		lambdas, err := ops.ListLambdas(cmd.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		j, _ := json.MarshalIndent(lambdas, "", "  ")
		fmt.Println(string(j))
	},
}

var lambdaStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start",
	Run: func(cmd *cobra.Command, args []string) {
		lambda, err := ops.StartLambda(cmd.Context(), lambdaID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		j, _ := json.MarshalIndent(lambda, "", "  ")
		fmt.Println(string(j))
	},
}

var lambdaDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy",
	Run: func(cmd *cobra.Command, args []string) {
		lambdaM := ops.CreateLambdaM{
			Name:     lambdaName,
			Runtime:  lambdaRuntime,
			Endpoint: lambdaEndpoint,
		}

		lambda, err := ops.DeployLambda(cmd.Context(), lambdaM, lambdaPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		j, _ := json.MarshalIndent(lambda, "", "  ")
		fmt.Println(string(j))
	},
}

var lambdaDestoyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy",
	Run: func(cmd *cobra.Command, args []string) {
		err := ops.DestroyLambda(cmd.Context(), lambdaID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
	},
}

func init() {
	RootCmd.AddCommand(lambdaCmd)
	lambdaCmd.AddCommand(lambdaDeployCmd)
	lambdaCmd.AddCommand(lambdaCreateCmd)
	lambdaCmd.AddCommand(lambdaListCmd)
	lambdaCmd.AddCommand(lambdaStartCmd)
	lambdaCmd.AddCommand(lambdaDestoyCmd)

	lambdaCreateCmd.Flags().StringVarP(&lambdaName, "name", "n", "", "name")
	lambdaCreateCmd.Flags().StringVarP(&lambdaRuntime, "runtime", "r", "", "runtime")
	lambdaCreateCmd.Flags().StringVarP(&lambdaEndpoint, "endpoint", "e", "", "endpoint")
	lambdaCreateCmd.Flags().StringVarP(&lambdaPath, "path", "p", "", "path")
	lambdaCreateCmd.MarkFlagRequired("name")
	lambdaCreateCmd.MarkFlagRequired("runtime")
	lambdaCreateCmd.MarkFlagRequired("endpoint")
	lambdaCreateCmd.MarkFlagRequired("path")

	lambdaCreateCmd.Flags().StringVarP(&lambdaID, "id", "i", "", "id")
	lambdaCreateCmd.MarkFlagRequired("id")

	lambdaStartCmd.Flags().StringVarP(&lambdaID, "id", "i", "", "id")
	lambdaStartCmd.MarkFlagRequired("id")

	lambdaDestoyCmd.Flags().StringVarP(&lambdaID, "id", "i", "", "id")
	lambdaDestoyCmd.MarkFlagRequired("id")

	lambdaDeployCmd.Flags().StringVarP(&lambdaName, "name", "n", "", "name")
	lambdaDeployCmd.Flags().StringVarP(&lambdaRuntime, "runtime", "r", "", "runtime")
	lambdaDeployCmd.Flags().StringVarP(&lambdaEndpoint, "endpoint", "e", "", "endpoint")
	lambdaDeployCmd.Flags().StringVarP(&lambdaPath, "path", "p", "", "path")
	lambdaDeployCmd.MarkFlagRequired("name")
	lambdaDeployCmd.MarkFlagRequired("runtime")
	lambdaDeployCmd.MarkFlagRequired("endpoint")
	lambdaDeployCmd.MarkFlagRequired("path")
}
