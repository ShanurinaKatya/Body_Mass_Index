import type { FC } from 'react'
import { useEffect, useState } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Button } from 'react-bootstrap'
import { ServiceCard } from '../components/patients_categories'
import { ROUTES, SERVICES_MOCK } from '../modules/patients_categories'
import type { PatientsCategorie } from '../modules/patients_categories'

const MAX_AGE = 100

export const FeedPage: FC = () => {
  const [searchParams] = useSearchParams()
  const requestedId = Number(searchParams.get('id'))

  const [index, setIndex] = useState(() => {
    const found = SERVICES_MOCK.findIndex((service) => service.id === requestedId)
    return found === -1 ? 0 : found
  })
  const service = SERVICES_MOCK[index]
  const [showFullDesc, setShowFullDesc] = useState(false)
  const [videoSrc, setVideoSrc] = useState(
    service.video_url || '/Default.mp4',
  )

  useEffect(() => {
    setVideoSrc(service.video_url || '/Default.mp4')
    window.scrollTo(0, 0)
  }, [index, service.video_url])

  const handleVideoError = () => {
    if (videoSrc !== '/Default.mp4') {
      setVideoSrc('/Default.mp4')
    }
  }

  const handleNext = () => {
    setIndex((current) => (current + 1) % SERVICES_MOCK.length)
    setShowFullDesc(false)
  }

  return (
    <>
      <header className="page-header">
        <h1 className="page-title">ИМТ Калькулятор</h1>
      </header>

      <div className="feed-container">
        <video
          className="feed-video"
          src={videoSrc}
          autoPlay
          muted
          loop
          playsInline
          onError={handleVideoError}
        />

        <div className="feed-overlay">
          <div className="feed-sidebar">
            <div className="feed-sidebar-item">
              <div className="feed-sidebar-icon">
                <img
                  src="/free-icon-weight.png"
                  alt="Вес"
                  style={{ width: '20px', height: '20px', objectFit: 'contain' }}
                />
              </div>
              <span className="feed-sidebar-label">{service.weight} кг</span>
            </div>

            <div className="feed-sidebar-item">
              <div className="feed-sidebar-icon">
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <circle cx="12" cy="8" r="4" strokeWidth="1.5" fill="none" />
                  <path
                    d="M4 20C4 16.6863 7.58172 14 12 14C16.4183 14 20 16.6863 20 20"
                    strokeWidth="1.5"
                    strokeLinecap="round"
                    fill="none"
                  />
                </svg>
              </div>
              <span className="feed-sidebar-label">{service.age} лет</span>
            </div>

            <div className="feed-sidebar-item">
              <div className="feed-sidebar-icon feed-sidebar-icon--like">
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path
                    d="M12 21.35L10.55 20.03C5.4 15.36 2 12.28 2 8.5C2 5.42 4.42 3 7.5 3C9.24 3 10.91 3.81 12 5.09C13.09 3.81 14.76 3 16.5 3C19.58 3 22 5.42 22 8.5C22 12.28 18.6 15.36 13.45 20.04L12 21.35Z"
                    strokeWidth="1.5"
                    fill="none"
                  />
                </svg>
              </div>
              <span className="feed-sidebar-label">{service.likes}</span>
            </div>

            <div
              className="feed-sidebar-item feed-sidebar-item--clickable"
              onClick={handleNext}
            >
              <div className="feed-sidebar-icon feed-sidebar-icon--next">
                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path
                    d="M5 12H19M19 12L12 5M19 12L12 19"
                    stroke="white"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              </div>
              <span className="feed-sidebar-label">Далее</span>
            </div>
          </div>

          <div className="feed-bottom">
            <h1 className="feed-title">{service.name}</h1>
            <p
              className={`feed-description${showFullDesc ? ' feed-description-expanded' : ''}`}
            >
              {service.description}
            </p>
            <span
              className="feed-more-btn"
              onClick={() => setShowFullDesc((value) => !value)}
            >
              {showFullDesc ? 'меньше' : 'ещё'}
            </span>
          </div>
        </div>
      </div>
    </>
  )
}

export const TilePage: FC = () => {
  const navigate = useNavigate()

  const [age, setAge] = useState(MAX_AGE)
  const [appliedAge, setAppliedAge] = useState(MAX_AGE)
  const [deletedIds, setDeletedIds] = useState<number[]>([])

  // Как в лабе 1 (handler.CardsHandler): возраст 100 — все карточки,
  // иначе строгое совпадение age = N, совпадений нет — пусто.
  const matchesAge =
    appliedAge >= 0 && appliedAge < MAX_AGE
      ? (service: PatientsCategorie) => service.age === appliedAge
      : () => true

  const services = SERVICES_MOCK.filter(
    (service) => matchesAge(service) && !deletedIds.includes(service.id),
  )

  const handleFilter = () => setAppliedAge(age)

  const handleOpen = (id: number) => navigate(`${ROUTES.FEED}?id=${id}`)

  const handleDelete = (id: number) =>
    setDeletedIds((current) => [...current, id])

  return (
    <>
      <header className="page-header">
        <h1 className="page-title">Каталог расчётов</h1>
        <p className="page-subtitle">Специализированные калькуляторы ИМТ</p>

        <div className="cards-filter">
          <div className="cards-filter-main">
            <label className="cards-filter-label" htmlFor="ageSlider">
              Фильтр по возрасту
            </label>
            <div className="cards-filter-track">
              <input
                id="ageSlider"
                type="range"
                min={0}
                max={MAX_AGE}
                value={age}
                className="cards-filter-input"
                onChange={(event) => setAge(Number(event.target.value))}
              />
              <div className="cards-filter-range-labels">
                <span>0</span>
                <span
                  className="cards-filter-bubble"
                  style={{ left: `calc(12px + (100% - 24px) * ${age / 100})` }}
                >
                  {age}
                </span>
                <span>{MAX_AGE}</span>
              </div>
            </div>
          </div>
          <Button className="cards-filter-btn" onClick={handleFilter}>
            Найти
          </Button>
        </div>
      </header>

      <div className="cards-container">
        <div className="cards-grid">
          {services.map((service) => (
            <ServiceCard
              key={service.id}
              service={service}
              onOpen={handleOpen}
              onDelete={handleDelete}
            />
          ))}
        </div>
      </div>
    </>
  )
}
