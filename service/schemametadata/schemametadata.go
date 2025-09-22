package schemametadata

type schemaRepo interface {
	SaveSchema(Schema) error
}

type Service struct {
	repo schemaRepo
}

func New(repo schemaRepo) Service {
	return Service{
		repo: repo,
	}
}

func (m Service) Save(schema Schema) error {
	return m.repo.SaveSchema(schema)
}
