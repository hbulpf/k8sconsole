package discovery

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/ingress"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/service"
	"k8s.io/client-go/kubernetes"
)

type Discovery struct {
	ServiceList service.ServiceList `json:"serviceList"`
	IngressList ingress.IngressList `json:"ingressList"`

	// List of non-critical errors, that occurred during resource retrieval
	Errors []error `json:"errors"`
}

func GetDiscovery(client kubernetes.Interface, nsQuery *common.NamespaceQuery,
	dsQuery *dataselect.DataSelectQuery) (*Discovery, error) {

		glog.Info("Getting discovery and load balancing category")
		channels := &common.ResourceChannels{
			ServiceList: common.GetServiceListChannel(client, nsQuery, 1),
			IngressList: common.GetIngressListChannel(client, nsQuery, 1),
		}

		return GetDiscoveryFromChannels(channels, dsQuery)
}

func GetDiscoveryFromChannels(channels *common.ResourceChannels,
	dsQuery *dataselect.DataSelectQuery) (*Discovery, error) {

	numErrs := 2
	errChan := make(chan error, numErrs)
	svcChan := make(chan *service.ServiceList)
	ingressChan := make(chan *ingress.IngressList)

	go func() {
		items, err := service.GetServiceListFromChannels(channels, dsQuery)
		errChan <- err
		svcChan <- items
	}()

	go func() {
		items, err := ingress.GetIngressListFromChannels(channels, dsQuery)
		errChan <- err
		ingressChan <- items
	}()

	for i := 0; i < numErrs; i++ {
		err := <- errChan
		if err != nil {
			return nil, err
		}
	}

	discovery := &Discovery{
		ServiceList: *(<-svcChan),
		IngressList: *(<-ingressChan),
	}

	discovery.Errors = errors.MergeErrors(discovery.ServiceList.Errors, discovery.IngressList.Errors)
	return discovery, nil
}
