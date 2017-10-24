/*
Copyright 2017 caicloud authors. All rights reserved.
*/

// This file was automatically generated by informer-gen

package v1alpha2

import (
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// LoadBalancers returns a LoadBalancerInformer.
	LoadBalancers() LoadBalancerInformer
}

type version struct {
	internalinterfaces.SharedInformerFactory
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory) Interface {
	return &version{f}
}

// LoadBalancers returns a LoadBalancerInformer.
func (v *version) LoadBalancers() LoadBalancerInformer {
	return &loadBalancerInformer{factory: v.SharedInformerFactory}
}
