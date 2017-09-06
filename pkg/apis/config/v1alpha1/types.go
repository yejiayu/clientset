/*
Copyright 2017 caicloud authors. All rights reserved.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigClaimStatusType is sync status of config claim
type ConfigClaimStatusType string

const (
	// Unknown means that config is sync not yet
	Unknown ConfigClaimStatusType = "Unknown"
	// Success means taht config is sync success
	Success ConfigClaimStatusType = "Success"
	// Failure  means taht config is sync failuer
	Failure ConfigClaimStatusType = "Failure"
)

// +genclient=true
// +genclientstatus=true

// ConfigClaim describes a config sync status
type ConfigClaim struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Most recently observed status of the Release
	// +optional
	Status ConfigClaimStatus `json:"status,omitempty"`
}

// ConfigClaimStatus describes the status of a ConfigClaim
type ConfigClaimStatus struct {
	// Status is sync status of Config
	Status ConfigClaimStatusType `json:"status"`
	// Reason describes success or Failure of status
	Reason string `json:"reason,omitempty"`
}

// ConfigClaimList describes an array of ConfigClaim instances
type ConfigClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of ConfigClaim
	Items []ConfigClaim `json:"items"`
}

// +genclient=true
// +genclientstatus=false

// ConfigReference describes the config reference list.
type ConfigReference struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// Specification of the desired behavior of the Release
	// +optional
	Spec ConfigReferenceSpec `json:"spec,omitempty"`
}

// ConfigReferenceSpec describes the config reference list.
type ConfigReferenceSpec struct {
	Refs []*Reference `json:"refs,omitempty"`
}

// Reference describes the config reference.
type Reference struct {
	Name       string       `json:"name"`
	Kind       string       `json:"kind"`
	APIVersion string       `json:"apiVersion"`
	Config     []Data       `json:"config,omitempty"`
	Children   []*Reference `json:"children,omitempty"`
}

// Data describes the config info.
type Data struct {
	Name string   `json:"name"`
	Keys []string `json:"keys"`
}

// ConfigReferenceList describes an array of ConfigReference instances.
type ConfigReferenceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of ConfigClaim
	Items []ConfigReference `json:"items"`
}
