Automation of AMI Build/Publish Process
Summary
Current CAPA Kubernetes AMIs are being built manually when a new Kubernetes release happens and published to a privately-owned AWS account.
We want to design a system to automate building and publishing these AMIs which will get triggered when
a new Kubernetes release happens
a problem occurs (such as missing packages) with the previously built-AMIs
a new OS or CNI patch comes out that needs to be in the AMIs. 
Also, CAPA maintainers now have access to a CNCF-owned AWS account that will be used to host AMIs.

Motivation
Manually building Kubernetes AMIs for all the new releases and all supported OSs is a cumbersome and error-prone process. We can automate some/all parts of building/publishing AMIs via Prow and image-builder to make this process consistent, easier, and faster.
Goals
Automate building AMIs on every new Kubernetes release.
Automate re-building AMIs when there is an error on the previously-built AMIs.
Automate re-building AMIs to pick up a critical bug/security fix in an OS package (e.g. kernel, containerd, glibc, etc)
Automate publishing AMIs to a CNCF-owned AWS account each time new AMIs are built.

Non-goals / Future work
Automate everything in the process and over-engineer this.
Delete older versions of AMIs from the AWS account.
Keep only a single AMI per each patch/OS/arch in case multiple builds are triggered.


Proposal
Some aspects to brainstorm for automating AMI publishing process are as follows:
Which tools will be utilized for building/publishing?
Prow, bots, image-builder, slack bots?
We can use pr-creator for resting PRs. 
Security aspects
Deciding on where to store and how to access AWS keys is an important issue for the automation. BOSKOS will not be a viable option for Prow as it nukes AWS resources. Keeping access keys in some safe location(?) is one option, and  another option is to provide access keys to the maintainer group, which is less favorable.
What will trigger and who can trigger the build/publish job?  
A file can be used to trigger the build/publish job. Prow pre/post-submit jobs can be triggered on any changes on this file. The file below will contain all necessary information about which APIs we want to have in the AWS account. Any changes should trigger a job for example deleting a release, should clean the AMIs from the account.

Apiversion: 
DEFAULT_GO_VERSION: 1.13.5
DEFAULT_IMAGE_BUILDER_VERSION: 0.1.5
DEFAULT_OS_LIST:
 - centos-7
 - ubuntu-1804
 - amazon-2


BUILDS:
  - '1.18.8':
	GO_VERSION: 1.13.5
	K8S_RELEASE:1.18.8
      IMAGE_BUILDER_VERSION: 0.1.6
      BUILD_REVISION: 0
      ARCHITECTURE: x86-64
	EXTRA_CONFIG_FILE: 'extra-configs.yaml'
      OS:
       - centos-7
       - ubuntu-1804


  - '1.19.1':
	CONFIG: '1.16'
	GO_VERSION: 1.13.5
	K8S_RELEASE:1.19.1




We want to trigger build either if a new release is cut or if an existing AMI is erroneous.
If there is an error in an AMI, we have to wait for the next image-builder release to rebuild it (?), hence the above file will be good to be used for triggering a re-build for erroneous AMIs too.


A possible workflow:
A bot can be set up that will open an issue notifying that a new Kubernetes release is cut, then maintainers can do the necessary changes on the file above and create a PR. A prow job is run on the PR as pre-submit job and build the new AMIs. Once the PR is merged, a post-submit job will push the images to the AWS account.

Other automation tools to trigger a build: Slack bot?

Testing built images 
Testing newly-created AMIs before publishing may be worthwhile to avoid possible problems. After the build is complete, we can run conformance (or more light-weight smoke tests) tests using the new AMIs by bringing up a k8s cluster on the CNCF-owned AWS account. 
Adding newly-created and pushed AMIs to the docs
There are two options here too:
Manually, contributors create a PR for adding the list of images that are pushed to the AWS account to https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/master/docs/amis.md.
Another job/bot can automatically create a PR with the changes or just an issue with the list of the new AMIs created and pushed to the AWS account.

