package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hedlx/doless/cli/ops"
	"github.com/hedlx/doless/cli/tui/endpoint"
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

type endpointOps struct {
	ctx context.Context
}

func (r *endpointOps) Create(name string, path string, lambda string) tea.Cmd {
	return func() tea.Msg {
		ops.CreateEndpoint(r.ctx, &api.CreateEndpoint{
			Name:   name,
			Path:   path,
			Lambda: lambda,
		})
		return nil
	}
}

func (r *endpointOps) List() tea.Cmd {
	return func() tea.Msg {
		endpts, err := ops.ListEndpoints(r.ctx)

		return endpoint.EndpointListResponseMsg{
			Resp: &endpoint.EndpointListResponse{
				Endpoints: endpts,
				Err:       err,
			},
		}
	}
}

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
		m := &endpoint.EndpointListModel{
			Lister: &endpointOps{
				ctx: cmd.Context(),
			},
		}
		p := tea.NewProgram(endpoint.InitEndpointListModel(m))

		if err := p.Start(); err != nil {
			fmt.Printf("Error: %s", err)
		}
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
