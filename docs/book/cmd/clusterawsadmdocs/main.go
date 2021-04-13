/*
Copyright 2021 The Kubernetes Authors.

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

package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"sigs.k8s.io/cluster-api-provider-aws/cmd/clusterawsadm/cmd"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

type byName []*cobra.Command

func (s byName) Len() int           { return len(s) }
func (s byName) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byName) Less(i, j int) bool { return s[i].Name() < s[j].Name() }

type commandLeaf struct {
	name        string
	link        string
	subcommands map[string]commandLeaf
}

func main() {
	root := cmd.RootCmd()
	err := doc.GenMarkdownTree(root, "./src/clusterawsadm")
	if err != nil {
		log.Fatal(err)
	}

	tree := commandLeaf{
		name:        "clusterawsadm command reference",
		subcommands: make(map[string]commandLeaf),
		link:        "clusterawsadm/clusterawsadm",
	}
	buildCommandTree(tree, root)
	commandSummary(tree, 0)
}

func commandSummary(tree commandLeaf, prefix int) {
	prefixStr := strings.Repeat(" ", prefix) + "- "
	title := "[" + tree.name + "]"
	link := "(" + tree.link + ".md)"
	fmt.Println(prefixStr + title + link)
	for _, cmds := range tree.subcommands {
		commandSummary(cmds, prefix+2)
	}
}

func buildCommandTree(tree commandLeaf, cmd *cobra.Command) {

	children := cmd.Commands()

	sort.Sort(byName(children))
	for _, child := range children {
		if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
			continue
		}
		name := child.Name()
		components := strings.Split(name, " ")
		leaf := tree
		for _, c := range components {
			newLeaf, ok := leaf.subcommands[c]
			if !ok {
				newLeaf = commandLeaf{
					name:        c,
					subcommands: make(map[string]commandLeaf),
					link:        leaf.link + "_" + c,
				}
				leaf.subcommands[c] = newLeaf
			}
			leaf = newLeaf
			buildCommandTree(leaf, child)
		}
	}
}
