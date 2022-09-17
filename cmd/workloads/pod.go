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
package workloads

import (
	"context"
	"fmt"
	"io"
	"kubectl-kount/client"
	common "kubectl-kount/cmd/common"
	"kubectl-kount/internal/printer"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
)

type podCountCmd struct {
	out    io.Writer
	errOut io.Writer
	opts   common.Options
}

func (p *podCountCmd) validate() error {

	if p.opts.AllNamespaces {
		p.opts.Namespace = ""
	}
	if p.opts.Selector != "" {
		common.ListOptions.LabelSelector = p.opts.Selector
	}
	return nil
}
func (p *podCountCmd) run() error {

	list, err := client.GetClient().CoreV1().Pods(p.opts.Namespace).List(context.Background(), common.ListOptions)
	if err != nil {
		printer.Error(p.out, err.Error())
		os.Exit(1)
	}
	if p.opts.Details {
		printTable(list)
	} else {
		_, err = fmt.Fprintln(p.out, len(list.Items))
	}
	return errors.Wrap(err, "write error")
}

func NewPodCountCmd(out io.Writer, errOut io.Writer) *cobra.Command {

	c := podCountCmd{
		out:    out,
		errOut: errOut,
	}
	cmd :=
		&cobra.Command{
			Use:     "pods",
			Aliases: []string{"po", "pod"},
			Short:   "Count pods in a namespace.",
			Long:    `Count pods in a namespace, optionally filtered by a label.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := c.validate(); err != nil {
					return err
				}
				return c.run()
			},
		}
	f := cmd.Flags()
	f.BoolVarP(&c.opts.AllNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current context is ignored even if specified with --namespace.")
	f.StringVarP(&c.opts.Namespace, "namespace", "n", "default", "resource namespace")
	f.StringVarP(&c.opts.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2). Matching objects must satisfy all of the specified label constraints.")
	f.BoolVarP(&c.opts.Details, "details", "d", false, "details")
	return cmd
}

func printTable(podList *v1.PodList) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	t.AppendHeader(table.Row{"Namespace", "Pod Name", "Container", "Container Name", "image", "Created"}, rowConfigAutoMerge)
	for _, v := range podList.Items {
		for _, container := range v.Spec.Containers {

			t.AppendRow(table.Row{v.Namespace, v.Name, len(v.Spec.Containers), container.Name, container.Image, v.CreationTimestamp}, rowConfigAutoMerge)
		}
	}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true, VAlign: text.VAlignMiddle},
		{Number: 2, AutoMerge: true, VAlign: text.VAlignMiddle},
		{Number: 3, AutoMerge: true},
		{Number: 4},
		{Number: 5},
		{Number: 6, AutoMerge: true},
	})
	t.Style().Options.SeparateRows = true
	t.Render()
}
