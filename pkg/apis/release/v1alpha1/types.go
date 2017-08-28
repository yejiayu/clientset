/*
Copyright 2017 caicloud authors. All rights reserved.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

// ReleaseRollbackConfig describes the rollback config of a release
type ReleaseRollbackConfig struct {
	// The version to rollback to. If set to 0, rollbck to the last version.
	Version int32 `json:"version,omitempty"`
}

// ReleaseSpec describes the basic info of a release
type ReleaseSpec struct {
	// Description is the description of current release
	Description string `json:"description,omitempty"`
	// Template is an archived template data
	Template []byte `json:"template,omitempty"`
	// Config is the config for parsing template
	Config string `json:"config,omitempty"`
	// The config this release is rolling back to. Will be cleared after rollback is done.
	RollbackTo *ReleaseRollbackConfig `json:"rollbackTo,omitempty"`
}

type ReleaseConditionType string

const (
	// ReleaseAvailable means the resources of release are available and can render service.
	ReleaseAvailable ReleaseConditionType = "Available"
	// ReleaseProgressing means release is playing a mutation. It occurs when create/update/rollback
	// release. If some bad thing was trigger, release transfers to ReleaseFailure.
	ReleaseProgressing ReleaseConditionType = "Progressing"
	// ReleaseFailure means some parts of release falled into wrong field. Some parts may work
	// as usual, but the release can't provide complete service.
	ReleaseFailure ReleaseConditionType = "Failure"
)

// ReleaseHistorySpec describes the history info of a release
type ReleaseCondition struct {
	// Type of release condition.
	Type ReleaseConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status apiv1.ConditionStatus `json:"status"`
	// Last time the condition transit from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// Human readable message indicating details about last transition.
	Message string `json:"message,omitempty"`
}

// ResourceCounter is a status counter
type ResourceCounter struct {
	// Available is the count of running target
	Available int32 `json:"available"`
	// Progressing is the count of mutating target
	Progressing int32 `json:"progressing"`
	// Failure is the count of wrong target
	Failure int32 `json:"failure"`
}

// ReleaseDetailStatus
type ReleaseDetailStatus struct {
	// Path is the path which resources from
	Path string `json:"path,omitempty"`
	// Resources contains a kind-counter map.
	// A kind should be a unique name of a group resources.
	Resources map[string]ResourceCounter `json:"resources,omitempty"`
}

// ReleaseStatus describes the status of a release
type ReleaseStatus struct {
	// Version is the version of current release
	Version int32 `json:"version,omitempty"`
	// Manifest is the generated kubernetes resources from template
	Manifest string `json:"manifest,omitempty"`
	// LastUpdateTime is the last update time of current release
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// Details contains all resources status of current release
	Details []ReleaseDetailStatus `json:"details,omitempty"`
	// Conditions is an array of current observed release conditions.
	Conditions []ReleaseCondition `json:"conditions,omitempty"`
}

// +genclient=true
// +genclientstatus=false

// Release describes a release wich chart and values
type Release struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the Release
	// +optional
	Spec ReleaseSpec `json:"spec,omitempty"`

	// Most recently observed status of the Release
	// +optional
	Status ReleaseStatus `json:"status,omitempty"`
}

// ReleaseList describes an array of Release instances
type ReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of releases
	Items []Release `json:"items"`
}

// ReleaseHistorySpec describes the history info of a release
type ReleaseHistorySpec struct {
	// Description is the description of current history
	Description string `json:"description,omitempty"`
	// Version is the version of a history
	Version int32 `json:"version,omitempty"`
	// Template is an archived template data
	Template []byte `json:"template,omitempty"`
	// Config is the config for parsing template
	Config string `json:"config,omitempty"`
	// Manifest is the generated kubernetes resources from template
	Manifest string `json:"manifest,omitempty"`
}

// +genclient=true
// +genclientstatus=false

// ReleaseHistory describes a history of a release version
type ReleaseHistory struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired behavior of the ReleaseHistory
	// +optional
	Spec ReleaseHistorySpec `json:"spec,omitempty"`
}

// ReleaseHistoryList describes an array of ReleaseHistory instances
type ReleaseHistoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of release histories
	Items []ReleaseHistory `json:"items"`
}
