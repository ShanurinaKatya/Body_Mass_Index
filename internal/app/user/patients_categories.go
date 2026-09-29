package user

type singleton struct {
	ID uint
}

var instance *singleton

func Get() *singleton {
	if instance == nil {
		instance = &singleton{
			ID: 1,
		}
	}
	return instance
}