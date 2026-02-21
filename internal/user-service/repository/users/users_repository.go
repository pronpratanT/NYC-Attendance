package repository

import (
	model "hr-program/internal/user-service/model/users"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	DB *gorm.DB
}

type UserRepositoryInterface interface {
	GetUserIDMapByEmployeeIDs(employeeIDs []string) (map[string]int64, error)
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) BulkInsert(data []model.Users) error {
	return r.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "employee_id"}},
			DoNothing: true,
		}).
		CreateInBatches(data, len(data)).Error
}

func (r *UserRepository) GetAllUsers() ([]model.Users, error) {
	sqlDB, err := r.DB.DB()
	if err != nil {
		log.Println("Failed to get raw DB connection:", err)
		return nil, err
	}
	rows, err := sqlDB.Query("SELECT id, employee_id, department_id, f_name, l_name, is_active, workday, created_at FROM users")
	if err != nil {
		log.Println("Failed to execute query:", err)
		return nil, err
	}
	defer rows.Close()

	var users []model.Users
	for rows.Next() {
		var user model.Users
		err := rows.Scan(
			&user.ID,
			&user.EmployeeID,
			&user.DepartmentID,
			&user.FName,
			&user.LName,
			&user.IsActive,
			&user.Workday,
			&user.CreatedAt,
		)
		if err != nil {
			log.Println("Failed to scan row:", err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) GetUserIDMapByEmployeeIDs(employeeIDs []string) (map[string]int64, error) {
	var users []model.Users

	err := r.DB.
		Select("id, employee_id").
		Where("employee_id IN ?", employeeIDs).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]int64)
	for _, u := range users {
		result[u.EmployeeID] = u.ID
	}

	return result, nil
}
