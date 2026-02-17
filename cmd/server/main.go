package main

import (
	"log"

	"github.com/pronpratanT/leave-system/internal/attendance-service/repository"
	"github.com/pronpratanT/leave-system/internal/attendance-service/service"
	"github.com/pronpratanT/leave-system/shared/config"
	db "github.com/pronpratanT/leave-system/shared/connection"
)

func main() {

	// Load ENV
	config.LoadConfig()

	// Connect DBs
	appDB := db.ConnectDB()
	cloudDB := db.ConnectCloudtime()

	// Init repositories
	appRepo := repository.NewAttendanceRepository(appDB)
	cloudRepo := repository.NewCloudtimeRepository(cloudDB)

	// Init service
	attendanceService := service.NewAttendanceService(cloudRepo, appRepo)

	// Run sync (2 worker parallel)
	err := attendanceService.SyncFullLoad()
	if err != nil {
		log.Fatal("Sync failed:", err)
	}

	log.Println("Sync completed successfully")
}
