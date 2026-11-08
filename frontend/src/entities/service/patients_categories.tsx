import type { FC, MouseEvent } from 'react'
import { Button, Card } from 'react-bootstrap'

import { toMediaUrl } from '../../shared/lib/media/patients_categories'
import type { PatientsCategorie } from '../../shared/model/patients_categories'

interface ServiceCardProps {
  service: PatientsCategorie
  onOpen: (id: number) => void
  onDelete?: (id: number) => void
}

export const ServiceCard: FC<ServiceCardProps> = ({
  service,
  onOpen,
  onDelete,
}) => {
  const handleAction = (event: MouseEvent, action: () => void) => {
    event.stopPropagation()
    action()
  }

  return (
    <Card className="card">
      <div className="card-link" onClick={() => onOpen(service.id)}>
        <Card.Img
          className="card-image"
          variant="top"
          src={toMediaUrl(service.image_url) || '/DefaultImage.jpg'}
          alt={service.name}
          onError={(event) => {
            event.currentTarget.src = '/DefaultImage.jpg'
          }}
        />
        <div className="card-content">
          <Card.Title className="card-title">{service.name}</Card.Title>
          <div className="card-row">
            <span className="card-meta">{service.weight} кг</span>
            <span className="card-meta">{service.age} лет</span>
          </div>
          {service.is_creator === 1 && onDelete ? (
            <Button
              type="button"
              className="card-delete"
              onClick={(event) =>
                handleAction(event, () => onDelete(service.id))
              }
            >
              Удалить
            </Button>
          ) : null}
        </div>
      </div>
    </Card>
  )
}
