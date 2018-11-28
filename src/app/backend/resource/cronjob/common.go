package cronjob

import (
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/common"
	"github.com/wzt3309/k8sconsole/src/app/backend/resource/dataselect"
	batch2 "k8s.io/api/batch/v1beta1"
)

// The code below allows to perform complex data section on []batch.CronJob

type CronJobCell batch2.CronJob

func (self CronJobCell) GetProperty(name dataselect.PropertyName) dataselect.ComparableValue {
	switch name {
	case dataselect.NameProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Name)
	case dataselect.CreationTimestampProperty:
		return dataselect.StdComparableTime(self.ObjectMeta.CreationTimestamp.Time)
	case dataselect.NamespaceProperty:
		return dataselect.StdComparableString(self.ObjectMeta.Namespace)
	default:
		// if name is not supported then just return a constant dummy value, sort will have no effect.
		return nil
	}
}

func toCells(std []batch2.CronJob) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = CronJobCell(std[i])
	}
	return cells
}

func fromCells(cells []dataselect.DataCell) []batch2.CronJob {
	std := make([]batch2.CronJob, len(cells))
	for i := range std {
		std[i] = batch2.CronJob(cells[i].(CronJobCell))
	}
	return std
}

func getStatus(list *batch2.CronJobList) common.ResourceStatus {
	info := common.ResourceStatus{}
	if list == nil {
		return info
	}

	for _, cronJob := range list.Items {
		if cronJob.Spec.Suspend != nil && !(*cronJob.Spec.Suspend) {
			info.Running++
		} else {
			info.Failed++
		}
	}

	return info
}
