package repository

import (
	"context"
	"database/sql"
	"hr-backend/internal/attendance-service/model"
)

type AttendanceRepository struct {
	DB *sql.DB
}

func NewAttendanceRepository(db *sql.DB) *AttendanceRepository {
	return &AttendanceRepository{DB: db}
}

func (r *AttendanceRepository) Insert(ctx context.Context, att model.Attendance) error {
	query := `
	INSERT INTO attendance (
		bh,
		user_serial,
		user_no,
		user_lname,
		dep_no,
		user_dep,
		user_depname,
		user_type,
		user_card,
		sj,
		iden,
		fx,
		jlzp_serial,
		dev_serial,
		mc,
		health_status
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.DB.ExecContext(ctx, query,
		att.BH,
		att.UserSerial,
		att.UserNo,
		att.UserLName,
		att.DepNo,
		att.UserDep,
		att.UserDepName,
		att.UserType,
		att.UserCard,
		att.SJ,
		att.Iden,
		att.FX,
		att.JLZPSerial,
		att.DevSerial,
		att.MC,
		att.HealthStatus,
	)

	return err
}
