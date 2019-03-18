package mapper

import (
	"time"

	msgraph "pkg/msgraph/types"

	logrus "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	aadgroupsyncv1 "pkg/aadsync/apis/aad.microsoft.com/v1"
)

// Mapper contains the internal mapper details
type Mapper struct {
	Log *logrus.Entry
}

// NewClient creates a new Mapper Client
func NewClient(log *logrus.Entry) *Mapper {

	mapper := &Mapper{
		Log: log,
	}
	log.Info("Created mapper")

	return mapper
}

// CreateFromMSGraphGroup creates a new aadgroupsyncs.aad.microsoft.com CRD entry from an AAD Group
func (m *Mapper) CreateFromMSGraphGroup(msGraphGroup *msgraph.Group) *aadgroupsyncv1.AADGroupSync {

	aadsyncGroup := &aadgroupsyncv1.AADGroupSync{
		ObjectMeta: metav1.ObjectMeta{
			Name: msGraphGroup.ID,
		},
		Spec: aadgroupsyncv1.AADGroupSyncSpec{
			Group: aadgroupsyncv1.Group{
				ObjectID:    msGraphGroup.ID,
				Name:        msGraphGroup.DisplayName,
				Description: msGraphGroup.Description,
				UserCount:   len(msGraphGroup.Users),
				Users:       []aadgroupsyncv1.User{},
			},
			LastSyncDateTime: time.Now().UTC().Format(time.RFC3339),
			LastSyncType:     "Scheduled",
		},
	}

	for _, user := range msGraphGroup.Users {

		aadsyncGroupUser := aadgroupsyncv1.User{
			ObjectID:          user.ID,
			UserPrincipalName: user.UserPrincipalName,
		}
		aadsyncGroup.Spec.Group.Users = append(aadsyncGroup.Spec.Group.Users, aadsyncGroupUser)
	}

	return aadsyncGroup
}

// UpdateFromMSGraphGroup updates an existing aadgroupsyncs.aad.microsoft.com CRD entry from an AAD Group
func (m *Mapper) UpdateFromMSGraphGroup(msGraphGroup *msgraph.Group, aadsyncGroup *aadgroupsyncv1.AADGroupSync) *aadgroupsyncv1.AADGroupSync {

	aadsyncGroup.Spec.Group.ObjectID = msGraphGroup.ID
	aadsyncGroup.Spec.Group.Name = msGraphGroup.DisplayName
	aadsyncGroup.Spec.Group.Description = msGraphGroup.Description
	aadsyncGroup.Spec.Group.UserCount = len(msGraphGroup.Users)
	aadsyncGroup.Spec.Group.Users = []aadgroupsyncv1.User{}
	aadsyncGroup.Spec.LastSyncDateTime = time.Now().UTC().Format(time.RFC3339)
	aadsyncGroup.Spec.LastSyncType = "Scheduled"

	for _, user := range msGraphGroup.Users {

		aadsyncGroupUser := aadgroupsyncv1.User{
			ObjectID:          user.ID,
			UserPrincipalName: user.UserPrincipalName,
		}
		aadsyncGroup.Spec.Group.Users = append(aadsyncGroup.Spec.Group.Users, aadsyncGroupUser)
	}

	return aadsyncGroup
}
