package handler

import (
	"github.com/emicklei/go-restful"
	kcErrors "github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"io"
)

func handleDownload(response *restful.Response, result io.ReadCloser) {
	response.AddHeader(restful.HEADER_ContentType, "text/plain")
	defer result.Close()
	_, err := io.Copy(response, result)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
}
