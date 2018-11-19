package errors

import (
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
)

const (
	MSG_TOKEN_EXPIRED_ERROR = "MSG_TOKEN_EXPIRED_ERROR"
)

// NonCriticalErrors is an slice of error statuses which are non-critical
// That means these errors can be passed to frontend as a warning
var NonCriticalErrors = []int32{http.StatusForbidden, http.StatusUnauthorized}

func HandleError(err error) ([]error, error) {
	nonCriticalErrors := make([]error, 0)
	return AppendError(err, nonCriticalErrors)
}

// AppendError will append non-critical error to slice (the first return value) and return the critical error
// as second value. We need to distinguish critical and non-critical errors because it is need to handle them
// in a different way
func AppendError(err error, nonCriticalErrors []error) ([]error, error) {
	if err != nil {
		if isErrorCritical(err) {
			return nonCriticalErrors, err
		} else {
			glog.Errorf("Non-critical error occured during resource retrieval: %s", err)
			nonCriticalErrors = appendMissing(nonCriticalErrors, err)
		}
	}
	return nonCriticalErrors, nil
}

// MergeErrors merges multiple non-critical error slice into one array
func MergeErrors(errorArrayToMerge ...[]error) (mergedErrors []error) {
	for _, errorArray := range errorArrayToMerge {
		mergedErrors = appendMissing(mergedErrors, errorArray...)
	}
	return
}

func isErrorCritical(err error) bool {
	status, ok := err.(*errors.StatusError)
	if !ok {
		// Assume, that error is critical if it cannot be mapped.
		return true
	}

	return !contains(NonCriticalErrors, status.ErrStatus.Code)
}

func appendMissing(slice []error, toAppend ...error) []error {
	m := make(map[string]bool, 0)

	for _, s := range slice {
		m[s.Error()] = true
	}

	for _, e := range toAppend {
		_, ok := m[e.Error()]
		if !ok {
			slice = append(slice, e)
			m[e.Error()] = true
		}
	}

	return slice
}

func contains(s []int32, e int32) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// HandleInternalError writes given error to the response
func HandleInternalError(response *restful.Response, err error) {
	glog.Error(err)
	statusCode := http.StatusInternalServerError
	statusError, ok := err.(*errors.StatusError)
	if ok && statusError.Status().Code > 0 {
		statusCode = int(statusError.Status().Code)
	}
	response.AddHeader("Content-Type", "text/plain")
	response.WriteErrorString(statusCode, err.Error()+"\n")
}

// Return a http status according to the error
func HandleHTTPError(err error) int {
	if err == nil {
		return http.StatusInternalServerError
	}

	if err.Error() == MSG_TOKEN_EXPIRED_ERROR {
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}
