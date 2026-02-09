package repository

import "strings"

func isDuplicateKeyError(err error, constraintName string) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "duplicate key") && strings.Contains(errStr, constraintName)
}
