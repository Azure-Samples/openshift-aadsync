package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AADGroupSync spec
type AADGroupSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AADGroupSyncSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AADGroupSyncList spec
type AADGroupSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AADGroupSync `json:"items"`
}

// AADGroupSyncSpec spec
type AADGroupSyncSpec struct {
	Group            Group  `json:"group"`
	LastSyncDateTime string `json:"lastSyncDateTime"`
	LastSyncType     string `json:"lastSyncType"`
}

// Group spec
type Group struct {
	ObjectID    string `json:"objectId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserCount   int    `json:"userCount"`
	Users       []User `json:"users"`
}

// User spec
type User struct {
	ObjectID          string `json:"objectId"`
	UserPrincipalName string `json:"userPrincipalName"`
}
