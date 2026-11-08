import type {
  ChangeEvent,
  FC,
  FormEvent,
  KeyboardEvent,
} from 'react'
import { useEffect, useRef, useState } from 'react'
import { Button, Form } from 'react-bootstrap'
import { useDispatch, useSelector } from 'react-redux'
import { Navigate, useNavigate } from 'react-router-dom'

import type { AppDispatch, RootState } from '../../app/store/patients_categories'
import {
  createService,
  fetchDraft,
  publishService,
} from '../../entities/service/model/patients_categories'
import { ROUTES } from '../../shared/config/patients_categories'
import { toMediaUrl } from '../../shared/lib/media/patients_categories'
import type { PatientsCategorie } from '../../shared/model/patients_categories'

const DEFAULT_IMAGE = '/DefaultImage.jpg'
const DEFAULT_VIDEO = '/Default.mp4'

interface MediaPlaceholderProps {
  preview: string
  label: string
  isVideo?: boolean
  onClick?: () => void
}

const MediaPlaceholder: FC<MediaPlaceholderProps> = ({
  preview,
  label,
  isVideo,
  onClick,
}) => {
  const handleKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    if (!onClick) {
      return
    }
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault()
      onClick()
    }
  }

  return (
    <div
      className={`add-placeholder has-preview${
        onClick ? '' : ' add-placeholder--static'
      }`}
      onClick={onClick}
      onKeyDown={handleKeyDown}
      role={onClick ? 'button' : undefined}
      tabIndex={onClick ? 0 : undefined}
    >
      {isVideo ? (
        <video
          className="add-preview"
          src={preview}
          muted
          playsInline
          preload="metadata"
          onError={(event) => {
            if (!event.currentTarget.src.endsWith(DEFAULT_VIDEO)) {
              event.currentTarget.src = DEFAULT_VIDEO
            }
          }}
        />
      ) : (
        <img
          className="add-preview"
          src={preview}
          alt={label}
          onError={(event) => {
            if (!event.currentTarget.src.endsWith(DEFAULT_IMAGE)) {
              event.currentTarget.src = DEFAULT_IMAGE
            }
          }}
        />
      )}
      <span className="add-placeholder-text">{label}</span>
    </div>
  )
}

interface PublishFormProps {
  draft: PatientsCategorie
  onPublished: () => void
}

const PublishForm: FC<PublishFormProps> = ({ draft, onPublished }) => {
  const dispatch = useDispatch<AppDispatch>()
  const loading = useSelector((state: RootState) => state.services.loading)

  const [description, setDescription] = useState(draft.description ?? '')
  const [age, setAge] = useState(draft.age ? String(draft.age) : '')
  const [weight, setWeight] = useState(draft.weight ? String(draft.weight) : '')

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    const result = await dispatch(
      publishService({
        id: draft.id,
        description,
        age: Number(age),
        weight: Number(weight),
      }),
    )
    if (publishService.fulfilled.match(result)) {
      onPublished()
    }
  }

  return (
    <Form onSubmit={handleSubmit}>
      <div className="add-placeholders">
        <MediaPlaceholder
          preview={toMediaUrl(draft.image_url) || DEFAULT_IMAGE}
          label="Обложка"
        />
        <MediaPlaceholder
          preview={toMediaUrl(draft.video_url) || DEFAULT_VIDEO}
          label="Видео"
          isVideo
        />
      </div>

      <div className="add-field">
        <span className="add-field-label">Название расчёта</span>
        <div className="add-name-readonly">{draft.name}</div>
      </div>

      <Form.Group className="add-field" controlId="description">
        <Form.Label className="add-field-label">Описание</Form.Label>
        <Form.Control
          className="add-textarea"
          as="textarea"
          rows={4}
          value={description}
          onChange={(event: ChangeEvent<HTMLTextAreaElement>) =>
            setDescription(event.target.value)
          }
        />
      </Form.Group>

      <div className="add-row">
        <Form.Group className="add-field" controlId="age">
          <Form.Label className="add-field-label">Возраст пациента</Form.Label>
          <Form.Control
            className="add-input"
            type="number"
            min={1}
            value={age}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              setAge(event.target.value)
            }
          />
        </Form.Group>

        <Form.Group className="add-field" controlId="weight">
          <Form.Label className="add-field-label">Вес (кг)</Form.Label>
          <Form.Control
            className="add-input"
            type="number"
            min={0}
            step="0.01"
            value={weight}
            onChange={(event: ChangeEvent<HTMLInputElement>) =>
              setWeight(event.target.value)
            }
          />
        </Form.Group>
      </div>

      <div className="add-actions">
        <Button className="add-submit-btn" type="submit" disabled={loading}>
          Опубликовать
        </Button>
      </div>
    </Form>
  )
}

