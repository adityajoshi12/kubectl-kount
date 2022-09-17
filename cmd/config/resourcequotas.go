package config

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

type resourcequotaCountCmd struct {
	out    io.Writer
	errOut io.Writer
	opts   common.Options
}

func (c *resourcequotaCountCmd) Validate() error {

	if c.opts.AllNamespaces {
		c.opts.Namespace = ""
	}
	if c.opts.Selector != "" {
		common.ListOptions.LabelSelector = c.opts.Selector
	}
	return nil
}

func (c *resourcequotaCountCmd) Run() error {
	list, err := client.GetClient().CoreV1().ResourceQuotas(c.opts.Namespace).List(context.Background(), common.ListOptions)
	if err != nil {
		printer.Error(c.out, err.Error())
		os.Exit(1)
	}
	_, err = fmt.Fprintln(c.out, len(list.Items))
	return errors.Wrap(err, "write error")

}
func NewResourcequotaCountCmdCountCMD(out io.Writer, errout io.Writer) *cobra.Command {
	c := resourcequotaCountCmd{
		out:    out,
		errOut: errout,
	}
	var cmd = &cobra.Command{
		Use:     "resourcequota",
		Aliases: []string{"quota", "resourcequotas"},
		Example: "kubectl kount resourcequotas -n kube-system",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := c.Validate(); err != nil {
				return err
			}
			return c.Run()
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&c.opts.AllNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	f.StringVarP(&c.opts.Namespace, "namespace", "n", "default", "resource namespace")
	f.StringVarP(&c.opts.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2). Matching objects must satisfy all of the specified label constraints.")
	return cmd
}
