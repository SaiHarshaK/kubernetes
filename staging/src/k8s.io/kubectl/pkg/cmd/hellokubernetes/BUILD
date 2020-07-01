load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "hello_kubernetes.go",
    ],
    importmap = "k8s.io/kubernetes/vendor/k8s.io/kubectl/pkg/cmd/hellokubernetes",
    importpath = "k8s.io/kubectl/pkg/cmd/hellokubernetes",
    visibility = ["//visibility:public"],
    deps = [
        "//staging/src/k8s.io/apimachinery/pkg/apis/meta:go_default_library",
        "//staging/src/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//staging/src/k8s.io/cli-runtime/pkg/genericclioptions:go_default_library",
        "//staging/src/k8s.io/cli-runtime/pkg/resource:go_default_library",
        "//vendor/k8s.io/klog/v2:go_default_library",
        "//staging/src/k8s.io/kubectl/pkg/cmd/util:go_default_library",
        "//staging/src/k8s.io/kubectl/pkg/scheme:go_default_library",
        "//staging/src/k8s.io/kubectl/pkg/util/i18n:go_default_library",
        "//staging/src/k8s.io/kubectl/pkg/util/templates:go_default_library",
        "//vendor/github.com/spf13/cobra:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "hello_kubernetes_test.go",
    ],
    data = ["//staging/src/k8s.io/kubectl/testdata"],
    embed = [":go_default_library"],
    deps = [
        "//staging/src/k8s.io/cli-runtime/pkg/genericclioptions:go_default_library",
        "//staging/src/k8s.io/kubectl/pkg/cmd/testing:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
