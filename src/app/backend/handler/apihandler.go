package handler

import (
	"github.com/emicklei/go-restful"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	"github.com/wzt3309/k8sconsole/src/app/backend/auth"
	authApi "github.com/wzt3309/k8sconsole/src/app/backend/auth/api"
	clientApi "github.com/wzt3309/k8sconsole/src/app/backend/client/api"
	kcErrors "github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/configmap"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/event"
	ns "github.com/wzt3309/k8sconsole/src/app/backend/resource/namespace"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/node"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/persistentvolume"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/persistentvolumeclaim"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/pod"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/secret"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/service"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/storageclass"
	"golang.org/x/net/xsrftoken"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type APIHandler struct {
	cManager clientApi.ClientManager
}

// CreateHTTPAPIHandler creates a new HTTP handler that handles all requests to the API of the backend.
func CreateHTTPAPIHandler(cManager clientApi.ClientManager, authManager authApi.AuthManager) (http.Handler, error) {
	apiHandler := APIHandler{cManager: cManager}

	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	apiV1Ws := new(restful.WebService)
	apiV1Ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	wsContainer.Add(apiV1Ws)

	authHandler := auth.NewAuthHandler(authManager)
	authHandler.Install(apiV1Ws)

	apiV1Ws.Route(
		apiV1Ws.GET("csrftoken/{action}").
			To(apiHandler.handleGetCsrfToken).
			Writes(api.CsrfToken{}))

	apiV1Ws.Route(
		apiV1Ws.GET("/pod").
			To(apiHandler.handleGetPods).
			Writes(pod.PodList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/pod/{namespace}").
			To(apiHandler.handleGetPods).
			Writes(pod.PodList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/pod/{namespace}/{pod}").
			To(apiHandler.handleGetPodDetail).
			Writes(pod.PodDetail{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/pod/{namespace}/{pod}/event").
			To(apiHandler.handleGetPodEvents).
			Writes(common.EventList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/pod/{namespace}/{pod}/persistentvolumeclaim").
			To(apiHandler.handleGetPodPersistentVolumeClaims).
			Writes(persistentvolumeclaim.PersistentVolumeClaimList{}))

	apiV1Ws.Route(
		apiV1Ws.POST("/namespace").
			To(apiHandler.handleCreateNamespace).
			Reads(ns.NamespaceSpec{}).
			Writes(ns.NamespaceSpec{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/namespace").
			To(apiHandler.handleGetNamespaces).
			Writes(ns.NamespaceList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/namespace/{name}").
			To(apiHandler.handleGetNamespaceDetail).
			Writes(ns.NamespaceDetail{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/namespace/{name}/event").
			To(apiHandler.handleGetNamespaceEvents).
			Writes(common.EventList{}))

	apiV1Ws.Route(
		apiV1Ws.GET("/secret").
			To(apiHandler.handleGetSecretList).
			Writes(secret.SecretList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/secret/{namespace}").
			To(apiHandler.handleGetSecretList).
			Writes(secret.SecretList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/secret/{namespace}/{name}").
			To(apiHandler.handleGetSecretDetail).
			Writes(secret.SecretDetail{}))
	apiV1Ws.Route(
		apiV1Ws.POST("/secret").
			To(apiHandler.handleCreateImagePullSecret).
			Reads(secret.ImagePullSecretSpec{}).
			Writes(secret.Secret{}))

	apiV1Ws.Route(
		apiV1Ws.GET("/configmap").
			To(apiHandler.handleGetConfigMapList).
			Writes(configmap.ConfigMapList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/configmap/{namespace}").
			To(apiHandler.handleGetConfigMapList).
			Writes(configmap.ConfigMapList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/configmap/{namespace}/{configmap}").
			To(apiHandler.handleGetConfigMapDetail).
			Writes(configmap.ConfigMapDetail{}))

	apiV1Ws.Route(
		apiV1Ws.GET("/service").
			To(apiHandler.handleGetServiceList).
			Writes(service.ServiceList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/service/{namespace}").
			To(apiHandler.handleGetServiceList).
			Writes(service.ServiceList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/service/{namespace}/{service}").
			To(apiHandler.handleGetServiceDetail).
			Writes(service.ServiceDetail{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/service/{namespace}/{service}/pod").
			To(apiHandler.handleGetServicePods).
			Writes(pod.PodList{}))

	apiV1Ws.Route(
		apiV1Ws.GET("/node").
			To(apiHandler.handleGetNodeList).
			Writes(node.NodeList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/node/{name}").
			To(apiHandler.handleGetNodeDetail).
			Writes(node.NodeDetail{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/node/{name}/event").
			To(apiHandler.handleGetNodeEvents).
			Writes(common.EventList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/node/{name}/pod").
			To(apiHandler.handleGetNodePods).
			Writes(pod.PodList{}))

	apiV1Ws.Route(
		apiV1Ws.GET("/persistentvolume").
			To(apiHandler.handleGetPersistentVolumeList).
			Writes(persistentvolume.PersistentVolumeList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/persistentvolume/{persistentvolume}").
			To(apiHandler.handleGetPersistentVolumeDetail).
			Writes(persistentvolume.PersistentVolumeDetail{}))

	apiV1Ws.Route(
		apiV1Ws.GET("/persistentvolumeclaim").
			To(apiHandler.handleGetPersistentVolumeClaimList).
			Writes(persistentvolumeclaim.PersistentVolumeClaimList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/persistentvolumeclaim/{namespace}").
			To(apiHandler.handleGetPersistentVolumeClaimList).
			Writes(persistentvolumeclaim.PersistentVolumeClaimList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/persistentvolumeclaim/{namespace}/{name}").
			To(apiHandler.handleGetPersistentVolumeClaimDetail).
			Writes(persistentvolumeclaim.PersistentVolumeClaimDetail{}))

	apiV1Ws.Route(
		apiV1Ws.GET("/storageclass").
			To(apiHandler.handleGetStorageClassList).
			Writes(storageclass.StorageClassList{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/storageclass/{storageclass}").
			To(apiHandler.handleGetStorageClass).
			Writes(storageclass.StorageClass{}))
	apiV1Ws.Route(
		apiV1Ws.GET("/storageclass/{storageclass}/persistentvolume").
			To(apiHandler.handleGetStorageClassPersistentVolumes).
			Writes(persistentvolume.PersistentVolumeList{}))
	return wsContainer, nil
}

func (apiHandler *APIHandler) handleGetCsrfToken(request *restful.Request, response *restful.Response) {
	action := request.PathParameter("action")
	token := xsrftoken.Generate(apiHandler.cManager.CSRFKey(), "none", action)
	response.WriteHeaderAndEntity(http.StatusOK, api.CsrfToken{Token: token})
}

func (apiHandler *APIHandler) handleGetPods(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespace := parseNamespacePathParameter(request)
	dsQuery := parseDataSelectPathParameter(request)
	result, err := pod.GetPodList(k8sClient, namespace, dsQuery)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPodDetail(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("pod")
	result, err := pod.GetPodDetail(k8sClient, namespace, name)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPodEvents(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	log.Println("Getting events related to a pod in namespace")
	namespace := request.PathParameter("namespace")
	name := request.PathParameter("pod")
	dsQuery := parseDataSelectPathParameter(request)
	result, err := pod.GetEventsForPod(k8sClient, dsQuery, namespace, name)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPodPersistentVolumeClaims(request *restful.Request,
	response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("pod")
	namespace := request.PathParameter("namespace")
	dataSelect := parseDataSelectPathParameter(request)
	result, err := persistentvolumeclaim.GetPodPersistentVolumeClaims(k8sClient,
		namespace, name, dataSelect)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleCreateNamespace(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespaceSpec := new(ns.NamespaceSpec)
	if err := request.ReadEntity(namespaceSpec); err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	if err := ns.CreateNamespace(namespaceSpec, k8sClient); err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, namespaceSpec)
}

func (apiHandler *APIHandler) handleGetNamespaces(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	dsQuery := parseDataSelectPathParameter(request)
	result, err := ns.GetNamespaceList(k8sClient, dsQuery)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetNamespaceDetail(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("name")
	result, err := ns.GetNamespaceDetail(k8sClient, name)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetNodeList(request *restful.Request, response *restful.Response) {
	k8sclient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	dsQuery := parseDataSelectPathParameter(request)
	result, err := node.GetNodeList(k8sclient, dsQuery)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetNodeDetail(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("name")
	dsQuery := parseDataSelectPathParameter(request)
	result, err := node.GetNodeDetail(k8sClient, name, dsQuery)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetNamespaceEvents(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("name")
	dsQuery := parseDataSelectPathParameter(request)
	result, err := event.GetNamespaceEvents(k8sClient, dsQuery, name)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetSecretList(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespace := parseNamespacePathParameter(request)
	dsQuery := parseDataSelectPathParameter(request)
	result, err := secret.GetSecretList(k8sClient, namespace, dsQuery)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetSecretDetail(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	result, err := secret.GetSecretDetail(k8sClient, namespace, name)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleCreateImagePullSecret(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	spec := new(secret.ImagePullSecretSpec)
	if err := request.ReadEntity(spec); err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	result, err := secret.CreateSecret(k8sClient, spec)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetConfigMapList(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespace := parseNamespacePathParameter(request)
	dataSelect := parseDataSelectPathParameter(request)
	result, err := configmap.GetConfigMapList(k8sClient, namespace, dataSelect)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetConfigMapDetail(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("configmap")
	result, err := configmap.GetConfigMapDetail(k8sClient, namespace, name)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetServiceList(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespace := parseNamespacePathParameter(request)
	dsQuery := parseDataSelectPathParameter(request)
	result, err := service.GetServiceList(k8sClient, namespace, dsQuery)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetServiceDetail(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("service")
	dsQuery := parseDataSelectPathParameter(request)
	result, err := service.GetServiceDetail(k8sClient, namespace, name, dsQuery)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetServicePods(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("service")
	dsQuery := parseDataSelectPathParameter(request)
	result, err := service.GetServicePods(k8sClient, namespace, name, dsQuery)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetNodeEvents(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("name")
	dsQuery := parseDataSelectPathParameter(request)
	result, err := event.GetNodeEvents(k8sClient, dsQuery, name)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetNodePods(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("name")
	dsQuery := parseDataSelectPathParameter(request)
	result, err := node.GetNodePods(k8sClient, dsQuery, name)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPersistentVolumeList(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	dataSelect := parseDataSelectPathParameter(request)
	result, err := persistentvolume.GetPersistentVolumeList(k8sClient, dataSelect)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPersistentVolumeDetail(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("persistentvolume")
	result, err := persistentvolume.GetPersistentVolumeDetail(k8sClient, name)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPersistentVolumeClaimList(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespace := parseNamespacePathParameter(request)
	dataSelect := parseDataSelectPathParameter(request)
	result, err := persistentvolumeclaim.GetPersistentVolumeClaimList(k8sClient, namespace, dataSelect)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetPersistentVolumeClaimDetail(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	namespace := request.PathParameter("namespace")
	name := request.PathParameter("name")
	result, err := persistentvolumeclaim.GetPersistentVolumeClaimDetail(k8sClient, namespace, name)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetStorageClassList(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	dsQuery := parseDataSelectPathParameter(request)
	result, err := storageclass.GetStorageClassList(k8sClient, dsQuery)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetStorageClass(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("storageclass")
	result, err := storageclass.GetStorageClassDetail(k8sClient, name)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (apiHandler *APIHandler) handleGetStorageClassPersistentVolumes(request *restful.Request, response *restful.Response) {
	k8sClient, err := apiHandler.cManager.Client(request)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}

	name := request.PathParameter("storageclass")
	dsQuery := parseDataSelectPathParameter(request)
	result, err := persistentvolume.GetStorageClassPersistentVolumes(k8sClient, name, dsQuery)
	if err != nil {
		kcErrors.HandleInternalError(response, err)
		return
	}
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// Get namespaces from path parameter
func parseNamespacePathParameter(request *restful.Request) *common.NamespaceQuery {
	namespace := request.PathParameter("namespace")
	namespaces := strings.Split(namespace, ",")
	var noBlankNamespaces []string
	for _, n := range namespaces {
		n = strings.Trim(n, " ")
		if len(n) > 0 {
			noBlankNamespaces = append(noBlankNamespaces, n)
		}
	}

	return common.NewNamespaceQuery(noBlankNamespaces)
}

func parsePaginationPathParameter(request *restful.Request) *dataselect.PaginationQuery {
	itemsPerPage, err := strconv.ParseInt(request.QueryParameter("itemsPerPage"), 10, 0)
	if err != nil {
		return dataselect.NoPagination
	}

	page, err := strconv.ParseInt(request.QueryParameter("page"), 10, 0)
	if err != nil {
		return dataselect.NoPagination
	}

	// Frontend page start from 1 and backend start from 0
	return dataselect.NewPaginationQuery(int(itemsPerPage), int(page - 1))
}

func parseFilterPathParameter(request *restful.Request) *dataselect.FilterQuery {
	return dataselect.NewFilterQuery(strings.Split(request.QueryParameter("filterBy"), ","))
}

func parseSortPathParameter(request *restful.Request) *dataselect.SortQuery {
	return dataselect.NewSortQuery(strings.Split(request.QueryParameter("sortBy"), ","))
}

// Parses query parameters of the request and returns a DataSelectQuery object
func parseDataSelectPathParameter(request *restful.Request) *dataselect.DataSelectQuery {
	paginationQuery := parsePaginationPathParameter(request)
	sortQuery := parseSortPathParameter(request)
	filterQuery := parseFilterPathParameter(request)
	return dataselect.NewDataSelectQuery(paginationQuery, sortQuery, filterQuery)
}