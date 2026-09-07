package repository

import "fmt"

type Service struct {
	ID          int
	Title       string
	Description string
	ImageKey    string
	VideoKey    string
	Weight      float64
	Height      float64
	Age         int
	Gender      string
	Likes       []int
	Status      string
}

type Repository struct {
	services []Service
}

func NewRepository() (*Repository, error) {
	services := []Service{
		{
			ID:          1,
			Title:       "ИМТ для девочек",
			Description: "Оценка индекса массы тела для девочек подросткового возраста. Учитывает особенности роста и развития.",
			ImageKey:    "girl_foto.jpg",
			VideoKey:    "Girl.mp4",
			Weight:      15,
			Height:      140,
			Age:         7,
			Gender:      "Женский",
			Likes:       []int{1},
			Status:      "published",
		},
		{
			ID:          2,
			Title:       "ИМТ для мальчиков",
			Description: "Расчёт индекса массы тела для мальчиков с учётом возрастных норм.",
			ImageKey:    "boy_foto.jpg",
			VideoKey:    "Boy.mp4",
			Weight:      20,
			Height:      150,
			Age:         10,
			Gender:      "Мужской",
			Likes:       []int{1, 2},
			Status:      "published",
		},
		{
			ID:          3,
			Title:       "ИМТ для женщин",
			Description: "Для женщин после менопаузы классическая норма ИМТ (18,5–24,9) уже не актуальна. Оптимальный диапазон — 22–27, он компенсирует потерю мышц и костной массы и снижает риски остеопороза.",
			ImageKey:    "woman_foto.avif",
			VideoKey:    "Woman.mp4",
			Weight:      75,
			Height:      165,
			Age:         35,
			Gender:      "Женский",
			Likes:       []int{1, 2},
			Status:      "published",
		},
		{
			ID:          4,
			Title:       "ИМТ для мужчин",
			Description: "Оценка соотношения массы тела и роста у мужчин разного возраста.",
			ImageKey:    "man_foto.avif",
			VideoKey:    "Man.mp4",
			Weight:      80,
			Height:      178,
			Age:         30,
			Gender:      "Мужской",
			Likes:       []int{1},
			Status:      "published",
		},
		{
			ID:          5,
			Title:       "ИМТ для пожилых мужчин",
			Description: "Для мужчин старше 65 лет нормы ИМТ отличаются от стандартных. Рекомендуется поддерживать ИМТ в диапазоне 23–28.",
			ImageKey:    "old_man_foto.png",
			VideoKey:    "Old_man.mp4",
			Weight:      78,
			Height:      170,
			Age:         78,
			Gender:      "Мужской",
			Likes:       []int{1, 2},
			Status:      "published",
		},
		{
			ID:          6,
			Title:       "ИМТ для пожилых женщин",
			Description: "Для женщин старше 65 лет классические нормы ИМТ требуют корректировки. Небольшой избыток массы тела может быть защитным фактором.",
			ImageKey:    "woman_foto.avif",
			VideoKey:    "Woman.mp4",
			Weight:      72,
			Height:      160,
			Age:         75,
			Gender:      "Женский",
			Likes:       []int{},
			Status:      "draft",
		},
	}

	return &Repository{services: services}, nil
}

func (r *Repository) GetServices() ([]Service, error) {
	var result []Service
	for _, s := range r.services {
		if s.Status == "published" {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *Repository) GetServiceByID(id int) (Service, error) {
	for _, s := range r.services {
		if s.ID == id && s.Status != "deleted" {
			return s, nil
		}
	}
	return Service{}, fmt.Errorf("service not found")
}

func (r *Repository) GetNextService(id int) (Service, error) {
	published := r.getPublished()
	for i, s := range published {
		if s.ID == id {
			nextIdx := (i + 1) % len(published)
			return published[nextIdx], nil
		}
	}
	if len(published) > 0 {
		return published[0], nil
	}
	return Service{}, fmt.Errorf("no services found")
}

func (r *Repository) GetDraft() (Service, error) {
	for _, s := range r.services {
		if s.Status == "draft" {
			return s, nil
		}
	}
	return Service{}, fmt.Errorf("draft not found")
}

func (r *Repository) GetServicesByAge(age int) ([]Service, error) {
	var result []Service
	for _, s := range r.services {
		if s.Status == "published" && s.Age == age {
			result = append(result, s)
		}
	}
	return result, nil
}

func (r *Repository) getPublished() []Service {
	var result []Service
	for _, s := range r.services {
		if s.Status == "published" {
			result = append(result, s)
		}
	}
	return result
}

func (r *Repository) GetAllPublished() []Service {
	return r.getPublished()
}
