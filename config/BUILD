genrule(
    name = "aws-provider-yaml",
    cmd = " ".join([
        "touch credentials &&",  # TODO(EKF) actual credentials?
        "install -D credentials cmd/clusterctl/examples/aws/out/credentials &&",
        "install -D $(location //cmd/clusterctl/examples/aws:out/aws_manager_image_patch.yaml)",
        "cmd/clusterctl/examples/aws/out/aws_manager_image_patch.yaml &&",
        "$(location @io_k8s_sigs_kustomize//:kustomize) build config/default > $@",
    ]),
    srcs = [
               ":default/kustomization.yaml",
               "//cmd/clusterctl/examples/aws:out/aws_manager_image_patch.yaml",
           ] +
           glob(["crds/*.yaml"]) +
           glob(["rbac/*.yaml"]) +
           glob(["manager/*.yaml"]),
    tools = ["@io_k8s_sigs_kustomize//:kustomize"],
    outs = ["aws_provider.yaml"],
    visibility = ["//visibility:public"],
)
