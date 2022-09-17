// Package cmd /*
package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"io"
	"kubectl-kount/client"
	"strconv"
)

type serviceCountCmd struct {
	out    io.Writer
	errOut io.Writer
	opts   Options
}

func (p *serviceCountCmd) validate() error {

	if p.opts.allNamespaces {
		p.opts.namespace = ""
	}
	if p.opts.selector != "" {
		listOptions.LabelSelector = p.opts.selector
	}
	return nil
}
func (p serviceCountCmd) run() error {

	list, err := client.GetClient().CoreV1().Services(p.opts.namespace).List(context.Background(), listOptions)
	if err != nil {
		panic(err)
	}

	_, err = p.out.Write([]byte(strconv.Itoa(len(list.Items))))
	if err != nil {
		return err
	}
	return nil
}

func NewServiceCountCmd(out io.Writer, errOut io.Writer) *cobra.Command {

	c := serviceCountCmd{
		out:    out,
		errOut: errOut,
	}
	cmd :=
		&cobra.Command{
			Use:     "service",
			Aliases: []string{"srv", "services"},
			Short:   "Count services in a namespace.",
			Long:    `Count services in a namespace, optionally filtered by a label.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := c.validate(); err != nil {
					return err
				}
				return c.run()
			},
		}
	f := cmd.Flags()
	f.BoolVarP(&c.opts.allNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	f.StringVarP(&c.opts.namespace, "namespace", "n", "default", "resource namespace")
	f.StringVarP(&c.opts.selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2). Matching objects must satisfy all of the specified label constraints.")
	return cmd
}
