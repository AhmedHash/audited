// Package audited is used to log last UpdatedBy and CreatedBy for models
package audited

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

// AuditedModel make Model Auditable, embed `audited.Model` into a model as anonymous field to make the model auditable
//    type User struct {
//      gorm.Model
//      audited.Model
//    }
type Model struct{
	CreatedAt		time.Time		`json:"created_at,omitempty"`
	CreatedByID		uuid.UUID		`json:"created_by_id,omitempty"`
	CreatedByRole	int64			`json:"created_by_role,omitempty"`
	CreatedByName	string			`gorm:"-" json:"created_by_name,omitempty"`		// to be used in joins

	UpdatedAt		time.Time		`json:"updated_at,omitempty"`
	UpdatedByID		uuid.UUID		`json:"updated_by_id,omitempty"`
	UpdatedByRole	int64			`json:"updated_by_role,,omitempty"`
	UpdatedByName	string			`gorm:"-" json:"updated_by_name,omitempty"`		// to be used in joins

	DeletedAt		*time.Time		`json:"deleted_at,omitempty" sql:"index"`
}

type User struct {
	ID		uuid.UUID
	Role	int64
}

// SetCreatedBy set created by
func (m *Model) SetCreatedBy(user User) {
	m.CreatedByID = user.ID
	m.CreatedByRole = user.Role
}

// GetCreatedBy get created by
func (m Model) GetCreatedBy() User {
	return User{
		ID:   m.CreatedByID,
		Role: m.CreatedByRole,
	}
}

// SetUpdatedBy set updated by
func (m *Model) SetUpdatedBy(user User) {
	m.UpdatedByID = user.ID
	m.UpdatedByRole = user.Role
}

// GetUpdatedBy get updated by
func (m Model) GetUpdatedBy() User {
	//return model.UpdatedByID, model.CreatedByRole
	return User{
		ID:   m.UpdatedByID,
		Role: m.UpdatedByRole,
	}
}