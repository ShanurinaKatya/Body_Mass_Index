import type { FC } from 'react'

interface NextServiceProps {
  onClick: () => void
}

export const NextService: FC<NextServiceProps> = ({ onClick }) => (
  <div
    className="feed-sidebar-item feed-sidebar-item--clickable"
    onClick={onClick}
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
)
