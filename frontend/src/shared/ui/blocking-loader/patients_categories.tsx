import type { FC } from 'react'
import { Spinner } from 'react-bootstrap'
import { useSelector } from 'react-redux'

import type { RootState } from '../../../app/store/patients_categories'

// Анимация на время запросов: перекрывает интерфейс и блокирует кнопки.
export const BlockingLoader: FC = () => {
  const loading = useSelector(
    (state: RootState) => state.services.loading || state.user.loading,
  )

  if (!loading) {
    return null
  }

  return (
    <div className="blocking-loader" role="status" aria-live="polite">
      <Spinner animation="border" variant="light" />
    </div>
  )
}
