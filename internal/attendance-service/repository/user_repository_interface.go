package repository

type UserRepositoryInterface interface {
	GetUserIDMapByEmployeeIDs(employeeIDs []string) (map[string]int64, error)
}
