package repositories

type Storage struct {
	UserRepository *UserRepository
}

func NewStorage() *Storage {
	newStorage := &Storage{
		UserRepository: &UserRepository{},
	}

	return newStorage
}
