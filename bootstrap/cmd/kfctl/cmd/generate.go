// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"

	kftypes "github.com/kubeflow/kubeflow/bootstrap/v3/pkg/apis/apps"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var generateCfg = viper.New()

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [all(=default)|k8s|platform]",
	Short: "Generate a kubeflow application where resources is one of 'platform|k8s|all'.",
	Long: `Generate a kubeflow application where resources is one of 'platform|k8s|all'.

  platform: non kubernetes resources (eg --platform gcp)
  k8s: kubernetes resources
  all: both platform and k8s

The default is 'all' for any selected platform.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetLevel(log.InfoLevel)
		log.Warn("please switch to new semantics")
		return errors.New("kfctl init has been deprecated")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCfg.SetConfigName("app")
	generateCfg.SetConfigType("yaml")

	// platform gcp
	generateCmd.Flags().String(string(kftypes.EMAIL), "",
		string(kftypes.EMAIL)+" if '--platform gcp'")
	bindErr := generateCfg.BindPFlag(string(kftypes.EMAIL), generateCmd.Flags().Lookup(string(kftypes.EMAIL)))
	if bindErr != nil {
		log.Errorf("couldn't set flag --%v: %v", string(kftypes.EMAIL), bindErr)
		return
	}

	// platform gcp
	generateCmd.Flags().String(string(kftypes.ZONE), "",
		string(kftypes.ZONE)+" if '--platform gcp'")
	bindErr = generateCfg.BindPFlag(string(kftypes.ZONE), generateCmd.Flags().Lookup(string(kftypes.ZONE)))
	if bindErr != nil {
		log.Errorf("couldn't set flag --%v: %v", string(kftypes.ZONE), bindErr)
		return
	}

	// platform gcp
	generateCmd.Flags().String(string(kftypes.IPNAME), "",
		string(kftypes.IPNAME)+" if '--platform gcp'")
	bindErr = generateCfg.BindPFlag(string(kftypes.IPNAME), generateCmd.Flags().Lookup(string(kftypes.IPNAME)))
	if bindErr != nil {
		log.Errorf("couldn't set flag --%v: %v", string(kftypes.IPNAME), bindErr)
		return
	}

	// platform gcp
	generateCmd.Flags().String(string(kftypes.HOSTNAME), "",
		string(kftypes.HOSTNAME)+" if '--platform gcp'")
	bindErr = generateCfg.BindPFlag(string(kftypes.HOSTNAME), generateCmd.Flags().Lookup(string(kftypes.HOSTNAME)))
	if bindErr != nil {
		log.Errorf("couldn't set flag --%v: %v", string(kftypes.HOSTNAME), bindErr)
		return
	}

	// platforms minikube
	generateCmd.Flags().Bool(string(kftypes.MOUNT_LOCAL), false,
		string(kftypes.MOUNT_LOCAL)+" if '--platform minikube'")
	bindErr = generateCfg.BindPFlag(string(kftypes.MOUNT_LOCAL), generateCmd.Flags().Lookup(string(kftypes.MOUNT_LOCAL)))
	if bindErr != nil {
		log.Errorf("couldn't set flag --%v: %v", string(kftypes.MOUNT_LOCAL), bindErr)
		return
	}

	// verbose output
	generateCmd.Flags().BoolP(string(kftypes.VERBOSE), "V", false,
		string(kftypes.VERBOSE)+" output default is false")
	bindErr = generateCfg.BindPFlag(string(kftypes.VERBOSE), generateCmd.Flags().Lookup(string(kftypes.VERBOSE)))
	if bindErr != nil {
		log.Errorf("couldn't set flag --%v: %v", string(kftypes.VERBOSE), bindErr)
		return
	}
}
