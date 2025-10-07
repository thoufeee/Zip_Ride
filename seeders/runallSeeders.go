package seeders

// for running all seeders

func RunAllSeeders() {
	SeedAdmin()
	SeedPermisiions()
	SeedRolePermissions()
	SeedRoles()
}
