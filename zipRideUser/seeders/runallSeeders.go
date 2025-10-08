package seeders

// for running all seeders

func RunAllSeeders() {
	SeedPermisiions()
	SeedRoles()
	SeedRolePermissions()
	SeedAdmin()
}
