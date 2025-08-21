/*
Copyright 2022 The Kubernetes Authors.

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

package registry

import (
	"fmt"
	"sort"
	"strings"

	"sigs.k8s.io/promo-tools/v4/types/image"
)

// Image holds information about an image. It's like an "Object" in the OOP
// sense, and holds all the information relating to a particular image that we
// care about.
type Image struct {
	Name image.Name `yaml:"name"`
	Dmap DigestTags `yaml:"dmap,omitempty"`
}

// Images is a slice of Image types.
type Images []Image

// A RegInvImage is a map containing all of the image names, and their
// associated digest-to-tags mappings. It is the simplest view of a Docker
// Registry, because the keys are just the image.Names (where each image.Name does
// *not* include the registry name, because we already key this by the
// RegistryName in MasterInventory).
//
// The image.Name is actually a path name, because it can be "foo/bar/baz", where
// the image name is the string after the last slash (in this case, "baz").
type RegInvImage map[image.Name]DigestTags

// DigestTags is a map where each digest is associated with a TagSlice. It is
// associated with a TagSlice because an image digest can have more than 1 tag
// pointing to it, even within the same image name's namespace (tags are
// namespaced by the image name).
type DigestTags map[image.Digest]TagSlice

// TagSlice is a slice of Tags.
type TagSlice []image.Tag

// TagSet is a set of Tags.
type TagSet map[image.Tag]interface{}

// ToYAML displays a RegInvImage as YAML, but with the map items sorted
// alphabetically.
func (a *RegInvImage) ToYAML(o YamlMarshalingOpts) string {
	images := a.ToSorted()

	var b strings.Builder
	for _, image := range images {
		fmt.Fprintf(&b, "- name: %s\n", image.Name)
		fmt.Fprintf(&b, "  dmap:\n")
		for _, digestEntry := range image.Digests {
			if o.BareDigest {
				fmt.Fprintf(&b, "    %s:", digestEntry.Hash)
			} else {
				fmt.Fprintf(&b, "    %q:", digestEntry.Hash)
			}

			switch len(digestEntry.Tags) {
			case 0:
				fmt.Fprintf(&b, " []\n")
			default:
				if o.SplitTagsOverMultipleLines {
					fmt.Fprintf(&b, "\n")
					for _, tag := range digestEntry.Tags {
						fmt.Fprintf(&b, "    - %s\n", tag)
					}
				} else {
					fmt.Fprintf(&b, " [")
					for i, tag := range digestEntry.Tags {
						if i == len(digestEntry.Tags)-1 {
							fmt.Fprintf(&b, "%q", tag)
						} else {
							fmt.Fprintf(&b, "%q, ", tag)
						}
					}
					fmt.Fprintf(&b, "]\n")
				}
			}
		}
	}

	return b.String()
}

// ToCSV is like ToYAML, but instead of printing things in an indented
// format, it prints one image on each line as a CSV. If there is a tag pointing
// to the image, then it is printed next to the image on the same line.
//
// Example:
// a@sha256:0000000000000000000000000000000000000000000000000000000000000000,a:1.0
// a@sha256:0000000000000000000000000000000000000000000000000000000000000000,a:latest
// b@sha256:1111111111111111111111111111111111111111111111111111111111111111,-
func (a *RegInvImage) ToCSV() string {
	images := a.ToSorted()

	var b strings.Builder
	for _, image := range images {
		for _, digestEntry := range image.Digests {
			if len(digestEntry.Tags) > 0 {
				for _, tag := range digestEntry.Tags {
					fmt.Fprintf(&b, "%s@%s,%s:%s\n",
						image.Name,
						digestEntry.Hash,
						image.Name,
						tag)
				}
			} else {
				fmt.Fprintf(&b, "%s@%s,-\n", image.Name, digestEntry.Hash)
			}
		}
	}

	return b.String()
}

// ToSorted converts a RegInvImage type to a sorted structure.
func (a *RegInvImage) ToSorted() []ImageWithDigestSlice {
	images := make([]ImageWithDigestSlice, 0)

	for name, dmap := range *a {
		var digests []Digest
		for k, v := range dmap {
			var tags []string
			for _, tag := range v {
				tags = append(tags, string(tag))
			}

			sort.Strings(tags)

			digests = append(digests, Digest{
				Hash: string(k),
				Tags: tags,
			})
		}
		sort.Slice(digests, func(i, j int) bool {
			return digests[i].Hash < digests[j].Hash
		})

		images = append(images, ImageWithDigestSlice{
			Name:    string(name),
			Digests: digests,
		})
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].Name < images[j].Name
	})

	return images
}

// ImageWithDigestSlice uses a slice of digests instead of a map, allowing its
// contents to be sorted.
type ImageWithDigestSlice struct {
	Name    string
	Digests []Digest
}

type Digest struct {
	Hash string
	Tags []string
}

// YamlMarshalingOpts holds options for tweaking the YAML output.
type YamlMarshalingOpts struct {
	// Render multiple tags on separate lines. I.e.,
	// prefer
	//
	//    sha256:abc...:
	//    - one
	//    - two
	//
	// over
	//
	//    sha256:abc...: ["one", "two"]
	//
	// If there is only 1 tag, it will be on one line in brackets (e.g.,
	// '["one"]').
	SplitTagsOverMultipleLines bool

	// Do not quote the digest. I.e., prefer
	//
	//    sha256:...:
	//
	// over
	//
	//    "sha256:...":
	//
	BareDigest bool
}
