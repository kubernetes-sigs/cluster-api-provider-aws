# Triage Party

This is the deployment set up for the Triage Party instance
temporarily at http://triag-triag-1txo9za20rglu-177532865.us-east-1.elb.amazonaws.com/ using the AWS CDK.

## Deploying

Ask @randomvariable to deploy using any of the following

 * `cdk deploy`      deploy this stack to your default AWS account/region
 * `cdk diff`        compare deployed stack with current state
 * `cdk synth`       emits the synthesized CloudFormation template
 * `go test`         run unit tests

The `TRIAGE_PARTY_GITHUB_TOKEN` must be set to a Github token
with the `public_repo` scope
