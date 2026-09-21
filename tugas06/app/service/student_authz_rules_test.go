package service

import (
	"testing"

	"tugas06/app/model"
	"tugas06/helper"
)

func TestCanAccessStudent(t *testing.T) {
	perms := helper.NewPermissionSet(map[string][]string{
		"admin": {"student:read:any"},
		"user":  {},
	})

	admin := model.AuthUser{StudentID: 1, Role: "admin"}
	user := model.AuthUser{StudentID: 5, Role: "user"}

	owner := 5
	var noOwner *int = nil
	otherOwner := 99

	cases := []struct {
		name     string
		current  model.AuthUser
		targetID int
		ownerID  *int
		want     bool
	}{
		{"user akses dirinya sendiri", user, 5, &owner, true},
		{"user akses student yang dia miliki", user, 7, &owner, true},
		{"user akses milik orang lain", user, 8, &otherOwner, false},
		{"admin akses siapa pun", admin, 8, nil, true},
		{"owner nil & bukan self & tidak punya permission", user, 9, noOwner, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := CanAccessStudent(tc.current, tc.targetID, tc.ownerID, perms, "student:read:any")
			if got != tc.want {
				t.Errorf("harap %v, dapat %v", tc.want, got)
			}
		})
	}
}

func TestValidateAssignRole(t *testing.T) {
	perms := helper.NewPermissionSet(map[string][]string{
		"admin": {"role:assign"},
		"staff": {},
		"user":  {},
	})

	admin := model.AuthUser{StudentID: 1, Role: "admin"}

	errs := ValidateAssignRole(admin, 5, model.AssignRoleRequest{Role: "staff"}, perms)
	if len(errs) != 0 {
		t.Errorf("harap lolos, dapat: %v", errs)
	}

	errs = ValidateAssignRole(admin, 1, model.AssignRoleRequest{Role: "user"}, perms)
	if _, ok := errs["role"]; !ok {
		t.Error("harap error role saat mengubah role diri sendiri")
	}

	errs = ValidateAssignRole(admin, 5, model.AssignRoleRequest{Role: "superadmin"}, perms)
	if _, ok := errs["role"]; !ok {
		t.Error("harap error role untuk role tidak dikenal")
	}
}