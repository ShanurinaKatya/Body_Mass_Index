import type { FC } from 'react'
import { Button } from 'react-bootstrap'

export const MAX_AGE = 100

interface AgeFilterProps {
  age: number
  onChange: (age: number) => void
  onApply: () => void
}

export const AgeFilter: FC<AgeFilterProps> = ({ age, onChange, onApply }) => (
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
          onChange={(event) => onChange(Number(event.target.value))}
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
    <Button className="cards-filter-btn" onClick={onApply}>
      Найти
    </Button>
  </div>
)
