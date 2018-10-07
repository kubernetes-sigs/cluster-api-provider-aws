package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

const (
	Failed  = "FAILED"
	Pending = "PENDING"
	Created = "CREATED"
	Stopped = "STOPPED"
	Deleted = "DELETED"
)

var Enc = NewTracer()

type tracer struct {
	enc zapcore.Encoder
}

func NewTracer() *tracer {
	return &tracer{
		enc: zapcore.NewJSONEncoder(zapcore.EncoderConfig{}),
	}
}

func (t *tracer) ClusterMachine(msg string, cluster *clusterv1.Cluster, machine *clusterv1.Machine, status string) string {

	bytes, err := t.enc.EncodeEntry(zapcore.Entry{},
		[]zapcore.Field{
			zap.String("msg", msg),
			zap.String("cluster", cluster.Name),
			zap.String("machine", machine.Name),
			zap.String("status", status),
		},
	)

	if err != nil {
		return "ERROR: Couldn't encode"
	}

	return bytes.String()
}

func CMEnc(msg string, cluster *clusterv1.Cluster, machine *clusterv1.Machine, status string) string {
	return Enc.ClusterMachine(msg, cluster, machine, status)
}
