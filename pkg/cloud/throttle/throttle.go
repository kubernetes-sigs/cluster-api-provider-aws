package throttle

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws/request"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/cloud/awserrors"
	"sigs.k8s.io/cluster-api-provider-aws/pkg/internal/rate"
)

type ServiceLimiters map[string]*ServiceLimiter

type ServiceLimiter []*OperationLimiter

func NewMultiOperationMatch(strs ...string) string {
	return "^" + strings.Join(strs, "|^")
}

type OperationLimiter struct {
	Operation  string
	RefillRate rate.Limit
	Burst      int
	regexp     *regexp.Regexp
	limiter    *rate.Limiter
}

func (o *OperationLimiter) Wait(r *request.Request) {
	o.getLimiter().Wait(r.Context())
}

func (o *OperationLimiter) Match(r *request.Request) (bool, error) {
	if o.regexp == nil {
		var err error
		o.regexp, err = regexp.Compile("^" + o.Operation)
		if err != nil {
			return false, err
		}
	}
	return o.regexp.Match([]byte(r.Operation.Name)), nil
}

func (s ServiceLimiter) LimitRequest(r *request.Request) {
	if ol, ok := s.matchRequest(r); ok {
		ol.Wait(r)
	}
}

func (o *OperationLimiter) getLimiter() *rate.Limiter {
	if o.limiter == nil {
		o.limiter = rate.NewLimiter(o.RefillRate, o.Burst)
	}
	return o.limiter
}

func (s ServiceLimiter) ReviewResponse(r *request.Request) {
	if r.Error != nil {
		if errorCode, ok := awserrors.Code(r.Error); ok {
			switch errorCode {
			case "Throttling", "RequestLimitExceeded":
				if ol, ok := s.matchRequest(r); ok {
					ol.limiter.ResetTokens()
				}
			}
		}
	}
}

func (s ServiceLimiter) matchRequest(r *request.Request) (*OperationLimiter, bool) {
	for _, ol := range s {
		match, err := ol.Match(r)
		if err != nil {
			return nil, false
		}
		if match {
			return ol, true
		}
	}
	return nil, false
}
