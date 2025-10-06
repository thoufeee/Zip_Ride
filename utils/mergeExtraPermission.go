package utils

import "strings"

// merge role and extra permissions
func MergePermissions(rolePerm []string, extraPerms []string) []string {
	perMap := make(map[string]bool)

	for _, p := range rolePerm {
		perMap[strings.ToUpper(p)] = true
	}

	for _, p := range extraPerms {
		perMap[strings.ToUpper(p)] = true
	}

	permissions := make([]string, 0, len(perMap))

	for p := range perMap {
		permissions = append(permissions, p)
	}

	return permissions
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
