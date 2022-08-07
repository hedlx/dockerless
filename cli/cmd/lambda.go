package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hedlx/doless/cli/ops"
	"github.com/hedlx/doless/cli/tui/lambda"
	"github.com/spf13/cobra"
)

var lambdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Lambda API methods",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var lambdaName string
var lambdaRuntime string
var lambdaType string

type lambdaOps struct {
	ctx context.Context
}

func (r *lambdaOps) List() tea.Cmd {
	return func() tea.Msg {
		lambdas, err := ops.ListLambdas(r.ctx)

		return lambda.LambdaListResponseMsg{
			Resp: &lambda.LambdaListResponse{
				Lambdas: lambdas,
				Err:     err,
			},
		}
	}
}

var lambdaCreateCmd = &cobra.Command{
	Use:   "create [path]",
	Short: "Create new lambda",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lambdaM := ops.CreateLambdaM{
			Name:       lambdaName,
			Runtime:    lambdaRuntime,
			LambdaType: lambdaType,
		}
		lambda, err := ops.CreateLambda(cmd.Context(), lambdaM, args[0])
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
	Short: "List lambdas",
	Run: func(cmd *cobra.Command, args []string) {
		m := &lambda.LambdaListModel{
			Lister: &lambdaOps{
				ctx: cmd.Context(),
			},
		}
		p := tea.NewProgram(lambda.InitLambdaListModel(m))

		if err := p.Start(); err != nil {
			fmt.Printf("Error: %s", err)
		}
	},
}

var lambdaStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start lambda",
	Run: func(cmd *cobra.Command, args []string) {
		lambda, err := ops.StartLambda(cmd.Context(), args[0])
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
	Short: "Deploy lambda, aka create + start",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lambdaM := ops.CreateLambdaM{
			Name:       lambdaName,
			Runtime:    lambdaRuntime,
			LambdaType: lambdaType,
		}

		lambda, err := ops.DeployLambda(cmd.Context(), lambdaM, args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		j, _ := json.MarshalIndent(lambda, "", "  ")
		fmt.Println(string(j))
	},
}

var lambdaDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy lambda",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := ops.DestroyLambda(cmd.Context(), args[0])
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
	lambdaCmd.AddCommand(lambdaDestroyCmd)

	lambdaCreateCmd.Flags().StringVarP(&lambdaName, "name", "n", "", "name")
	lambdaCreateCmd.Flags().StringVarP(&lambdaRuntime, "runtime", "r", "", "runtime")
	lambdaCreateCmd.Flags().StringVarP(&lambdaType, "type", "t", "ENDPOINT", "type of lambda (ENDPOINT | INTERNAL)")
	lambdaCreateCmd.MarkFlagRequired("name")
	lambdaCreateCmd.MarkFlagRequired("runtime")

	lambdaDeployCmd.Flags().StringVarP(&lambdaName, "name", "n", "", "name")
	lambdaDeployCmd.Flags().StringVarP(&lambdaRuntime, "runtime", "r", "", "runtime")
	lambdaDeployCmd.Flags().StringVarP(&lambdaType, "type", "e", "ENDPOINT", "type of lambda (ENDPOINT | INTERNAL)")
	lambdaDeployCmd.MarkFlagRequired("name")
	lambdaDeployCmd.MarkFlagRequired("runtime")
}
