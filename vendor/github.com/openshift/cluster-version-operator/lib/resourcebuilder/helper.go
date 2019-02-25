package resourcebuilder

import "k8s.io/client-go/rest"

// withProtobuf makes a client use protobuf.
func withProtobuf(config *rest.Config) *rest.Config {
	config = rest.CopyConfig(config)
	config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	config.ContentType = "application/vnd.kubernetes.protobuf"
	return config
}
