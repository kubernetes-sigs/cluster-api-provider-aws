genrule(
    name = "aws-provider-yaml",
    cmd = " ".join([
        "install -D $(location :credential_file) cmd/clusterctl/examples/aws/out/credentials &&",
        "install -D $(location //cmd/clusterctl/examples/aws:out/aws_manager_image_patch.yaml)",
        "cmd/clusterctl/examples/aws/out/aws_manager_image_patch.yaml &&",
        "$(location @io_k8s_sigs_kustomize//:kustomize) build config/default > $@",
    ]),
    srcs = [
               ":default/kustomization.yaml",
               ":credential_file",
               "//cmd/clusterctl/examples/aws:out/aws_manager_image_patch.yaml",
           ] +
           glob(["crds/*.yaml"]) +
           glob(["rbac/*.yaml"]) +
           glob(["manager/*.yaml"]),
    tools = ["@io_k8s_sigs_kustomize//:kustomize"],
    outs = ["aws_provider.yaml"],
    visibility = ["//visibility:public"],
)

genrule(
  name = "credential_file",
  outs = ["credentials"],
  cmd = " ".join([
      "touch $@ &&",
      "export AWS_CONFIG_FILE=$$(grep ^AWS_CONFIG_FILE bazel-out/volatile-status.txt | cut -f2 -d\" \") && ",
      "if [ -e $$AWS_CONFIG_FILE ]; then cat $$AWS_CONFIG_FILE >$@; fi"
  ]),
  stamp = 1,
)
