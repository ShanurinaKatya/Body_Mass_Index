import type { FC } from 'react'

interface ToggleDescriptionProps {
  expanded: boolean
  onToggle: () => void
}

export const ToggleDescription: FC<ToggleDescriptionProps> = ({
  expanded,
  onToggle,
}) => (
  <span className="feed-more-btn" onClick={onToggle}>
    {expanded ? 'меньше' : 'ещё'}
  </span>
)
