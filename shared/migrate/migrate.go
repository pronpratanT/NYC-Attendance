package migrate

func AutoMigrate() {
	return db.AutoMigrate(
		&model.Departments{},
		&model.Employees{},
		&model.Attendance{},
		&model.AttendanceDaily{},
	)
}
