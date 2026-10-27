import type { FC } from 'react'
import { Card } from 'react-bootstrap'
import type { PatientsCategorie } from '../../shared/model/patients_categories'

interface ServiceCardProps {
  service: PatientsCategorie
  onOpen: (id: number) => void
}

export const ServiceCard: FC<ServiceCardProps> = ({ service, onOpen }) => (
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
  </Card>
)
