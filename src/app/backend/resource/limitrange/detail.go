package limitrange

import "k8s.io/api/core/v1"

type limitRangesMap map[v1.LimitType]rangeMap

type rangeMap map[v1.ResourceName]*LimitRangeItem

func (rMap rangeMap) getRange(resource v1.ResourceName) *LimitRangeItem {
	r, ok := rMap[resource]
	if !ok {
		rMap[resource] = &LimitRangeItem{}
		return rMap[resource]
	} else {
		return r
	}
}

type LimitRangeItem struct {
	// ResourceName usage constraints on this kind by resource name
	ResourceName string `json:"resourceName,omitempty"`

	// ResourceType of resource that this limit applies to
	ResourceType string `json:"resourceType,omitempty"`

	// Min usage constraints on this kind by resource name
	Min string `json:"min,omitempty"`

	// Max usage constraints on this kind by resource name
	Max string `json:"max,omitempty"`

	// Default resource requirement limit value by resource name.
	Default string `json:"default,omitempty"`

	// DefaultRequest resource requirement request value by resource name.
	DefaultRequest string `json:"defaultRequest,omitempty"`

	// MaxLimitRequestRatio represents the max burst value for the named resource
	MaxLimitRequestRatio string `json:"maxLimitRequestRatio,omitempty"`
}

func toLimitRangesMap(rawLimitRange *v1.LimitRange) limitRangesMap {
	rawLimitRanges := rawLimitRange.Spec.Limits

	limitRanges := make(limitRangesMap, len(rawLimitRanges))

	for _, rawLimitRange := range rawLimitRanges {
		rangeMap := make(rangeMap)

		for resource, min := range rawLimitRange.Min {
			rangeMap.getRange(resource).Min = min.String()
		}

		for resource, max := range rawLimitRange.Max {
			rangeMap.getRange(resource).Max = max.String()
		}

		for resource, df := range rawLimitRange.Default {
			rangeMap.getRange(resource).Default = df.String()
		}

		for resource, dfR := range rawLimitRange.DefaultRequest {
			rangeMap.getRange(resource).DefaultRequest = dfR.String()
		}

		for resource, mLR := range rawLimitRange.MaxLimitRequestRatio {
			rangeMap.getRange(resource).MaxLimitRequestRatio = mLR.String()
		}

		limitRanges[rawLimitRange.Type] = rangeMap
	}

	return limitRanges
}

func ToLimitRanges(rawLimitRange *v1.LimitRange) []LimitRangeItem {
	limitRangeMap := toLimitRangesMap(rawLimitRange)
	limitRangeList := make([]LimitRangeItem, 0)

	for limitType, rangeMap := range limitRangeMap {
		for resource, limit := range rangeMap {
			limit.ResourceName = resource.String()
			limit.ResourceType = string(limitType)
			limitRangeList = append(limitRangeList, *limit)
		}
	}

	return limitRangeList
}