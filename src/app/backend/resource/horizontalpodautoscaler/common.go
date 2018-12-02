package horizontalpodautoscaler

// Simple mapping of an autoscaling.CrossVersionObjectReference
type ScaleTargetRef struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

