package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
	"sigs.k8s.io/cluster-api-provider-aws/ci/ami/custom"
)

var (
	OWNER_REPO  = os.Getenv("GITHUB_REPOSITORY")
	OWNER, REPO = strings.Split(OWNER_REPO, "/")[0], strings.Split(OWNER_REPO, "/")[1]
)

func main() {
	var m2, m3 string
	url := "https://storage.googleapis.com/kubernetes-release/release/stable.txt"
	k8sReleaseResponse, err := http.Get(url)
	custom.CheckError(err)

	min1, err := ioutil.ReadAll(k8sReleaseResponse.Body)
	custom.CheckError(err)

	min1Release := custom.BuildReleaseVersion(string(min1))
	log.Print("Info: min1Release: Major ", min1Release.Major, ", Minor ", min1Release.Minor, ", Patch ", min1Release.Patch)

	if min1Release.Minor >= 2 {
		m2 = strconv.Itoa(min1Release.Major) + "." + strconv.Itoa(min1Release.Minor-1)
		m3 = strconv.Itoa(min1Release.Major) + "." + strconv.Itoa(min1Release.Minor-2)
	}

	url = fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/stable-%s.txt", m2)
	k8sReleaseResponse, err = http.Get(url)
	custom.CheckError(err)

	min2, err := ioutil.ReadAll(k8sReleaseResponse.Body)
	custom.CheckError(err)

	min2Release := custom.BuildReleaseVersion(string(min2))
	log.Print("Info: min2Release: Major ", min2Release.Major, ", Minor ", min2Release.Minor, ", Patch ", min2Release.Patch)

	url = fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/stable-%s.txt", m3)
	k8sReleaseResponse, err = http.Get(url)
	custom.CheckError(err)

	min3, err := ioutil.ReadAll(k8sReleaseResponse.Body)
	custom.CheckError(err)

	min3Release := custom.BuildReleaseVersion(string(min3))
	log.Print("Info: min3Release: Major ", min3Release.Major, ", Minor ", min3Release.Minor, ", Patch ", min3Release.Patch)

	latestAMIBuildConfig := &custom.AMIBuildConfig{
		K8sReleases: map[string]string{
			"min1": string(min1),
			"min2": string(min2),
			"min3": string(min3),
		},
	}

	latestAMIBuildConfigFileBytes, err := json.MarshalIndent(latestAMIBuildConfig, "", "  ")
	custom.CheckError(err)

	AMIBuildConfigFilename := os.Getenv("AMI_BUILD_CONFIG_FILENAME")
	dat, err := os.ReadFile(AMIBuildConfigFilename)
	if err != nil {
		if os.IsNotExist(err) {
			Action(latestAMIBuildConfigFileBytes, AMIBuildConfigFilename)
			log.Printf("Info: Created \"AMIBuildConfig.json\" K8s versions \"%s\"", latestAMIBuildConfig.K8sReleases)
			return
		} else {
			log.Fatal(err)
		}
	}

	currentAMIBuildConfig := new(custom.AMIBuildConfig)
	err = json.Unmarshal(dat, currentAMIBuildConfig)
	custom.CheckError(err)
	if !cmp.Equal(currentAMIBuildConfig, latestAMIBuildConfig) {
		prCreated := Action(latestAMIBuildConfigFileBytes, AMIBuildConfigFilename)
		if prCreated {
			log.Printf("Info: Updated \"%s\" with K8s versions from \"%s\" to \"%s\"", AMIBuildConfigFilename, currentAMIBuildConfig.K8sReleases, latestAMIBuildConfig.K8sReleases)
		}
	} else {
		log.Printf("Info: \"%s\" is up-to-date.", AMIBuildConfigFilename)
	}
}

