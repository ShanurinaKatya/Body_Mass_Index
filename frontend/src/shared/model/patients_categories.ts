import { MAX_AGE } from '../config/patients_categories'

export interface PatientsCategorie {
  id: number
  name: string
  description: string
  image_url: string
  video_url: string
  age: number
  weight: number
  likes: number
  creator_id: number
  created_at: string
  updated_at: string
  is_creator?: number
  is_liked?: number
}

// Mock-данные: используются только гостем, когда бэкенд недоступен.
// Поле image_url/video_url пустое — применяется дефолтное медиа проекта.
export const SERVICES_MOCK: PatientsCategorie[] = [
  {
    id: 1,
    name: 'ИМТ для девочек',
    description:
      'Оценка индекса массы тела для девочек подросткового возраста. Учитывает особенности роста и развития.',
    image_url: '',
    video_url: '',
    age: 8,
    weight: 17.5,
    likes: 2,
    creator_id: 0,
    created_at: '2026-09-14T09:08:35Z',
    updated_at: '2026-09-14T09:08:35Z',
  },
  {
    id: 2,
    name: 'ИМТ для мальчиков',
    description:
      'Расчёт индекса массы тела для мальчиков с учётом возрастных норм.',
    image_url: '',
    video_url: '',
    age: 10,
    weight: 20,
    likes: 2,
    creator_id: 0,
    created_at: '2026-09-14T09:08:35Z',
    updated_at: '2026-09-14T09:08:35Z',
  },
  {
    id: 3,
    name: 'ИМТ для женщин',
    description:
      'Для женщин после менопаузы классическая норма ИМТ (18,5–24,9) уже не актуальна. Оптимальный диапазон — 22–27, он компенсирует потерю мышц и костной массы и снижает риски остеопороза.',
    image_url: '',
    video_url: '',
    age: 35,
    weight: 75,
    likes: 2,
    creator_id: 0,
    created_at: '2026-09-14T09:08:35Z',
    updated_at: '2026-09-14T09:08:35Z',
  },
  {
    id: 4,
    name: 'ИМТ для мужчин',
    description:
      'Оценка соотношения массы тела и роста у мужчин разного возраста.',
    image_url: '',
    video_url: '',
    age: 30,
    weight: 80,
    likes: 1,
    creator_id: 0,
    created_at: '2026-09-14T09:08:35Z',
    updated_at: '2026-09-14T09:08:35Z',
  },
  {
    id: 5,
    name: 'ИМТ для пожилых мужчин',
    description:
      'Для мужчин старше 65 лет нормы ИМТ отличаются от стандартных. Рекомендуется поддерживать ИМТ в диапазоне 23–28.',
    image_url: '',
    video_url: '',
    age: 78,
    weight: 78,
    likes: 2,
    creator_id: 0,
    created_at: '2026-09-14T09:08:35Z',
    updated_at: '2026-09-14T09:08:35Z',
  },
  {
    id: 8,
    name: 'ИТМ для женщин 60+',
    description: 'Обновлённое описание',
    image_url: '',
    video_url: '',
    age: 75,
    weight: 80.5,
    likes: 1,
    creator_id: 0,
    created_at: '2026-09-14T09:43:33Z',
    updated_at: '2026-09-14T09:43:33Z',
  },
  {
    id: 24,
    name: 'ИМТ для учеников',
    description:
      'Расчёт индекса массы тела для учеников с учётом состава тела.',
    image_url: '',
    video_url: '',
    age: 17,
    weight: 62.5,
    likes: 0,
    creator_id: 0,
    created_at: '2026-09-29T09:38:09Z',
    updated_at: '2026-09-29T09:38:09Z',
  },
  {
    id: 25,
    name: 'ИМТ тест',
    description: 'Новое описание',
    image_url: '',
    video_url: '',
    age: 18,
    weight: 62,
    likes: 0,
    creator_id: 0,
    created_at: '2026-10-10T19:53:01Z',
    updated_at: '2026-10-10T19:53:01Z',
  },
]

export const getMockTiles = (age: number): PatientsCategorie[] =>
  age < MAX_AGE
    ? SERVICES_MOCK.filter((service) => service.age === age)
    : SERVICES_MOCK

export const getMockFeed = (id?: number, next?: boolean): PatientsCategorie => {
  if (id === undefined || Number.isNaN(id)) {
    return SERVICES_MOCK[0]
  }
  const index = SERVICES_MOCK.findIndex((service) => service.id === id)
  if (index === -1) {
    return SERVICES_MOCK[0]
  }
  return next
    ? SERVICES_MOCK[(index + 1) % SERVICES_MOCK.length]
    : SERVICES_MOCK[index]
}
