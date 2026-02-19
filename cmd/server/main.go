package main

import (
	attrepo "hr-program/internal/attendance-service/repository"
	deprepo "hr-program/internal/user-service/repository/departments"
	usrrepo "hr-program/internal/user-service/repository/users"

	"hr-program/shared/config"
	db "hr-program/shared/connection"

	attservice "hr-program/internal/attendance-service/service"
	usrservice "hr-program/internal/user-service/service"
	"log"
)

func main() {

	// Load ENV
	config.LoadConfig()

	// Connect DBs
	appDB := db.ConnectDB()
	cloudDB := db.ConnectCloudtime()

	// Init repositories for attendance service
	attAppRepo := attrepo.NewAttendanceRepository(appDB)
	attCloudRepo := attrepo.NewCloudtimeRepository(cloudDB)
	// Init repositories for user service
	usrAppRepo := usrrepo.NewUserRepository(appDB)
	usrCloudRepo := usrrepo.NewCloudtimeUserRepository(cloudDB)
	// Init repositories for user service - departments
	depAppRepo := deprepo.NewDepartmentsRepository(appDB)
	depCloudRepo := deprepo.NewCloudtimeDepartmentsRepository(cloudDB)

	// Init services
	attendanceService := attservice.NewAttendanceService(attCloudRepo, attAppRepo)
	userService := usrservice.NewUserService(usrCloudRepo, usrAppRepo)
	departmentService := usrservice.NewDepartmentsService(depCloudRepo, depAppRepo)

	// Run sync (2 worker parallel)
	err := attendanceService.SyncFullLoad()
	if err != nil {
		log.Fatal("Sync attendance failed:", err)
	}

	err = userService.SyncFullLoad()
	if err != nil {
		log.Fatal("Sync users failed:", err)
	}

	err = departmentService.SyncFullLoad()
	if err != nil {
		log.Fatal("Sync departments failed:", err)
	}

	log.Println("Sync completed successfully")
}
