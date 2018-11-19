package node

import (
	"github.com/golang/glog"
	"github.com/wzt3309/k8sconsole/src/app/backend/api"
	kcErrors "github.com/wzt3309/k8sconsole/src/app/backend/errors"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	"k8s.io/api/core/v1"
	client "k8s.io/client-go/kubernetes"
)

// NodeList contains a list of nodes in the cluster.
type NodeList struct {
	ListMeta 	api.ListMeta 	`json:"listMeta"`
	Nodes 		[]Node 				`json:"nodes"`
	// List of non-critical errors, that occurred during resource retrieval.
	Errors 		[]error 			`json:"errors"`
}

// Node is a presentation layer view of kubernetes nodes.
type Node struct {
	ObjectMeta 					api.ObjectMeta `json:"objectMeta"`
	TypeMeta 						api.TypeMeta	 `json:"typeMeta"`
	Ready 							v1.ConditionStatus `json:"ready"`
	AllocatedResources	NodeAllocatedResources `json:"allocatedResources"`
}

// GetNodeListFromChannels returns a list of all Nodes in the cluster.
func GetNodeListFromChannels(client client.Interface, channels *common.ResourceChannels,
	dsQuery *dataselect.DataSelectQuery) (*NodeList, error) {

	nodes := <- channels.NodeList.List
	err := <- channels.NodeList.Error

	nonCriticalErrors, criticalError := kcErrors.HandleError(err)
	if criticalError != nil {
		return nil, criticalError
	}

	return toNodeList(client, nodes.Items, nonCriticalErrors, dsQuery), nil
}

// GetNodeList returns a list of all Nodes in the cluster.
func GetNodeList(client client.Interface, dsQuery *dataselect.DataSelectQuery) (*NodeList, error) {
	nodes, err := client.CoreV1().Nodes().List(api.ListEverything)

	nonCriticalErrors, criticalErrors := kcErrors.HandleError(err)
	if criticalErrors != nil {
		return nil, criticalErrors
	}

	return toNodeList(client, nodes.Items, nonCriticalErrors, dsQuery), nil
}

// GetNodeList returns a list of all Nodes in the cluster.
func toNodeList(client client.Interface, nodes []v1.Node,
	nonCriticalErrors []error, dsQuery *dataselect.DataSelectQuery) *NodeList {

	nodeList := &NodeList{
		Nodes: make([]Node, 0),
		ListMeta: api.ListMeta{TotalItems: len(nodes)},
		Errors: nonCriticalErrors,
	}

	nodeCells, filteredTotal := dataselect.GenericDataSelectWithFilter(toCells(nodes), dsQuery)
	nodes = fromCells(nodeCells)
	nodeList.ListMeta = api.ListMeta{TotalItems: filteredTotal}

	for _, node := range nodes {
		pods, err := getNodePods(client, node)
		if err != nil {
			glog.Errorf("Couldn't get pods of %s node: %s\n", node.Name, err)
		}

		nodeList.Nodes = append(nodeList.Nodes, toNode(node, pods))
	}

	return nodeList
}

func toNode(node v1.Node, pods *v1.PodList) Node {
	allocatedResources, err := getNodeAllocatedResources(node, pods)
	if err != nil {
		glog.Errorf("Couldn't get allocated resources of %s node: %s\n", node.Name, err)
	}

	return Node{
		ObjectMeta:         api.NewObjectMeta(node.ObjectMeta),
		TypeMeta:           api.NewTypeMeta(api.ResourceKindNode),
		Ready:              getNodeConditionStatus(node, v1.NodeReady),
		AllocatedResources: allocatedResources,
	}
}

func getNodeConditionStatus(node v1.Node, conditionType v1.NodeConditionType) v1.ConditionStatus {
	for _, condition := range node.Status.Conditions {
		if condition.Type == conditionType {
			return condition.Status
		}
	}
	return v1.ConditionUnknown
}


