// Copyright © 2018 The Kubernetes Authors.
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

package config

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/shurcool/httpfs/vfsutil"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/template"
	"strings"
)

// RootCmd is the root of the `alpha bootstrap command`
func RootCmd() *cobra.Command {

	force := false

	newCmd := &cobra.Command{
		Use:   "config",
		Short: "configure a cluster",
		Long:  `Cluster configuration commands`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Help(); err != nil {
				return err
			}
			return nil
		},
	}

	newCmd.PersistentFlags().BoolVarP(&force, "force", "f", force, "Overwrite if files already exist")

	newCmd.AddCommand(InitCmd(&force))
	newCmd.AddCommand(NewOverlayCmd())
	newCmd.AddCommand(RenderCmd())
	newCmd.AddCommand(CreateCmd())

	return newCmd
}

// InitCmd is the root of the `alpha bootstrap command`
func InitCmd(force *bool) *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "init",
		Short: "initialize new cluster",
		Long:  `Create directory with editable manifests`,
		RunE: func(cmd *cobra.Command, args []string) error {

			empty, err := currentDirisEmpty()
			if err != nil {
				fmt.Println("Error checking directory.")
				return errors.WithStack(err)
			}

			if !empty && !*force {
				fmt.Println("cannot write to empty directory without --force flag")
				return errors.WithStack(err)
			}

			fmt.Println("Initializing this directory for cluster-api-provider-aws configuration using kustomize")

			fmt.Println("\nWriting ./clusterawsadm")
			ioutil.WriteFile(rootIdentifier, []byte{}, os.ModePerm)
			c, err := newConfiguration("")
			if err != nil {
				return errors.WithStack(err)
			}
			err = c.walk(c.hfs, "", "", false)
			if err != nil {
				return errors.Wrap(err, "cannot walk asset hierarchy and initialize directory")
			}

			fmt.Println("\nDirectory initialized for kustomize. Run the following to create a new overlay called test:")
			fmt.Println("\nclusterawsadm alpha config new-overlay test")

			return nil
		},
	}

	return newCmd
}

// NewOverlayCmd is the root of the `alpha bootstrap command`
func NewOverlayCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "new-overlay",
		Short: "create a new configuration overlay",
		Long:  `Create a new configuration overlay`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Printf("Error: requires overlay name as an argument\n\n")
				if err := cmd.Help(); err != nil {
					return err
				}
				os.Exit(200)
			}

			c, err := newConfiguration(args[0])
			if err != nil {
				return errors.WithStack(err)
			}
			fmt.Printf("Creating new overlay %s:\n\n", args[0])
			err = c.walk(c.templateFS, "/", c.envDir, true)
			if err != nil {
				return errors.Wrap(err, "cannot walk asset hierarchy and create overlay")
			}

			fmt.Println("\nNew overlay created. Edit cluster-configuration.yaml to customize your cluster.")
			fmt.Println("To render component configuration for clusterctl, run:")
			fmt.Printf("\nclusterawsadm alpha config render %s\n\n", args[0])

			return nil
		},
	}

	return newCmd
}

// RenderCmd is the root of the `alpha bootstrap command`
func RenderCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "render",
		Short: "render cluster files",
		Long:  `Render cluster files`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				fmt.Printf("Error: requires environment as an argument\n\n")
				if err := cmd.Help(); err != nil {
					return err
				}
				os.Exit(200)
			}

			c, err := newConfiguration(args[0])
			if err != nil {
				return errors.WithStack(err)
			}

			err = c.kustomizeBuild("provider-components", false)
			if err != nil {
				return errors.Wrap(err, "Could not generate provider components")
			}

			err = c.kustomizeBuild("clusters", true)
			if err != nil {
				return errors.Wrap(err, "Could not generate cluster configuration")
			}

			err = c.kustomizeBuild("machines", true)
			if err != nil {
				return errors.Wrap(err, "Could not generate machines configuration")
			}

			err = c.kustomizeBuild("addons", false)
			if err != nil {
				return errors.Wrap(err, "Could not generate addons configuration")
			}

			fmt.Println("\nYou can now run the following to create the cluster:")

			fmt.Printf(
				"\nclusterctl create --provider-components %s --cluster %s --machines %s --addons %s\n\n",
				c.componentPath("provider-components"),
				c.componentPath("clusters"),
				c.componentPath("machines"),
				c.componentPath("addons"),
			)

			return nil
		},
	}

	return newCmd
}

// CreateCmd is the root of the `alpha bootstrap command`
func CreateCmd() *cobra.Command {
	newCmd := &cobra.Command{
		Use:   "create",
		Short: "render cluster files",
		Long:  `Render cluster files`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				fmt.Printf("Error: requires environment as an argument\n\n")
				if err := cmd.Help(); err != nil {
					return err
				}
				os.Exit(200)
			}
			c, err := newConfiguration(args[0])
			if err != nil {
				return errors.WithStack(err)
			}

			c.createCluster(args[1:])

			return nil
		},
	}

	return newCmd
}

type configuration struct {
	// EnvName exposes the environment name
	EnvName      string
	hfs          http.FileSystem
	rootDir      string
	envDir       string
	kustomizeDir string
	envOutDir    string
	templateDir  string
	templateFS   http.FileSystem
}

const (
	overlaysDir          = "/overlays"
	overlaysTemplatesDir = "/overlays/template"
	kustomizeConfigDir   = "kustomize-config"
	rootIdentifier       = ".clusterawsadm"
)

