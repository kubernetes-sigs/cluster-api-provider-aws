package awserrors

// DryRunOperation is the error returned by AWS when a DryRun succeeds.
const DryRunOperation = "DryRunOperation"

// IsDryRunOperationError returns whether the error is DryRunOperation, which signifies the DryRun operation
// was successful.
func IsDryRunOperationError(err error) bool {
	if code, ok := Code(err); ok {
		switch code {
		case DryRunOperation:
			return true
		default:
			return false
		}
	}

	return false
}
