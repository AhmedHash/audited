package audited

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strconv"

	"github.com/jinzhu/gorm"
)

type auditableInterface interface {
	SetCreatedBy(User)
	GetCreatedBy() User
	SetUpdatedBy(User)
	GetUpdatedBy() User
}

func isAuditable(scope *gorm.Scope) (isAuditable bool) {
	if scope.GetModelStruct().ModelType == nil {
		return false
	}
	_, isAuditable = reflect.New(scope.GetModelStruct().ModelType).Interface().(auditableInterface)
	return
}

func getCurrentUser(scope *gorm.Scope) (currentUser User, ok bool) {
	var user interface{}
	var hasUser bool
	var err error

	user, hasUser = scope.DB().Get("audited:current_user")

	if hasUser {
		if userID, ok := scope.New(user).FieldByName("id"); ok {
			id := fmt.Sprintf("%v", userID.Field.Interface())
			currentUser.ID, err =  uuid.FromString(id)
			if err != nil {
				return currentUser, false
			}
		}
		if userRole, ok := scope.New(user).FieldByName("role"); ok {
			role :=  fmt.Sprintf("%v", userRole.Field.Interface())
			currentUser.Role, err = strconv.ParseInt(role, 10, 64)
			if err != nil {
				return currentUser, false
			}
		}
		return currentUser, true
	}

	return currentUser, false
}

func assignCreatedBy(scope *gorm.Scope) {
	if isAuditable(scope) {
		if user, ok := getCurrentUser(scope); ok {
			scope.SetColumn("CreatedByID", user.ID)
			scope.SetColumn("CreatedByRole", user.Role)
		}
	}
}

func assignUpdatedBy(scope *gorm.Scope) {
	if isAuditable(scope) {
		if user, ok := getCurrentUser(scope); ok {
			if attrs, ok := scope.InstanceGet("gorm:update_attrs"); ok {
				updateAttrs := attrs.(map[string]interface{})
				updateAttrs["updated_by_id"] = user.ID
				updateAttrs["updated_by_role"] = user.Role
				scope.InstanceSet("gorm:update_attrs", updateAttrs)
			} else {
				scope.SetColumn("UpdatedByID", user.ID)
				scope.SetColumn("UpdatedByRole", user.Role)
			}
		}
	}
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *gorm.DB) {
	callback := db.Callback()
	if callback.Create().Get("audited:assign_created_by") == nil {
		callback.Create().After("gorm:before_create").Register("audited:assign_created_by", assignCreatedBy)
	}
	if callback.Update().Get("audited:assign_updated_by") == nil {
		callback.Update().After("gorm:before_update").Register("audited:assign_updated_by", assignUpdatedBy)
	}
}