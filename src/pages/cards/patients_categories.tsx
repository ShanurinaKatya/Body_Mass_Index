import type { FC } from 'react'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { ServiceCard } from '../../entities/service/patients_categories'
import { AgeFilter, MAX_AGE } from '../../features/filter/patients_categories'
import { ROUTES } from '../../shared/config/patients_categories'
import { SERVICES_MOCK } from '../../shared/model/patients_categories'
import type { PatientsCategorie } from '../../shared/model/patients_categories'

export const TilePage: FC = () => {
  const navigate = useNavigate()

  const [age, setAge] = useState(MAX_AGE)
  const [appliedAge, setAppliedAge] = useState(MAX_AGE)

  // Как в лабе 1 (handler.CardsHandler): возраст 100 — все карточки,
  // иначе строгое совпадение age = N, совпадений нет — пусто.
  const matchesAge =
    appliedAge >= 0 && appliedAge < MAX_AGE
      ? (service: PatientsCategorie) => service.age === appliedAge
      : () => true

  const services = SERVICES_MOCK.filter((service) => matchesAge(service))

  const handleFilter = () => setAppliedAge(age)

  const handleOpen = (id: number) => navigate(`${ROUTES.FEED}?id=${id}`)

  return (
    <>
      <header className="page-header">
        <h1 className="page-title">Каталог расчётов</h1>
        <p className="page-subtitle">Специализированные калькуляторы ИМТ</p>

        <AgeFilter age={age} onChange={setAge} onApply={handleFilter} />
      </header>

      <div className="cards-container">
        <div className="cards-grid">
          {services.map((service) => (
            <ServiceCard
              key={service.id}
              service={service}
              onOpen={handleOpen}
            />
          ))}
        </div>
      </div>
    </>
  )
}
