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
package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"io"
	"kubectl-kount/client"
	"strconv"
)

type daemonsetCountCmd struct {
	out    io.Writer
	errOut io.Writer
	opts   Options
}

func (p *daemonsetCountCmd) run() error {
	list, err := client.GetClient().AppsV1().DaemonSets(p.opts.namespace).List(context.Background(), listOptions)
	if err != nil {
		panic(err)
	}
	_, err = p.out.Write([]byte(strconv.Itoa(len(list.Items))))
	if err != nil {
		return err
	}
	return nil
}
func (p *daemonsetCountCmd) validate() error {

	if p.opts.allNamespaces {
		p.opts.namespace = ""
	}
	if p.opts.selector != "" {
		listOptions.LabelSelector = p.opts.selector
	}
	return nil
}
func NewDaemonsetCountCmd(out io.Writer, errOut io.Writer) *cobra.Command {

	c := daemonsetCountCmd{
		out:    out,
		errOut: errOut,
	}
	cmd :=
		&cobra.Command{
			Use:     "daemonset",
			Aliases: []string{"ds", "daemonsets"},
			Short:   "Count daemonsets in a namespace.",
			Long:    `Count daemonsets in a namespace, optionally filtered by a label.`,
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
