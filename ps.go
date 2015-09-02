// Copyright 2015 - António Meireles  <antonio.meireles@reformi.st>
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
//

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	psCmd = &cobra.Command{
		Use:     "ps",
		Aliases: []string{"status"},
		Short:   "lists running CoreOS instances",
		Run:     psCommand,
	}
)

func psCommand(cmd *cobra.Command, args []string) {
	viper.BindPFlags(cmd.Flags())
	ls, _ := ioutil.ReadDir(filepath.Join(SessionContext.configDir, "running"))
	if len(ls) > 0 {
		for _, d := range ls {
			vm, err := getSavedConfig(d.Name())
			if err == nil && vm.isAlive() {
				fmt.Printf("- %v (PID %v/detached=%v), up %v\n",
					vm.Name, vm.Pid, vm.Detached,
					time.Now().Sub(d.ModTime()))
				if buf, _ := ioutil.ReadFile(
					filepath.Join(SessionContext.configDir,
						fmt.Sprintf("running/%s/%s",
							d.Name(), "ip"))); buf != nil {
					fmt.Println("  - IP (public):",
						strings.TrimSpace(string(buf)))
				}
				if viper.GetBool("all") {
					pp, _ := json.MarshalIndent(vm, "", "    ")
					// FIXME get a PrettyPrint f
					fmt.Println(string(pp))
				}
			}
		}
	}
}

func getSavedConfig(uuid string) (VMInfo, error) {
	vm := VMInfo{}
	buf, err := ioutil.ReadFile(filepath.Join(SessionContext.configDir,
		fmt.Sprintf("running/%s/config", uuid)))
	if err != nil {
		return vm, err
	}
	json.Unmarshal(buf, &vm)
	if buf, err := ioutil.ReadFile(
		filepath.Join(SessionContext.configDir,
			fmt.Sprintf("running/%s/%s",
				vm.UUID, "ip"))); err == nil && buf != nil {
		vm.PublicIP = strings.TrimSpace(string(buf))
	}
	return vm, err
}

func init() {
	psCmd.Flags().BoolP("all", "a", false,
		"shows extended info about running instances")
	RootCmd.AddCommand(psCmd)
}
