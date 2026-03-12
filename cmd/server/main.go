package main

import (
	attrepo "hr-program/internal/attendance-service/repository"
	reqrepo "hr-program/internal/request-service/repository"
	deprepo "hr-program/internal/user-service/repository/departments"
	shfrepo "hr-program/internal/user-service/repository/shifts"
	usrrepo "hr-program/internal/user-service/repository/users"

	attroute "hr-program/internal/attendance-service/app/router"
	atthandler "hr-program/internal/attendance-service/handler"
	reqroute "hr-program/internal/request-service/app/router"
	reqhandler "hr-program/internal/request-service/handler"
	usrroute "hr-program/internal/user-service/app/router"
	usrhandler "hr-program/internal/user-service/handler"

	"hr-program/shared/config"
	db "hr-program/shared/connection"

	attservice "hr-program/internal/attendance-service/service"
	reqservice "hr-program/internal/request-service/service"
	usrservice "hr-program/internal/user-service/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	// Load ENV
	config.LoadConfig()

	// Connect DBs
	appDB := db.ConnectDB()
	cloudDB := db.ConnectCloudtime()
	econsDB := db.ConnectEcons()
	sqlExpressDB := db.ConnectSQLExpress()

	redisClient := db.ConnectRedis()
	_ = redisClient // Currently not used in this service, but connected for future use

	// Init repositories for attendance service
	attAppRepo := attrepo.NewAttendanceRepository(appDB)
	attCloudRepo := attrepo.NewCloudtimeRepository(cloudDB)
	// Init repositories for user service
	usrAppRepo := usrrepo.NewUserRepository(appDB)
	usrCloudRepo := usrrepo.NewCloudtimeUserRepository(cloudDB)
	// Init repositories for user service - departments
	depAppRepo := deprepo.NewDepartmentsRepository(appDB)
	depCloudRepo := deprepo.NewCloudtimeDepartmentsRepository(cloudDB)
	// Init repositories for user service - shifts
	shiftAppRepo := shfrepo.NewShiftsRepository(appDB)
	shiftSqlExpressRepo := shfrepo.NewSQLExpressShiftRepository(sqlExpressDB)

	// Init repositories for request service - OT
	otAppRepo := reqrepo.NewOTRepository(appDB)
	holidayRepo := reqrepo.NewHolidayRepository(appDB)
	econsRepo := reqrepo.NewEconsRepository(econsDB)

	// Init services
	attendanceService := attservice.NewAttendanceService(attCloudRepo, attAppRepo, usrAppRepo, shiftAppRepo, otAppRepo, holidayRepo)
	userService := usrservice.NewUserService(usrCloudRepo, usrAppRepo, depAppRepo, depCloudRepo, shiftSqlExpressRepo, shiftAppRepo)
	requestService := reqservice.NewRequestService(otAppRepo, econsRepo, usrAppRepo, holidayRepo)

	// handler + router
	attendanceHandler := atthandler.NewAttendanceHandler(attendanceService)
	userHandler := usrhandler.NewUserHandler(userService)
	reqHandler := reqhandler.NewRequestHandler(requestService)

	r := gin.Default()

	attroute.AttendanceRouter(r, attendanceHandler)
	usrroute.UserRouter(r, userHandler)
	reqroute.RequestRouter(r, reqHandler)

	// Initial sync เบื้องหลังครั้งแรกตอน start service
	go func() {
		if err := attendanceService.SyncFullLoadAttendance(); err != nil {
			log.Println("Initial sync attendance failed:", err)
		}
		if err := userService.SyncFullLoadUsers(); err != nil {
			log.Println("Initial sync users failed:", err)
		}
		if err := userService.SyncFullLoadDeps(); err != nil {
			log.Println("Initial sync departments failed:", err)
		}
		if err := attendanceService.GenerateAndSaveAttendanceDaily(); err != nil {
			log.Println("Initial process attendance daily failed:", err)
		}
		if err := requestService.SyncFullLoadOT(); err != nil {
			log.Println("Initial sync OT requests failed:", err)
		}
		if err := requestService.GenerateAndSaveOT(); err != nil {
			log.Println("Initial process OT logs to OT docs failed:", err)
		}
		if err := requestService.SyncHolidays(); err != nil {
			log.Println("Initial sync holidays failed:", err)
		}
		// กะการทำงานจาก SQL Express ข้อมูลจาก Bplus
		// if err := userService.GenerateAndSaveShifts(); err != nil {
		// 	log.Println("Initial process shifts failed:", err)
		// }
		// UserShifts กะการทำงานของพนักงาน
		if err := userService.ProcessUserShifts(); err != nil {
			log.Println("Initial process user shifts failed:", err)
		}
		log.Println("Initial sync completed successfully")
	}()

	// Scheduler รัน sync + generate attendance_daily ซ้ำทุก ๆ 10 นาที โดยไม่ต้อง restart container
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := attendanceService.SyncFullLoadAttendance(); err != nil {
				log.Println("Scheduled sync attendance failed:", err)
				continue
			}
			if err := attendanceService.GenerateAndSaveAttendanceDaily(); err != nil {
				log.Println("Scheduled process attendance daily failed:", err)
				continue
			}
			if err := requestService.SyncFullLoadOT(); err != nil {
				log.Println("Scheduled sync OT requests failed:", err)
				continue
			}
			if err := requestService.GenerateAndSaveOT(); err != nil {
				log.Println("Scheduled process OT logs to OT docs failed:", err)
				continue
			}
			// if err := syncService.SyncHolidays(); err != nil {
			// 	log.Println("Scheduled sync holidays failed:", err)
			// 	continue
			// }
			log.Println("Scheduled sync + attendance daily completed")
		}
	}()

	r.Run(":8080")
}
