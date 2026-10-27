import type { FC } from 'react'
import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { SERVICES_MOCK } from '../../shared/model/patients_categories'
import { NextService } from '../../features/next-service/patients_categories'
import { ToggleDescription } from '../../features/toggle-description/patients_categories'

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

            <NextService onClick={handleNext} />
          </div>

          <div className="feed-bottom">
            <h1 className="feed-title">{service.name}</h1>
            <p
              className={`feed-description${showFullDesc ? ' feed-description-expanded' : ''}`}
            >
              {service.description}
            </p>
            <ToggleDescription
              expanded={showFullDesc}
              onToggle={() => setShowFullDesc((value) => !value)}
            />
          </div>
        </div>
      </div>
    </>
  )
}
