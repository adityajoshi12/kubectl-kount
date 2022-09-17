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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	countDesc = `kubectl plugin to count kubernetes object.`
)

// NewCmdCount adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func NewCmdCount() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:          "count",
		Short:        "Count Kubernetes resource instances.",
		Long:         countDesc,
		SilenceUsage: true,
	}

	log.SetLevel(log.DebugLevel)
	rootCmd.AddCommand(
		NewPodCountCmd(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()),
		NewDeploymentCountCmd(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()),
		NewStatefulSetCountCmd(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()),
		NewDaemonsetCountCmd(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()),
		NewReplicasetCountCmd(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()),
		NewJobCountCmd(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()),
		NewCronJobCountCmd(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()),

		NewServiceCountCmd(rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()),
	)
	return rootCmd
}
