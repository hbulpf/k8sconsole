package deployment

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	apps "k8s.io/api/apps/v1beta2"
	"k8s.io/api/core/v1"
)

func FindOldReplicaSets(deployment *apps.Deployment, rsList []*apps.ReplicaSet) (
	[]*apps.ReplicaSet, []*apps.ReplicaSet, error) {
	var requiredRSs []*apps.ReplicaSet
	var allRSs []*apps.ReplicaSet
	newRS, err := FindNewReplicaSet(deployment, rsList)
	if err != nil {
		return nil, nil, err
	}
	for _, rs := range rsList {
		if newRS != nil && rs.UID == newRS.UID {
			continue
		}
		allRSs = append(allRSs, rs)
		if *(rs.Spec.Replicas) != 0 {
			requiredRSs = append(requiredRSs, rs)
		}
	}
	return requiredRSs, allRSs, nil
}

func FindNewReplicaSet(deployment *apps.Deployment, rsList []*apps.ReplicaSet) (
	*apps.ReplicaSet, error) {
	newRSTemplate := GetNewReplicaSetTemplate(deployment)
	for i := range rsList {
		if common.EqualIgnoreHash(rsList[i].Spec.Template, newRSTemplate) {
			return rsList[i], nil
		}
	}
	// new ReplicaSet does not exist
	return nil, nil
}

// GetNewReplicaSetTemplate returns the desired PodTemplateSpec for the new ReplicaSet corresponding to the given ReplicaSet.
// Callers of this helper need to set the DefaultDeploymentUniqueLabelKey k/v pair.
func GetNewReplicaSetTemplate(deployment *apps.Deployment) v1.PodTemplateSpec {
	// newRS will have the same template as in deployment spec.
	return v1.PodTemplateSpec{
		ObjectMeta: deployment.Spec.Template.ObjectMeta,
		Spec:       deployment.Spec.Template.Spec,
	}
}