package service

//type Repos interface {
//	TodoRepo
//}

type Services struct {
	repo TodoRepo
}

func NewService(repo TodoRepo) *Services {
	return &Services{repo: repo}
}