func newConfiguration(env string) (*configuration, error) {
	rootDir, err := findCLIFile(".")
	if err != nil {
		return nil, errors.Wrap(err, "cannot find .clusterawsadm in this or any parent directory")
	}
	envDir := filepath.Join(rootDir, overlaysDir, env)
	envOutDir := filepath.Join(envDir, "out")
	kustomizeDir := filepath.Join(rootDir, kustomizeConfigDir)
	templateDir := filepath.Join(rootDir, overlaysTemplatesDir)
	templateFS := http.Dir(templateDir)

	return &configuration{
		EnvName:      env,
		hfs:          assets,
		rootDir:      rootDir,
		envDir:       envDir,
		kustomizeDir: kustomizeDir,
		envOutDir:    envOutDir,
		templateDir:  templateDir,
		templateFS:   templateFS,
	}, nil
}

func (c configuration) ensureOverlay() {
	os.MkdirAll(c.envDir, os.ModePerm)
}

func (c configuration) walk(hfs http.FileSystem, stripDir string, replaceDir string, templated bool) error {

	w := walkConfiguration{
		hfs:        hfs,
		stripDir:   stripDir,
		replaceDir: replaceDir,
		templated:  templated,
		c:          c,
	}

	err := vfsutil.Walk(w.hfs, "/", w.walkFn())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type walkConfiguration struct {
	// EnvName exposes the environment name
	hfs        http.FileSystem
	stripDir   string
	replaceDir string
	templated  bool
	c          configuration
}

func (w walkConfiguration) walkFn() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		newPath := path
		if w.stripDir != "" {
			newPath = strings.Replace(path, w.stripDir, w.replaceDir+"/", 1)
		}
		err = w.writeFile(newPath, path, info)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
}

func (w walkConfiguration) writeFile(dst string, src string, info os.FileInfo) error {
	if info.IsDir() {
		return nil
	}
	fmt.Printf("Writing %s\n", dst)
	dir := filepath.Dir(dst)
	os.MkdirAll(filepath.Join(w.c.rootDir, dir), os.ModePerm)
	bytes, err := vfsutil.ReadFile(w.hfs, src)
	if err != nil {
		return errors.Wrapf(err, "cannot read asset: %s", src)
	}

	rendered := bytes

	if w.templated {
		renderedStr, err := template.Generate(src, string(bytes), w.c)
		if err != nil {
			return errors.Wrapf(err, "cannot render template %s", src)
		}
		rendered = []byte(renderedStr)
	}

	err = ioutil.WriteFile(filepath.Join(w.c.rootDir, dst), rendered, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "cannot write file: %s", dst)
	}
	return nil
}

func (c configuration) createCluster(args []string) error {

	componentArgs := []string{"create",
		"cluster",
		"--provider",
		"aws",
		"-a",
		c.componentPath("addons"),
		"-c",
		c.componentPath("clusters"),
		"-m",
		c.componentPath("machines"),
		"-p",
		c.componentPath("provider-components")}

	runCommandWithWait("clusterctl", append(componentArgs, args...)...)
	return nil
}

func runCommandWithWait(cmd string, args ...string) bool {
	command := runCommand(cmd, args...)
	if err := command.Wait(); err != nil {
		fmt.Println(err)
	}
	return command.ProcessState.Success()
}

func runShell(cmd string) *exec.Cmd {
	return runCommand("sh", "-c", cmd)
}

func runCommand(cmd string, args ...string) *exec.Cmd {
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Start(); err != nil {
		fmt.Println(err)
	}
	return command
}

func (c configuration) componentPath(component string) string {
	return filepath.Join(c.envOutDir, fmt.Sprintf("%s.yaml", component))
}

func (c configuration) kustomizeBuild(component string, filter bool) error {
	componentDir := filepath.Join(c.envDir, component)
	cmd := exec.Command("kustomize", "build", "-t", c.kustomizeDir, componentDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(string(stdout.Bytes()))
		fmt.Println(string(stderr.Bytes()))
		return errors.Wrap(err, "running Kustomize failed")
	}
	os.MkdirAll(c.envOutDir, os.ModePerm)
	outFile := c.componentPath(component)

	var outStr string

	if filter {
		outStr = filterYAML(string(stdout.Bytes()))
	} else {
		outStr = string(stdout.Bytes())
	}

	fmt.Printf("Writing %s\n", outFile)
	err = ioutil.WriteFile(outFile, []byte(outStr), os.ModePerm)
	if err != nil {
		fmt.Printf("cannot write file: %s\n", outFile)
		return errors.Wrapf(err, "cannot write file: %s", outFile)
	}

	return nil
}

func currentDirisEmpty() (bool, error) {
	cwd, err := os.Open(".")

	if err != nil {
		errors.WithStack(err)
	}

	entries, err := cwd.Readdir(1)

	if err != nil {
		errors.WithStack(err)
	}

	if len(entries) == 1 {
		return false, nil
	}
	return true, nil
}

func findCLIFile(path string) (string, error) {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return "", errors.WithStack(err)
	}

	if absPath == "/" {
		return "", errors.New("Reached / root of file system")
	}
	_, err = os.Stat(fmt.Sprintf("%s/.clusterawsadm", absPath))
	if err != nil {
		newPath := filepath.Clean(filepath.Join(absPath, ".."))
		return findCLIFile(newPath)
	}

	return path, nil
}

func filterYAML(docs string) string {
	return strings.Split(docs, "---\n")[1]
}
