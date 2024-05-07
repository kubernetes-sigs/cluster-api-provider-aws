# Nightly Builds

Nightly builds are regular automated builds of the CAPA source code that occur every night.

These builds are generated directly from the latest commit o source code on the main branch.

Nightly builds serve several purposes:

- **Early Testing**: They provide an opportunity for developers and testers to access the most recent changes in the codebase and identify any issues or bugs that may have been introduced.
- **Feedback Loop**: They facilitate a rapid feedback loop, enabling developers to receive feedback on their changes quickly, allowing them to iterate and improve the code more efficiently.
- **Preview of New Features**: Users and can get a preview of upcoming features or changes by testing nightly builds, although these builds may not always be stable enough for production use.

Overall, nightly builds play a crucial role in software development by promoting user testing, early bug detection, and rapid iteration.

CAPA Nightly build jobs run in Prow. For details on how this is configured you can check the [Periodics Jobs section](../topics/reference/jobs.md#periodics).

## Usage

To try a nightly build, you can download the latest built nightly CAPA manifests, you can find the available ones by executing the following command:
```bash
curl -sL -H 'Accept: application/json' "https://storage.googleapis.com/storage/v1/b/k8s-staging-cluster-api-aws/o" | jq -r '.items | map(select(.name | startswith("components/nightly_main"))) | .[] | [.timeCreated,.mediaLink] | @tsv'
```
The output should look something like this:
```
2024-05-03T08:03:09.087Z        https://storage.googleapis.com/download/storage/v1/b/k8s-staging-cluster-api-aws/o/components%2Fnightly_main_2024050x?generation=1714723389033961&alt=media
2024-05-04T08:02:52.517Z        https://storage.googleapis.com/download/storage/v1/b/k8s-staging-cluster-api-aws/o/components%2Fnightly_main_2024050y?generation=1714809772486582&alt=media
2024-05-05T08:02:45.840Z        https://storage.googleapis.com/download/storage/v1/b/k8s-staging-cluster-api-aws/o/components%2Fnightly_main_2024050z?generation=1714896165803510&alt=media
```

Now visit the link for the manifest you want to download. This will automatically download the manifest for you.

Once downloaded you can apply the manifest directly to your testing CAPI management cluster/namespace (e.g. with kubectl), as the downloaded CAPA manifest
will already contain the correct, corresponding CAPA nightly image reference.
