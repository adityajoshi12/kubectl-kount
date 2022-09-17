package nodes

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"kubectl-kount/client"
	common "kubectl-kount/cmd/common"
	"kubectl-kount/internal/printer"
	"os"
)

type nodesCountCmd struct {
	out    io.Writer
	errOut io.Writer
	opts   common.Options
}

func (c *nodesCountCmd) Validate() error {

	if c.opts.AllNamespaces {
		c.opts.Namespace = ""
	}
	if c.opts.Selector != "" {
		common.ListOptions.LabelSelector = c.opts.Selector
	}
	return nil
}

func (c *nodesCountCmd) Run() error {
	list, err := client.GetClient().CoreV1().Nodes().List(context.Background(), common.ListOptions)
	if err != nil {
		printer.Error(c.out, err.Error())
		os.Exit(1)
	}
	_, err = fmt.Fprintln(c.out, len(list.Items))
	return errors.Wrap(err, "write error")

}
func NewNodesCountCMD(out io.Writer, errout io.Writer) *cobra.Command {
	c := nodesCountCmd{
		out:    out,
		errOut: errout,
	}
	var cmd = &cobra.Command{
		Use:     "node",
		Aliases: []string{"no", "nodes"},
		Example: "kubectl kount nodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.Validate(); err != nil {
				return err
			}
			return c.Run()
		},
	}
	f := cmd.Flags()
	f.StringVarP(&c.opts.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2). Matching objects must satisfy all of the specified label constraints.")
	return cmd
}
