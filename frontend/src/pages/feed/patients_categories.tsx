import type { FC } from 'react'
import { useEffect, useState } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useSearchParams } from 'react-router-dom'

import type { AppDispatch, RootState } from '../../app/store/patients_categories'
import {
  fetchFeed,
  likeService,
} from '../../entities/service/model/patients_categories'
import { NextService } from '../../features/next-service/patients_categories'
import { ToggleDescription } from '../../features/toggle-description/patients_categories'
import { toMediaUrl } from '../../shared/lib/media/patients_categories'
import type { PatientsCategorie } from '../../shared/model/patients_categories'

interface FeedContentProps {
  feed: PatientsCategorie
  isAuthenticated: boolean
  onLike: () => void
  onNext: () => void
}

const FeedContent: FC<FeedContentProps> = ({
  feed,
  isAuthenticated,
  onLike,
  onNext,
}) => {
  const [showFullDesc, setShowFullDesc] = useState(false)
  const [videoError, setVideoError] = useState(false)

  const videoSrc = videoError
    ? '/Default.mp4'
    : toMediaUrl(feed.video_url) || '/Default.mp4'

  return (
    <div className="feed-container">
      <video
        className="feed-video"
        src={videoSrc}
        autoPlay
        muted
        loop
        playsInline
        onError={() => setVideoError(true)}
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
            <span className="feed-sidebar-label">{feed.weight} кг</span>
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
            <span className="feed-sidebar-label">{feed.age} лет</span>
          </div>

          <div
            className={`feed-sidebar-item feed-sidebar-item--like${
              isAuthenticated ? ' feed-sidebar-item--clickable' : ''
            }`}
            onClick={isAuthenticated ? onLike : undefined}
          >
            <div className="feed-sidebar-icon feed-sidebar-icon--like">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path
                  className={feed.is_liked ? 'feed-like-active' : undefined}
                  d="M12 21.35L10.55 20.03C5.4 15.36 2 12.28 2 8.5C2 5.42 4.42 3 7.5 3C9.24 3 10.91 3.81 12 5.09C13.09 3.81 14.76 3 16.5 3C19.58 3 22 5.42 22 8.5C22 12.28 18.6 15.36 13.45 20.04L12 21.35Z"
                  strokeWidth="1.5"
                />
              </svg>
            </div>
            <span className="feed-sidebar-label">{feed.likes}</span>
          </div>

          <NextService onClick={onNext} />
        </div>

        <div className="feed-bottom">
          <h1 className="feed-title">{feed.name}</h1>
          <p
            className={`feed-description${showFullDesc ? ' feed-description-expanded' : ''}`}
          >
            {feed.description}
          </p>
          <ToggleDescription
            expanded={showFullDesc}
            onToggle={() => setShowFullDesc((value) => !value)}
          />
        </div>
      </div>
    </div>
  )
}

export const FeedPage: FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const [searchParams, setSearchParams] = useSearchParams()

  // Отсутствующий, пустой или некорректный id означает «первая услуга».
  // Number(null) === 0, поэтому просто Number() использовать нельзя.
  const idParam = searchParams.get('id')
  const parsedId = Number(idParam)
  const requestedId = idParam && parsedId > 0 ? parsedId : NaN

  const isAuthenticated = useSelector(
    (state: RootState) => state.user.isAuthenticated,
  )
  const { feed, error } = useSelector(
    (state: RootState) => state.services,
  )

  useEffect(() => {
    dispatch(
      fetchFeed(Number.isNaN(requestedId) ? {} : { id: requestedId }),
    )
  }, [dispatch, requestedId])

  const handleNext = () => {
    if (!feed) {
      return
    }
    window.scrollTo(0, 0)
    dispatch(fetchFeed({ id: feed.id, next: true }))
      .unwrap()
      .then((result) => {
        setSearchParams({ id: String(result.item.id) })
      })
      .catch(() => undefined)
  }

  const handleLike = () => {
    if (!feed) {
      return
    }
    dispatch(likeService({ id: feed.id, like: feed.is_liked ? 0 : 1 }))
  }

  return (
    <>
      <header className="page-header">
        <h1 className="page-title">ИМТ Калькулятор</h1>
        {error ? <p className="page-error">{error}</p> : null}
      </header>

      {feed ? (
        <FeedContent
          key={feed.id}
          feed={feed}
          isAuthenticated={isAuthenticated}
          onLike={handleLike}
          onNext={handleNext}
        />
      ) : null}
    </>
  )
}
