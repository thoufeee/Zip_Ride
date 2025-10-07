package utils

import (
	"strings"
	"zipride/internal/models"
)

// convert slices of permission struct to slice of permission names
func PermissionToString(perms []models.Permission) []string {
	var names []string

	for _, p := range perms {
		names = append(names, p.Name)
	}

	return names
}

// merge role and extra permissions
func MergePermissions(a, b []string) []string {
	set := make(map[string]struct{})

	for _, p := range a {
		set[p] = struct{}{}
	}

	for _, p := range b {
		set[p] = struct{}{}
	}

	var merges []string

	for k := range set {
		merges = append(merges, k)
	}

	return merges
}

// check permission exists

func CheckPermission(permissions []string, perm string) bool {
	perm = strings.ToUpper(perm)

	for _, p := range permissions {

		if p == perm {
			return true
		}
	}

	return false
}
