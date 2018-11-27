package container

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/logs"
	"io"
	"io/ioutil"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

// maximum number of lines loaded from the apiserver
var lineReadLimit int64 = 5000

// maximum number of bytes loaded from the apiserver
var byteReadLimit int64 = 500000

func GetLogDetails(client kubernetes.Interface, namespace, podID string, container string,
	logSelector *logs.Selector, usePreviousLogs bool) (*logs.LogDetails, error) {
	pod, err := client.CoreV1().Pods(namespace).Get(podID, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if len(container) == 0 {
		container = pod.Spec.Containers[0].Name
	}

	logOptions := mapToLogOptions(container, logSelector, usePreviousLogs)
	rawLogs, err := readRawLogs(client, namespace, podID, logOptions)
	if err != nil {
		return nil, err
	}
	details := ConstructLogDetails(podID, rawLogs, container, logSelector)
	return details, nil
}

// Maps the log selection to the corresponding api object
func mapToLogOptions(container string, logSelector *logs.Selector, previous bool) *v1.PodLogOptions {
	logOptions := &v1.PodLogOptions{
		Container: container,
		Follow: false,
		Previous: previous,
		Timestamps: true,
	}

	if logSelector.LogFilePosition == logs.Beginning {
		logOptions.LimitBytes = &byteReadLimit
	} else {
		logOptions.TailLines = &lineReadLimit
	}

	return logOptions
}

// Construct a request for getting the logs for a pod and retrieves the logs.
func readRawLogs(client kubernetes.Interface, namespace, podID string, logOptions *v1.PodLogOptions) (
	string, error) {
	readCloser, err := openStream(client, namespace, podID, logOptions)
	if err != nil {
		return err.Error(), nil
	}

	defer readCloser.Close()

	result, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// GetLogFile returns a stream to the log file which can be piped directly to the response
func GetLogFile(client kubernetes.Interface, namespace, podID string, container string, usePreviousLogs bool) (
	io.ReadCloser, error) {
	logOptions := &v1.PodLogOptions{
		Container: container,
		Follow: false,
		Previous: usePreviousLogs,
		Timestamps: false,
	}
	logStream, err := openStream(client, namespace, podID, logOptions)
	return logStream, err
}

func openStream(client kubernetes.Interface, namespace, podID string, logOptions *v1.PodLogOptions) (
	io.ReadCloser, error) {
	return client.CoreV1().RESTClient().Get().
		Namespace(namespace).
		Name(podID).
		Resource("pods").
		SubResource("log").
		VersionedParams(logOptions, scheme.ParameterCodec).Stream()
}

func ConstructLogDetails(podID string, rawLogs string, container string, logSelector *logs.Selector) *logs.LogDetails {
	parsedLines := logs.ToLogLines(rawLogs)
	logLines, fromDate, toDate, logSelection, lastPage := parsedLines.SelectLogs(logSelector)

	readLimitReached := isReadLimitReached(int64(len(rawLogs)), int64(len(parsedLines)), logSelector.LogFilePosition)
	truncated := readLimitReached && lastPage

	info := logs.LogInfo{
		PodName: 				podID,
		ContainerName: 	container,
		FromDate: 			fromDate,
		ToDate: 				toDate,
		Truncated: 			truncated,
	}
	return &logs.LogDetails{
		Info: info,
		Selector: logSelection,
		LogLines: logLines,
	}
}

// Checks if the amount of log file returned from the apiserver is equal to the read limits
func isReadLimitReached(bytesLoaded int64, linesLoaded int64, logFilePosition string) bool {
	return (logFilePosition == logs.Beginning && bytesLoaded >= byteReadLimit) ||
		(logFilePosition == logs.End && linesLoaded >= lineReadLimit)
}