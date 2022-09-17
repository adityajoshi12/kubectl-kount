/*
Copyright © 2022 Aditya Joshi <adityaprakashjoshi1@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package network

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

type serviceCountCmd struct {
	out    io.Writer
	errOut io.Writer
	opts   common.Options
}

func (p *serviceCountCmd) Validate() error {

	if p.opts.AllNamespaces {
		p.opts.Namespace = ""
	}
	if p.opts.Selector != "" {
		common.ListOptions.LabelSelector = p.opts.Selector
	}
	return nil
}
func (p serviceCountCmd) Run() error {

	list, err := client.GetClient().CoreV1().Services(p.opts.Namespace).List(context.Background(), common.ListOptions)
	if err != nil {
		printer.Error(p.out, err.Error())
		os.Exit(1)
	}
	_, err = fmt.Fprintln(p.out, len(list.Items))
	return errors.Wrap(err, "write error")
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
