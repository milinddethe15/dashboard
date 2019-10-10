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
	"fmt"
	"os"

	kftypes "github.com/kubeflow/kubeflow/bootstrap/v3/pkg/apis/apps"
	"github.com/kubeflow/kubeflow/bootstrap/v3/pkg/kfapp/coordinator"
	"github.com/kubeflow/kubeflow/bootstrap/v3/pkg/kfupgrade"
	"github.com/kubeflow/kubeflow/bootstrap/v3/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var applyCfg = viper.New()
var kfApp kftypes.KfApp
var err error

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply [all(=default)|k8s|platform]",
	Short: "Deploy a generated kubeflow application.",
	Long:  `Deploy a generated kubeflow application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.SetLevel(log.InfoLevel)
		if applyCfg.GetBool(string(kftypes.VERBOSE)) != true {
			log.SetLevel(log.WarnLevel)
		}
		if configFilePath != "" {
			kind, err := utils.GetObjectKindFromUri(configFilePath)
			if err != nil {
				return fmt.Errorf("Cannot determine the object kind: %v", err)
			}

			if kind == string(kftypes.KFDEF) {
				kfApp, err = coordinator.BuildKfAppFromURI(configFilePath)
			} else if kind == string(kftypes.KFUPGRADE) {
				kfUpgrade, err := kfupgrade.NewKfUpgrade(configFilePath)
				if err != nil {
					return fmt.Errorf("couldn't load KfUpgrade: %v", err)
				}

				err = kfUpgrade.Apply()
				if err != nil {
					return fmt.Errorf("couldn't apply KfUpgrade: %v", err)
				}
				return nil
			} else {
				return fmt.Errorf("Unsupported object kind: %v", kind)
			}
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot fetch current directory for apply: %v", err)
			}
			kfApp, err = coordinator.GetKfAppFromCfgFile(cwd+"/app.yaml", false)
			if err != nil || kfApp == nil {
				return fmt.Errorf("error loading kfapp: %v", err)
			}
		}
		if kfApp == nil {
			return fmt.Errorf("kfApp is nil")
		}
		applyErr := kfApp.Apply(kftypes.ALL)
		if applyErr != nil {
			return fmt.Errorf("couldn't apply KfApp: %v", applyErr)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCfg.SetConfigName("app")
	applyCfg.SetConfigType("yaml")

	// Config file option
	applyCmd.PersistentFlags().StringVarP(&configFilePath, string(kftypes.FILE), "f", "",
		`Static config file to use. Can be either a local path or a URL.
For example:
--config=https://raw.githubusercontent.com/kubeflow/kubeflow/master/bootstrap/config/kfctl_platform_existing.yaml
--config=kfctl_platform_gcp.yaml`)
	bindErr := applyCfg.BindPFlag(string(kftypes.CONFIG), applyCmd.Flags().Lookup(string(kftypes.CONFIG)))
	if bindErr != nil {
		log.Errorf("couldn't set flag --%v: %v", string(kftypes.CONFIG), bindErr)
		return
	}

	// verbose output
	applyCmd.Flags().BoolP(string(kftypes.VERBOSE), "V", false,
		string(kftypes.VERBOSE)+" output default is false")
	bindErr = applyCfg.BindPFlag(string(kftypes.VERBOSE), applyCmd.Flags().Lookup(string(kftypes.VERBOSE)))
	if bindErr != nil {
		log.Errorf("couldn't set flag --%v: %v", string(kftypes.VERBOSE), bindErr)
		return
	}
}
