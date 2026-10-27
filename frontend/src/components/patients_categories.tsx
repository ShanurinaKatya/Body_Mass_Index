import type { FC } from 'react'
import { NavLink } from 'react-router-dom'
import { Card } from 'react-bootstrap'
import { ROUTES } from '../modules/patients_categories'
import type { PatientsCategorie } from '../modules/patients_categories'

export const TabBar: FC = () => (
  <nav className="tab-bar">
    <NavLink
      to={ROUTES.FEED}
      end
      className={({ isActive }) => (isActive ? 'tab-item active' : 'tab-item')}
    >
      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M3 9.5L12 3L21 9.5V20C21 20.55 20.55 21 20 21H4C3.45 21 3 20.55 3 20V9.5Z"
          stroke="currentColor"
          strokeWidth="1.5"
          strokeLinejoin="round"
        />
        <path d="M9 21V14H15V21" stroke="currentColor" strokeWidth="1.5" strokeLinejoin="round" />
      </svg>
      <span>Лента</span>
    </NavLink>

    <span className="tab-item tab-item--disabled" aria-disabled="true">
      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle cx="12" cy="12" r="9" stroke="currentColor" strokeWidth="1.5" />
        <path d="M12 8V16M8 12H16" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      </svg>
      <span>Добавление</span>
    </span>

    <NavLink
      to={ROUTES.TILE}
      className={({ isActive }) => (isActive ? 'tab-item active' : 'tab-item')}
    >
      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <rect x="3" y="3" width="8" height="8" rx="1.5" stroke="currentColor" strokeWidth="1.5" />
        <rect x="13" y="3" width="8" height="8" rx="1.5" stroke="currentColor" strokeWidth="1.5" />
        <rect x="3" y="13" width="8" height="8" rx="1.5" stroke="currentColor" strokeWidth="1.5" />
        <rect x="13" y="13" width="8" height="8" rx="1.5" stroke="currentColor" strokeWidth="1.5" />
      </svg>
      <span>Плитка</span>
    </NavLink>

    <span className="tab-item tab-item--disabled" aria-disabled="true">
      <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle cx="12" cy="8" r="4" stroke="currentColor" strokeWidth="1.5" fill="none" />
        <path
          d="M4 20C4 16.6863 7.58172 14 12 14C16.4183 14 20 16.6863 20 20"
          stroke="currentColor"
          strokeWidth="1.5"
          strokeLinecap="round"
          fill="none"
        />
      </svg>
      <span>Вход</span>
    </span>
  </nav>
)

interface ServiceCardProps {
  service: PatientsCategorie
  onOpen: (id: number) => void
  onDelete: (id: number) => void
}

export const ServiceCard: FC<ServiceCardProps> = ({ service, onOpen, onDelete }) => (
  <Card className="card">
    <div className="card-link" onClick={() => onOpen(service.id)}>
      <Card.Img
        className="card-image"
        variant="top"
        src={service.image_url || '/DefaultImage.jpg'}
        alt={service.name}
        onError={(event) => {
          event.currentTarget.src = '/DefaultImage.jpg'
        }}
      />
      <div className="card-content">
        <Card.Title className="card-title">{service.name}</Card.Title>
        <div className="card-row">
          <div className="card-likes">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path
                d="M12 21.35L10.55 20.03C5.4 15.36 2 12.28 2 8.5C2 5.42 4.42 3 7.5 3C9.24 3 10.91 3.81 12 5.09C13.09 3.81 14.76 3 16.5 3C19.58 3 22 5.42 22 8.5C22 12.28 18.6 15.36 13.45 20.04L12 21.35Z"
                stroke="#e43f44"
                strokeWidth="1.5"
                fill="none"
              />
            </svg>
            <span style={{ color: '#333333' }}>{service.likes}</span>
          </div>
          <span className="card-meta">{service.weight} кг</span>
          <span className="card-meta">{service.age} лет</span>
        </div>
      </div>
    </div>
    <form
      className="card-delete"
      onSubmit={(event) => {
        event.preventDefault()
        onDelete(service.id)
      }}
    >
      <button type="submit" className="card-delete-btn">
        Удалить
      </button>
    </form>
  </Card>
)
