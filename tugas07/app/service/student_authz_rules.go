package service

import (
	"strings"

	"tugas07/app/model"
	"tugas07/helper"
)

func CanAccessStudent(
	current model.AuthUser,
	targetID int,
	ownerID *int,
	perms *helper.PermissionSet,
	anyPermission string,
) bool {
	if current.StudentID == targetID {
		return true
	}
	if ownerID != nil && *ownerID == current.StudentID {
		return true
	}
	return perms.Can(current.Role, anyPermission)
}

func ValidateAssignRole(
	current model.AuthUser,
	targetID int,
	req model.AssignRoleRequest,
	perms *helper.PermissionSet,
) map[string]string {
	errs := map[string]string{}

	role := strings.TrimSpace(req.Role)
	if role == "" {
		errs["role"] = "wajib diisi"
		return errs
	}

	if !perms.IsKnownRole(role) {
		errs["role"] = "role tidak dikenal, pilih salah satu dari: " +
			strings.Join(perms.KnownRoles(), ", ")
	}

	if current.StudentID == targetID {
		errs["role"] = "tidak boleh mengubah role diri sendiri"
	}

	return errs
}
