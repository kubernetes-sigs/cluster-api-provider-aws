/*
Copyright 2019 The Kubernetes Authors.

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
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

/*
required parameters

https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line
github token for modifying release

# dependencies
* requires a build environment (make artifact-release must work)
* github.com/itchio/gothub
*/

const (
	// TODO(chuckha): figure this out based on directory name
	repository = "cluster-api-provider-aws"

	// TODO move these into config
	registry         = "gcr.io"
	managerImageName = "cluster-api-aws-controller"
	pullPolicy       = "IfNotPresent"
)

func main() {
	var remote, user, version string

	fs := flag.NewFlagSet("main", flag.ExitOnError)
	fs.Usage = documentation
	fs.StringVar(&remote, "remote", "origin", "name of the remote in local git")
	fs.StringVar(&user, "user", "kubernetes-sigs", "the github user/organization of the repo to release a version to")

	// TODO(chuckha): it would be ideal if we could release major/minor/patch and have it
	// automatically bump the latest tag git finds
	// until then, hard code it
	fs.StringVar(&version, "version", "", "the version number, for example v0.2.0")
	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	required("version", version)

	checkDependencies()

	cfg := config{
		version:     version,
		artifactDir: "out",
		artifacts: []string{
			"cluster-api-provider-aws-examples.tar",
			"clusterawsadm-darwin-amd64",
			"clusterawsadm-linux-amd64",
			"clusterctl-darwin-amd64",
			"clusterctl-linux-amd64",
		},
		registry:         fmt.Sprintf("%s/%s", registry, repository),
		imageName:        managerImageName,
		pullPolicy:       pullPolicy,
		githubRepository: repository,
		githubUser:       user,
		gitRemote:        remote,
	}

	logger := &stdoutlogger{}
	run := &runner{
		builder: makebuilder{
			registry:   cfg.registry,
			imageTag:   cfg.version,
			pullPolicy: cfg.pullPolicy,
			logger:     logger,
		},
		releaser: gothubReleaser{
			artifactsDir: cfg.artifactDir,
			user:         cfg.githubUser,
			repository:   cfg.githubRepository,
			logger:       logger,
		},
		tagger: git{
			repository: cfg.githubRepository,
			remote:     cfg.gitRemote,
			logger:     logger,
		},
		config: cfg,
		logger: logger,
	}

	if err := run.run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("The next steps are:")
	fmt.Println()
	// TODO(chuckha): automate writing the release notes
	fmt.Println("- Write the release notes. It is a manual process")
	// TODO(chuckha): something something docker container
	fmt.Println("- push the container images")
	// TODO(chuckha): send an email or at least print out the contents of an
	// email to
	fmt.Println("- Email kubernetes-dev@googlegroups.com to announce that a release happened")
}

func required(arg, val string) {
	if val == "" {
		fmt.Printf("%v is a required parameter\n", arg)
		os.Exit(1)
	}
}

func checkDependencies() {
	_, err := exec.LookPath("gothub")
	if err != nil {
		fmt.Println("Please install gothub:")
		fmt.Println()
		fmt.Println("    go get github.com/itchio/gothub")
		os.Exit(1)
	}

	if ght := os.Getenv("GITHUB_TOKEN"); len(ght) == 0 {
		fmt.Println("Please set the GITHUB_TOKEN environment variable.")
		fmt.Println("Read this guide for more information on creating Personal Access Tokens")
		fmt.Println()
		fmt.Println("https://help.github.com/en/articles/creating-a-personal-access-token-for-the-command-line")
		os.Exit(1)
	}
}

func documentation() {
	fmt.Println(`
This tool is designed to help automate the release process.

Usage:

	release [-remote <remote>] [-user <user>] -version <version>

Options:

	-remote
		The local name of the remote repository.
		The default is origin.

	-user
		The github username or organization of the repository to publish a release to.
		The default is kubernetes-sigs.

	-version
		The name of the version to release.

Examples:

1.
	origin is github.com/kubernetes-sigs/cluster-api-provider-aws
	local master points to origin/master branch

	To release:

		./release -remote origin -version v1.1.1

2.
	origin is github.com/myuser/cluster-api-provider-aws
	upstream is github.com/kubernetes-sigs/cluster-api-provider-aws
	local master points to upstream/master branch

	To release:

		./release -remote upstream -version v1.1.1

3.
	To test release

		./release -remote YOUR_FORK -user YOUR_GITHUB_USER_NAME -version v1.1.1`)
}

// config defines all configuration needed to get a release going
type config struct {
	// version is the version being released
	version string
	// artifacts are the list of artifacts to attach to the release
	artifacts []string
	// artifactsDir is the directory where all artifacts will be found
	artifactDir string
	// registry is the image registry where the container image will live
	registry string
	// imageName is the name of the container image
	imageName string
	// pullPolicy defines the pull policy of the manager in the provider-components
	pullPolicy string
	// githubRepository is the name of the repository on github https://github.com/<org or user>/<repository>
	githubRepository string
	// githubUser is the user/org name on github
	githubUser string
	// gitRemote is the local name of the remote to push the tag to
	gitRemote string
}

type runner struct {
	builder  builder
	releaser releaser
	tagger   tagger
	config   config
	logger   logger
}