interface CreateFormProps {
  loading: boolean
  onCreate: (data: { name: string; image?: File; video?: File }) => void
}

const CreateForm: FC<CreateFormProps> = ({ loading, onCreate }) => {
  const [name, setName] = useState('')
  const [image, setImage] = useState<File>()
  const [video, setVideo] = useState<File>()
  const [imagePreview, setImagePreview] = useState('')
  const [videoPreview, setVideoPreview] = useState('')

  const imageInputRef = useRef<HTMLInputElement>(null)
  const videoInputRef = useRef<HTMLInputElement>(null)

  useEffect(
    () => () => {
      if (imagePreview) {
        URL.revokeObjectURL(imagePreview)
      }
    },
    [imagePreview],
  )

  useEffect(
    () => () => {
      if (videoPreview) {
        URL.revokeObjectURL(videoPreview)
      }
    },
    [videoPreview],
  )

  const handleImageChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) {
      return
    }
    setImage(file)
    setImagePreview(URL.createObjectURL(file))
  }

  const handleVideoChange = (event: ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) {
      return
    }
    setVideo(file)
    setVideoPreview(URL.createObjectURL(file))
  }

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault()
    if (!name.trim()) {
      return
    }
    onCreate({ name: name.trim(), image, video })
  }

  return (
    <Form onSubmit={handleSubmit}>
      <div className="add-placeholders">
        <MediaPlaceholder
          preview={imagePreview || DEFAULT_IMAGE}
          label="Обложка"
          onClick={() => imageInputRef.current?.click()}
        />
        <MediaPlaceholder
          preview={videoPreview || DEFAULT_VIDEO}
          label="Видео"
          isVideo
          onClick={() => videoInputRef.current?.click()}
        />
      </div>

      <input
        ref={imageInputRef}
        className="add-file-input"
        type="file"
        accept="image/*"
        hidden
        onChange={handleImageChange}
      />
      <input
        ref={videoInputRef}
        className="add-file-input"
        type="file"
        accept="video/*"
        hidden
        onChange={handleVideoChange}
      />

      <Form.Group className="add-field" controlId="name">
        <Form.Label className="add-field-label">Название расчёта</Form.Label>
        <Form.Control
          className="add-input"
          type="text"
          value={name}
          onChange={(event: ChangeEvent<HTMLInputElement>) =>
            setName(event.target.value)
          }
          placeholder="Например, ИМТ для подростков"
        />
      </Form.Group>

      <div className="add-actions">
        <Button className="add-submit-btn" type="submit" disabled={loading}>
          Далее
        </Button>
      </div>
    </Form>
  )
}

export const AddPage: FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()

  const isAuthenticated = useSelector(
    (state: RootState) => state.user.isAuthenticated,
  )
  const { draft, error, loading } = useSelector(
    (state: RootState) => state.services,
  )

  const hasDraft = Boolean(draft && draft.id !== 0)

  useEffect(() => {
    if (isAuthenticated) {
      dispatch(fetchDraft())
    }
  }, [dispatch, isAuthenticated])

  if (!isAuthenticated) {
    return <Navigate to={ROUTES.LOGIN} replace />
  }

  const handleCreate = (data: { name: string; image?: File; video?: File }) => {
    dispatch(createService(data))
  }

  return (
    <>
      <header className="page-header">
        <h1 className="page-title">Новый расчёт ИМТ</h1>
      </header>

      <div className="add-container">
        {error ? <p className="page-error">{error}</p> : null}

        {hasDraft && draft ? (
          <PublishForm
            key={draft.id}
            draft={draft}
            onPublished={() => navigate(ROUTES.TILE)}
          />
        ) : (
          <CreateForm loading={loading} onCreate={handleCreate} />
        )}
      </div>
    </>
  )
}