Transition steps from the current account to the new account
Can we move the default AWS account id in a minor release? Any concerns regarding its effect on existing clusters that uses AMIs on the current AWS account?
Existing clusters with older CAPA version will keep using the current account until they get updated. Do we keep pushing new AMIs to the private account also to continue supporting those Clusters?

How long should the images need to be in the private AWS account after the move is complete to the CNCF account?




Some more brainstorming for an example workflow (with input from @nikhita)
Create a github CAPA-bot similar to https://github.com/k8s-release-robot. 
Create amis.yaml to keep a list of Kubernetes releases and image-builder versions etc. like the one above.
A Kubernetes release is cut.
kubernetes/release team creates an issue in CAPA repo or periodically we watch for new issues created for new releases or a periodic job checks releases page for new releases.
After knowing that there is a new release, if a periodic job detects the new release, that job creates a PR using CAPA-bot that adds the new release to amis.yaml.
A different pre-submit build job is triggered by that PR, and detects only the changed part.
The build job builds AMIs with a staging tag which is not picked up by CAPA by default, then runs conformance tests on them, passing in the AMI ID manually. After tests pass, that PR needs to be merged manually.
If the PR is merged:
A new “published” AMI is created referencing the same EBS snapshot, which acts as “promotion” of the image.
Get all “published” images and publish JSON blob to S3. Javascript on website renders this as list of AMIs.
Rebuilding workflow can initially be more manual and triggered when amis.yaml is updated by maintainers.



Detecting when a new release is cut:
In K8s repo, release-bot creates an issue when a new release is cut (e.g., https://github.com/kubernetes/kubernetes/issues/92944). 
kubernetes/release uses publishing-bot to create this issue. (Logic is here: https://github.com/kubernetes/release/blob/master/lib/gitlib.sh)
→ We can ask the Kubernetes release team to have an issue created in CAPA too or we can have a periodic job that quickly checks if a release happened since the previous time it checked.


Automatically creating PRs for CAPA
We can automate creating PRs for i) Adding the new release to the amis.yaml, ii) Create a doc PR that has new created AMIs.
1- Setup a github CAPA-bot similar to: https://github.com/k8s-release-robot
2- In a pre-submit job, when a new issue is created using run-if-changed, make CAPA-bot create a PR that adds the new release to the list of releases file.
As a very relevant example, https://github.com/kubernetes/test-infra/tree/master/prow/cmd/autobump is run periodically, updates images with the new K8s release and creates a PR. autobump.sh checks if K8s version is changed and updates images if changed.

How to detect changes in the PR created by CAPA-bot?
Preferably, we want to create the AMIs just for the updated/added versions looking at the changes in amis.yaml. To do that, we need to have a way to get the diff of the PR and interpret it. There is no such example in other prow jobs.
This may be useful:
comparison=$(extract-commit "${old_version}")...$(extract-commit "${version}")
https://github.com/kubernetes/test-infra/compare/${comparison}

Separating Build and Publish jobs
There is no easy way to store AMIs created in one job that can be used by another job. So, an easier option may be to have a single job that builds, tests, and publishes.

Create issue for failed AMIs
If some AMIs do not pass conformance tests, CAPA-bot can create an issue listing the problematic AMIs.

Access Secrets
K8s-infra-prow-build-trusted Prow cluster is used to keep secret files. A very small group of people have access to that account. We can share CNCF-AWS account credentials with them. Or we can set the Secrets as environment variables in https://github.com/kubernetes/test-infra/blob/master/config/prow/config.yaml

Will use the same Prow cluster
There are only 3-clusters in Prow (k8s-infra-prow-build-trusted, k8s-infra-prow-build, and one for the rest). So, we will be using the same cluster we are using for testing to build AMIs.


Future Enhancements
Keeping only one set of AMIs for the latest patch of each Kubernetes release and cleaning up the patches that are no longer latest from the AWS account. 
Have a deletion policy for the AMIs and have another Prow job to delete out-dated images per deletion policy.
Relevant Issues
https://github.com/kubernetes-sigs/cluster-api-provider-aws/issues/1861
https://github.com/kubernetes/k8s.io/issues/1116