// TODO sha512 the artifacts!
func (r runner) run() error {
	r.logger.Infof("tagging repository %q\n", r.config.version)
	if err := r.tagger.tag(r.config.version); err != nil {
		return err
	}
	r.logger.Infof("checking out tag %q\n", r.config.version)
	if err := r.tagger.checkout(r.config.version); err != nil {
		return err
	}
	r.logger.Infof("building artifacts %v\n", r.config.artifacts)
	if err := r.builder.build(); err != nil {
		return err
	}
	r.logger.Infof("building container image: %s/%s:%s\n", r.config.registry, r.config.imageName, r.config.version)
	if err := r.builder.images(); err != nil {
		return err
	}
	r.logger.Infof("pushing tag %q\n", r.config.version)
	if err := r.tagger.pushTag(r.config.version); err != nil {
		return err
	}
	r.logger.Infof("drafting a release for tag %q\n", r.config.version)
	if err := r.releaser.draft(r.config.version); err != nil {
		return err
	}
	for _, artifact := range r.config.artifacts {
		r.logger.Infof("uploading %q\n", artifact)
		if err := r.releaser.upload(r.config.version, artifact); err != nil {
			return err
		}
	}
	return nil
}

type releaser interface {
	draft(string) error
	upload(string, string) error
}

type gothubReleaser struct {
	// github user or organization (kubernetes-sigs, heptiolabs, chuckha, etc)
	user string
	// repository is the name of the repository
	repository string

	artifactsDir string
	logger       logger
}

func (g gothubReleaser) draft(version string) error {
	cmd := exec.Command("gothub", "release", "--tag", version, "--user", g.user, "--repo", g.repository, "--draft")
	out, err := cmd.CombinedOutput()
	if err != nil {
		g.logger.Info(string(out))
	}
	return err
}
func (g gothubReleaser) upload(version, file string) error {
	cmd := exec.Command("gothub", "upload", "--tag", version, "--user", g.user, "--repo", g.repository, "--file", path.Join(g.artifactsDir, file), "--name", file)
	out, err := cmd.CombinedOutput()
	if err != nil {
		g.logger.Info(string(out))
	}
	return err
}

type tagger interface {
	tag(string) error
	pushTag(string) error
	checkout(string) error
}

type git struct {
	// repository is the name of the repository
	repository string

	// remote is the local name of the remote
	remote string
	logger logger
}

func (g git) gitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		g.logger.Info(string(out))
	}
	return err
}

func (g git) getUserConfirmation(message string) {
	g.logger.Info(message)
	reader := bufio.NewReader(os.Stdin)
	g.logger.Info("Press return to confirm.")
	// throwing away the error
	reader.ReadString('\n')
}

func (g git) forceUpdateTag(version string) error {
	g.getUserConfirmation(fmt.Sprintf("You are about to force update the tag %v to reference HEAD.", version))
	return g.updateTag(version, "-f")
}

func (g git) updateTag(version string, extraArgs ...string) error {
	args := append([]string{"tag"}, extraArgs...)
	args = append(args, "-s", "-m", g.tagMessage(version), version)
	return g.gitCommand(args...)
}

func (g git) tagMessage(version string) string {
	return fmt.Sprintf("A release of %q for version %q", g.repository, version)
}

func (g git) tag(version string) error {
	oldRev, err := g.tagCommit(version)
	// an error getting the tag commit implies the tag does not exist
	if err != nil {
		if err := g.updateTag(version); err != nil {
			g.logger.Infof("failed to create a new tag: %v", version)
			return err
		}
		return nil
	}

	// if there is no error getting the tag commit that means the tag exists and must be force updated
	g.logger.Infof("Tag %q references commit %q\n", version, oldRev)
	if err := g.forceUpdateTag(version); err != nil {
		g.logger.Infof("failed to force update tag: %v", version)
		return err
	}
	return nil
}

func (g git) tagCommit(version string) (string, error) {
	cmd := exec.Command("git", "rev-parse", version)
	out, err := cmd.Output()
	if err != nil {
		// TODO print stderr as well
		return "", err
	}

	return strings.TrimSpace(string(out)), err
}

func (g git) remoteTagExists(version string) bool {
	cmd := exec.Command("git", "ls-remote", "--tags", "--exit-code", g.remote, version)
	// the output is unimportant for now. It will contain the ref of the tag if it exists.
	_, err := cmd.Output()
	return err == nil
}

func (g git) push(version string, extraArgs ...string) error {
	args := append([]string{"push"}, extraArgs...)
	args = append(args, g.remote, version)
	return g.gitCommand(args...)
}

func (g git) forcePush(version string) error {
	g.getUserConfirmation(fmt.Sprintf("You are about to force push the tag %v.", version))
	return g.push(version, "-f")
}

func (g git) pushTag(version string) error {
	if g.remoteTagExists(version) {
		if err := g.forcePush(version); err != nil {
			g.logger.Infof("failed to force push the tag %q", version)
			return err
		}
		return nil
	}

	if err := g.push(version); err != nil {
		g.logger.Infof("failed to push tag %q", version)
		return err
	}
	return nil
}

func (g git) checkout(version string) error {
	cmd := exec.Command("git", "checkout", version)
	out, err := cmd.CombinedOutput()
	if err != nil {
		g.logger.Info(string(out))
	}
	return err
}

type builder interface {
	build() error
	images() error
}

type makebuilder struct {
	registry   string
	imageTag   string
	pullPolicy string
	logger     logger
}

func (m makebuilder) cmdWithEnv(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("MANAGER_IMAGE_TAG=%v", m.imageTag),
		fmt.Sprintf("REGISTRY=%v", m.registry),
		fmt.Sprintf("PULL_POLICY=%v", m.pullPolicy))

	out, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.Info(string(out))
	}
	return err
}

func (m makebuilder) build() error {
	return m.cmdWithEnv("make", "release-artifacts")
}

func (m makebuilder) images() error {
	return m.cmdWithEnv("make", "docker-build")
}

type logger interface {
	Infof(string, ...interface{})
	Info(...interface{})
}

type stdoutlogger struct{}

func (s *stdoutlogger) Infof(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
}
func (s *stdoutlogger) Info(msgs ...interface{}) {
	fmt.Println(msgs...)
}
