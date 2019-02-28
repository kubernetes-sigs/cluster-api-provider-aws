package framework

import "time"

const (
	WorkerRoleLabel = "node-role.kubernetes.io/worker"
	WaitShort       = 1 * time.Minute
	WaitMedium      = 3 * time.Minute
	WaitLong        = 10 * time.Minute
	RetryMedium     = 5 * time.Second
)
