package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hedlx/doless/cli/ops"
	api "github.com/hedlx/doless/client"
	"github.com/spf13/cobra"
)

var endpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Endpoint API methods",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var endpointName string
var endpointLambdaID string

var endpointCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		endpoint, err := ops.CreateEndpoint(cmd.Context(), &api.CreateEndpoint{
			Name:   endpointName,
			Path:   args[0],
			Lambda: endpointLambdaID,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		j, _ := json.MarshalIndent(endpoint, "", "  ")
		fmt.Println(string(j))
	},
}

var endpointListCmd = &cobra.Command{
	Use:   "list",
	Short: "List",
	Run: func(cmd *cobra.Command, args []string) {
		endpoints, err := ops.ListEndpoints(cmd.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		j, _ := json.MarshalIndent(endpoints, "", "  ")
		fmt.Println(string(j))
	},
}

func init() {
	RootCmd.AddCommand(endpointCmd)
	endpointCmd.AddCommand(endpointCreateCmd)
	endpointCmd.AddCommand(endpointListCmd)

	endpointCreateCmd.Flags().StringVarP(&endpointName, "name", "n", "", "endpoint name")
	endpointCreateCmd.Flags().StringVarP(&endpointLambdaID, "lambda-id", "l", "", "lambda id")
	endpointCreateCmd.MarkFlagRequired("name")
	endpointCreateCmd.MarkFlagRequired("lambda-id")
}
