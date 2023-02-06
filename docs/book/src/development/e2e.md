# Developing E2E tests

Visit the [Cluster API documentation on E2E][cluster_api_e2e] for information on how to develop and run e2e tests. 

[cluster_api_e2e]: https://cluster-api.sigs.k8s.io/developer/e2e.html

## Set up

It's recommended to create a separate AWS account to run E2E tests. This ensures it does not conflict with
your other cluster API environment.

## Running from CLI

e2e tests can be run using Makefile targets:

```bash
$ make test-e2e
$ make test-e2e-eks
```

The following useful env variables can help to speed up the runs:

- `E2E_ARGS="--skip-cloudformation-creation --skip-cloudformation-deletion"` - in case the cloudformation stack is already properly set up, this ensures a quicker start and tear down.
- `GINKGO_FOCUS='\[PR-Blocking\]'` - only run a subset of tests
- `USE_EXISTING_CLUSTER` - use an existing management cluster (useful if you have a [Tilt][tilt-setup] setup)

[tilt-setup]: ./tilt-setup.md

## Running in IDEs

The following example assumes you run a management cluster locally (e.g. using [Tilt][tilt-setup]). 

[tilt-setup]: ./tilt-setup.md

### IntelliJ/GoLand

The following run configuration can be used:

```xml
<component name="ProjectRunConfigurationManager">
  <configuration default="false" name="capa e2e: unmanaged PR-Blocking" type="GoTestRunConfiguration" factoryName="Go Test">
    <module name="cluster-api-provider-aws" />
    <working_directory value="$PROJECT_DIR$/test/e2e/suites/unmanaged" />
    <parameters value="-ginkgo.focus=&quot;\[PR-Blocking\]&quot; -ginkgo.v=true -artifacts-folder=$PROJECT_DIR$/_artifacts --data-folder=$PROJECT_DIR$/test/e2e/data -use-existing-cluster=true -config-path=$PROJECT_DIR$/test/e2e/data/e2e_conf.yaml" />
    <envs>
      <env name="AWS_REGION" value="SET_AWS_REGION" />
      <env name="AWS_PROFILE" value="IF_YOU_HAVE_MULTIPLE_PROFILES" />
      <env name="AWS_ACCESS_KEY_ID" value="REPLACE_ACCESS_KEY" />
      <env name="AWS_SECRET_ACCESS_KEY" value="2W2RlZmFZSZnRg==" />
    </envs>
    <kind value="PACKAGE" />
    <package value="sigs.k8s.io/cluster-api-provider-aws/v2/test/e2e/suites/unmanaged" />
    <directory value="$PROJECT_DIR$" />
    <filePath value="$PROJECT_DIR$" />
    <framework value="gotest" />
    <pattern value="^\QTestE2E\E$" />
    <method v="2" />
  </configuration>
</component>
```

### Visual Studio Code

With the example above, you can configure a [launch configuration for VSCode][msft_vscode]. 

[msft_vscode]: https://go.microsoft.com/fwlink/?linkid=830387
