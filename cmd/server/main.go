package main

import (
	"hr-program/internal/attendance-service/app/router"
	"hr-program/internal/attendance-service/handler"
	attrepo "hr-program/internal/attendance-service/repository"
	deprepo "hr-program/internal/user-service/repository/departments"
	usrrepo "hr-program/internal/user-service/repository/users"

	"hr-program/shared/config"
	db "hr-program/shared/connection"

	attservice "hr-program/internal/attendance-service/service"
	usrservice "hr-program/internal/user-service/service"
	"log"
	"time"
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
	attendanceService := attservice.NewAttendanceService(attCloudRepo, attAppRepo, usrAppRepo)
	userService := usrservice.NewUserService(usrCloudRepo, usrAppRepo)
	departmentService := usrservice.NewDepartmentsService(depCloudRepo, depAppRepo)

	// handler + router
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	r := router.AttendanceRouter(attendanceHandler)

	// Initial sync เบื้องหลังครั้งแรกตอน start service
	go func() {
		if err := attendanceService.SyncFullLoad(); err != nil {
			log.Println("Initial sync attendance failed:", err)
		}
		if err := userService.SyncFullLoad(); err != nil {
			log.Println("Initial sync users failed:", err)
		}
		if err := departmentService.SyncFullLoad(); err != nil {
			log.Println("Initial sync departments failed:", err)
		}
		if err := attendanceService.GenerateAndSaveAttendanceDaily(); err != nil {
			log.Println("Initial process attendance daily failed:", err)
		}
		log.Println("Initial sync completed successfully")
	}()

	// Scheduler รัน sync + generate attendance_daily ซ้ำทุก ๆ 5 นาที โดยไม่ต้อง restart container
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := attendanceService.SyncFullLoad(); err != nil {
				log.Println("Scheduled sync attendance failed:", err)
				continue
			}
			if err := attendanceService.GenerateAndSaveAttendanceDaily(); err != nil {
				log.Println("Scheduled process attendance daily failed:", err)
				continue
			}
			log.Println("Scheduled sync + attendance daily completed")
		}
	}()

	r.Run(":8080")
}
