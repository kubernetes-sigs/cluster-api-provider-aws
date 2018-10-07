package machine

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func (a *Actuator) logEntry(msg string, cluster *clusterv1.Cluster, machine *clusterv1.Machine, status string) []zapcore.Field {
	return []zapcore.Field{
		zap.String("msg", msg),
		zap.String("cluster", cluster.Name),
		zap.String("machine", machine.Name),
		zap.String("status", status),
	}
}
