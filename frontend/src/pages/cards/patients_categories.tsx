import type { FC } from 'react'
import { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useNavigate } from 'react-router-dom'

import type { AppDispatch, RootState } from '../../app/store/patients_categories'
import { ServiceCard } from '../../entities/service/patients_categories'
import {
  deleteService,
  fetchTiles,
} from '../../entities/service/model/patients_categories'
import { AgeFilter } from '../../features/filter/patients_categories'
import { setAge } from '../../features/filter/model/patients_categories'
import { ROUTES } from '../../shared/config/patients_categories'

export const TilePage: FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()

  const age = useSelector((state: RootState) => state.filter.age)
  const { tiles, error } = useSelector(
    (state: RootState) => state.services,
  )
  const isAuthenticated = useSelector(
    (state: RootState) => state.user.isAuthenticated,
  )

  useEffect(() => {
    dispatch(fetchTiles())
  }, [dispatch])

  const handleFilter = () => {
    dispatch(fetchTiles())
  }

  const handleOpen = (id: number) => navigate(`${ROUTES.FEED}?id=${id}`)

  const handleDelete = (id: number) => {
    dispatch(deleteService(id))
  }

  return (
    <>
      <header className="page-header">
        <h1 className="page-title">Каталог расчётов</h1>
        <p className="page-subtitle">Специализированные калькуляторы ИМТ</p>

        <AgeFilter
          age={age}
          onChange={(value) => dispatch(setAge(value))}
          onApply={handleFilter}
        />
      </header>

      <div className="cards-container">
        {error ? <p className="page-error">{error}</p> : null}

        <div className="cards-grid">
          {tiles.map((service) => (
            <ServiceCard
              key={service.id}
              service={service}
              onOpen={handleOpen}
              onDelete={isAuthenticated ? handleDelete : undefined}
            />
          ))}
        </div>
      </div>
    </>
  )
}
