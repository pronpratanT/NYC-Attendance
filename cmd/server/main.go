package main

import (
	"hr-program/internal/attendance-service/handler"
	attrepo "hr-program/internal/attendance-service/repository"
	deprepo "hr-program/internal/user-service/repository/departments"
	usrrepo "hr-program/internal/user-service/repository/users"

	"hr-program/shared/config"
	db "hr-program/shared/connection"

	attservice "hr-program/internal/attendance-service/service"
	usrservice "hr-program/internal/user-service/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	// Load ENV
	config.LoadConfig()

	// Connect DBs
	appDB := db.ConnectDB()
	cloudDB := db.ConnectCloudtime()
	r := gin.Default()

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

	// handler
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	handler.RegisterAttendanceRoutes(r, attendanceHandler)

	// Run sync (2 worker parallel) เบื้องหลัง
	go func() {
		if err := attendanceService.SyncFullLoad(); err != nil {
			log.Println("Sync attendance failed:", err)
		}
		if err := userService.SyncFullLoad(); err != nil {
			log.Println("Sync users failed:", err)
		}
		if err := departmentService.SyncFullLoad(); err != nil {
			log.Println("Sync departments failed:", err)
		}
		log.Println("Sync completed successfully")
	}()

	r.Run(":8080")
}
