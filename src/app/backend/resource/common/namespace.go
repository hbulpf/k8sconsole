package common

import api "k8s.io/api/core/v1"

// NamespaceQuery is a query for namespace of a list of object.
// cases 3:
//	1. No namespace selected. Means querying all namespace except 'kube-system'
//  2. single namespace selected.
//	3. More than one namespaces selected.
type NamespaceQuery struct {
	namespaces []string
}

// NewOneNamespaceQuery creates new namespace query that queries one namespace.
func NewOneNamespaceQuery(namespace string) *NamespaceQuery {
	return &NamespaceQuery{[]string{namespace}}
}

// NewNamespaceQuery creates new query for given namespaces.
func NewNamespaceQuery(namespaces []string) *NamespaceQuery {
	return &NamespaceQuery{namespaces}
}

// ToRequestParam returns K8S api namespace query for list of objects.
// If nsQuery's namespace is single, just return the namespace, otherwise
// return all namespaces and then use `Matches` to filter
func (n *NamespaceQuery) ToRequestParam() string {
	if len(n.namespaces) == 1 {
		return n.namespaces[0]
	}

	return api.NamespaceAll
}

// Matches returns true when the given namespace matches this query.
func (n *NamespaceQuery) Matches(namespace string) bool {
	if len(n.namespaces) == 0 {
		return true
	}

	for _, queryNamespace := range n.namespaces {
		if namespace == queryNamespace {
			return true
		}
	}

	return false
}