func Action(blobBytes []byte, AMIBuildConfigFilename string) bool {
	// create a github api client and context using our action's auto-generated github token
	client, ctx := GetGithubClientCtx(os.Getenv("GITHUB_TOKEN"))

	// define references
	baseRef := "refs/heads/" + os.Getenv("CAPA_ACTION_BASE_BRANCH")
	headRef := "refs/heads/" + os.Getenv("CAPA_ACTION_HEAD_BRANCH")
	prHeadRef := OWNER + ":" + headRef
	prBaseRef := baseRef

	log.Print("Info: line 109 getref call STARTED")

	// check if the required head branch already exists
	ref, _, err := client.Git.GetRef(ctx, OWNER, REPO, headRef)
	log.Print("Info: line 109 getref call COMPLETED")

	if err == nil {
		prListOpts := github.PullRequestListOptions{
			Head: prHeadRef,
			Base: prBaseRef,
		}
		
		log.Print("Info: line 123 PullRequests call STARTED")

		prList, _, err := client.PullRequests.List(ctx, OWNER, REPO, &prListOpts)
		if err != nil {
			if len(prList) != 0 {
				log.Fatal(err)
			}
		}
		log.Print("Info: line 123 PullRequests call COMPLETED")

		if len(prList) == 0 {
			_, err := client.Git.DeleteRef(ctx, OWNER, REPO, headRef)
			custom.CheckError(err)

			CreateRef(client, ctx, baseRef, headRef)
			log.Printf("Info: Recreated existing head reference: %s", headRef)
		} else {
			log.Printf("Info: PR #%d corresponding to the specified base branch \"%s\" and head branch \"%s\" is still open. Exiting.", *prList[0].Number, baseRef, headRef)
			return false
		}
	} else {
		if ref == nil {
			log.Print("Info: line 143 CreateRef call STARTED")

			CreateRef(client, ctx, baseRef, headRef)
			log.Print("Info: line 143 CreateRef call COMPLETED")

		} else {
			log.Fatal(err)
		}
	}

	// get the reference to the head branch
	ref, _, err = client.Git.GetRef(ctx, OWNER, REPO, headRef)
	custom.CheckError(err)

	// get the commit pointed by the head branch
	parentCommit, _, err := client.Git.GetCommit(ctx, OWNER, REPO, *ref.Object.SHA)
	custom.CheckError(err)

	// upload the base64 encoded blob for updated amibuildconfig to github server
	blob, err := CreateBlob(client, ctx, "base64", blobBytes)
	custom.CheckError(err)

	// get the tree pointed by the head branch
	baseTree, _, err := client.Git.GetTree(ctx, OWNER, REPO, *parentCommit.Tree.SHA, true)
	custom.CheckError(err)

	// create a new tree with the updated amibuildconfig
	newTree, err := CreateTree(client, ctx, AMIBuildConfigFilename, "100644", *baseTree.SHA, *blob.SHA)
	custom.CheckError(err)

	// create a new commit with our newly created tree
	commitMsg := fmt.Sprintf("⚓️ Updating `%s`", AMIBuildConfigFilename)
	newCommit := github.Commit{
		Message: &commitMsg,
		Tree:    newTree,
		Parents: []*github.Commit{parentCommit},
	}
	commit, _, err := client.Git.CreateCommit(ctx, OWNER, REPO, &newCommit)
	custom.CheckError(err)

	// update the head to point to our newly created commit
	_, err = UpdateRef(client, ctx, ref, commit)
	custom.CheckError(err)

	// create pr to update the amibuildconfig
	prTitle := fmt.Sprintf("[CAPA-Action] ⚓️ Updating `%s`", AMIBuildConfigFilename)
	prBody := fmt.Sprintf("Updated config:\n```json\n%s\n```", string(blobBytes))
	prCreated, err := CreatePR(client, ctx, false, prTitle, prHeadRef, prBaseRef, prBody)
	custom.CheckError(err)

	// add labels to the newly created pr
	labels := []string{"ami-build-action"}
	_, _, err = client.Issues.AddLabelsToIssue(ctx, OWNER, REPO, *prCreated.Number, labels)
	custom.CheckError(err)

	// request reviewers for the newly created pr
	_, err = RequestReviewers(client, ctx, *prCreated.Number)
	custom.CheckError(err)

	// add assignees to the newly created pr
	_, _, err = client.Issues.AddAssignees(ctx, OWNER, REPO, *prCreated.Number, strings.Split(os.Getenv("CAPA_ACTION_PR_ASSIGNEES"), ","))
	custom.CheckError(err)

	return true
}

func GetGithubClientCtx(token string) (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc), ctx
}

func CreateRef(client *github.Client, ctx context.Context, fromRef, toRef string) *github.Reference {
	log.Print("Info: line 208 getref call STARTED")
	ref, _, err := client.Git.GetRef(ctx, OWNER, REPO, fromRef)
	log.Print("Info: line 208 getref call COMPLETED")

	custom.CheckError(err)

	newRef := github.Reference{
		Ref:    &toRef,
		URL:    ref.URL,
		Object: ref.Object,
	}
	
	log.Print("Info: line 232 CreateRef call STARTED")
	refNew, _, err := client.Git.CreateRef(ctx, OWNER, REPO, &newRef)
	log.Print("Info: line 232 CreateRef call COMPLETED")

	custom.CheckError(err)

	return refNew
}

func CreateBlob(client *github.Client, ctx context.Context, encoding string, blobBytes []byte) (*github.Blob, error) {
	blobContent := base64.RawStdEncoding.EncodeToString(blobBytes)
	newBlob := github.Blob{
		Content:  &blobContent,
		Encoding: &encoding,
	}
	blob, _, err := client.Git.CreateBlob(
		ctx,
		OWNER,
		REPO,
		&newBlob,
	)

	return blob, err
}

func CreateTree(client *github.Client, ctx context.Context, filename string, mode string, baseSHA, blobSHA string) (*github.Tree, error) {
	treePath := "ci/ami/" + filename
	treeMode := "100644"
	newTreeEntry := github.TreeEntry{
		Path: &treePath,
		Mode: &treeMode,
		SHA:  &blobSHA,
	}
	newTree, _, err := client.Git.CreateTree(ctx, OWNER, REPO, baseSHA, []*github.TreeEntry{&newTreeEntry})

	return newTree, err
}

func CreatePR(client *github.Client, ctx context.Context, prModify bool, prTitle, prHeadRef, prBaseRef, prBody string) (*github.PullRequest, error) {
	newPR := github.NewPullRequest{
		Title:               &prTitle,
		Head:                &prHeadRef,
		Base:                &prBaseRef,
		Body:                &prBody,
		MaintainerCanModify: &prModify,
	}
	prCreated, _, err := client.PullRequests.Create(ctx, OWNER, REPO, &newPR)

	return prCreated, err
}

func UpdateRef(client *github.Client, ctx context.Context, ref *github.Reference, commit *github.Commit) (*github.Reference, error) {
	refObjType := "commit"
	ref.Object.SHA = commit.SHA
	ref.Object.URL = commit.URL
	ref.Object.Type = &refObjType

	newRef, _, err := client.Git.UpdateRef(ctx, OWNER, REPO, ref, true)

	return newRef, err
}

func RequestReviewers(client *github.Client, ctx context.Context, prNum int) (*github.PullRequest, error) {
	reqReviewers := github.ReviewersRequest{
		Reviewers: strings.Split(os.Getenv("CAPA_ACTION_PR_REVIEWERS"), ","),
	}
	pr, _, err := client.PullRequests.RequestReviewers(ctx, OWNER, REPO, prNum, reqReviewers)

	return pr, err
}
