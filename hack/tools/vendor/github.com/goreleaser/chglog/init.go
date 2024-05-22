package chglog

import (
	"fmt"
	"os"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func versionsInRepo(gitRepo *git.Repository) (map[plumbing.Hash]*semver.Version, error) {
	tagRefs, err := gitRepo.Tags()
	if err != nil {
		return nil, err
	}

	defer tagRefs.Close()

	tags := make(map[plumbing.Hash]*semver.Version)

	err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		var (
			version *semver.Version
			tag     *object.Tag
		)

		tagName := t.Name().Short()
		hash := t.Hash()

		if version, err = semver.NewVersion(tagName); err != nil || version == nil {
			fmt.Fprintf(os.Stderr, "Warning: unable to parse version from tag: %s : %v\n", tagName, err)
			return nil
		}

		// If this is an annotated tag look up the hash of the commit and use that.
		if tag, err = gitRepo.TagObject(t.Hash()); err == nil {
			var c *object.Commit

			if c, err = tag.Commit(); err != nil {
				return fmt.Errorf("cannot dereference annotated tag: %s : %w", tagName, err)
			}
			hash = c.Hash
		}

		tags[hash] = version

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tags, nil
}

func versionsOnBranch(gitRepo *git.Repository) (map[*semver.Version]plumbing.Hash, error) {
	repoVersions, err := versionsInRepo(gitRepo)
	if err != nil {
		return nil, err
	}

	refs, err := gitRepo.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	defer refs.Close()

	versions := make(map[*semver.Version]plumbing.Hash)

	err = refs.ForEach(func(c *object.Commit) error {
		if v, ok := repoVersions[c.Hash]; ok {
			versions[v] = c.Hash
		}
		return nil
	})

	return versions, err
}

// InitChangelog create a new ChangeLogEntries from a git repo.
func InitChangelog(gitRepo *git.Repository, owner string, notes *ChangeLogNotes, deb *ChangelogDeb, useConventionalCommits bool) (cle ChangeLogEntries, err error) {
	var start, end plumbing.Hash

	cle = make(ChangeLogEntries, 0)
	end = plumbing.ZeroHash

	versions, err := versionsOnBranch(gitRepo)
	if err != nil {
		return nil, err
	}

	tags := make([]*semver.Version, 0, len(versions))
	for v := range versions {
		tags = append(tags, v)
	}

	sort.Slice(tags, func(i, j int) bool { return tags[i].LessThan(tags[j]) })

	for _, version := range tags {
		var (
			commits      []*object.Commit
			commitObject *object.Commit
		)

		if version.Prerelease() != "" {
			// Do not need change logs for pre-release entries
			continue
		}

		start = versions[version]

		if commitObject, err = gitRepo.CommitObject(start); err != nil {
			return nil, fmt.Errorf("unable to fetch commit from tag %v: %w", version.Original(), err)
		}

		if owner == "" {
			owner = fmt.Sprintf("%s <%s>", commitObject.Committer.Name, commitObject.Committer.Email)
		}
		if commits, err = CommitsBetween(gitRepo, start, end); err != nil {
			return nil, fmt.Errorf("unable to find commits between %s & %s: %w", end, start, err)
		}

		changelog := CreateEntry(commitObject.Committer.When, version, owner, notes, deb, commits, useConventionalCommits)
		cle = append(cle, changelog)
		end = start
	}

	sort.Sort(sort.Reverse(cle))

	return cle, nil
}
