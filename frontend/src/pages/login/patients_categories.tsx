import type { ChangeEvent, FC, FormEvent } from 'react'
import { useState } from 'react'
import { Alert, Button, Form } from 'react-bootstrap'
import { useDispatch, useSelector } from 'react-redux'
import { Navigate, useNavigate } from 'react-router-dom'

import type { AppDispatch, RootState } from '../../app/store/patients_categories'
import {
  setError,
  setLoading,
  setUser,
} from '../../entities/user/model/patients_categories'
import { loginRequest } from '../../shared/api/auth/patients_categories'
import { ROUTES } from '../../shared/config/patients_categories'
import { setToken } from '../../shared/lib/storage/patients_categories'

export const LoginPage: FC = () => {
  const dispatch = useDispatch<AppDispatch>()
  const navigate = useNavigate()

  const isAuthenticated = useSelector(
    (state: RootState) => state.user.isAuthenticated,
  )
  const loading = useSelector((state: RootState) => state.user.loading)
  const error = useSelector((state: RootState) => state.user.error)

  const [form, setForm] = useState({ username: '', password: '' })

  if (isAuthenticated) {
    return <Navigate to={ROUTES.FEED} replace />
  }

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [event.target.name]: event.target.value })
  }

  // Авторизация — вручную через axios, состояние в Redux кладём сами.
  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    if (!form.username || !form.password) {
      dispatch(setError('Заполните логин и пароль'))
      return
    }

    dispatch(setLoading(true))
    dispatch(setError(null))

    try {
      const token = await loginRequest(form.username, form.password)
      setToken(token)
      dispatch(setUser(form.username))
      navigate(ROUTES.FEED)
    } catch {
      dispatch(setError('Неверный логин или пароль'))
    } finally {
      dispatch(setLoading(false))
    }
  }

  return (
    <>
      <header className="page-header">
        <h1 className="page-title">Вход</h1>
      </header>

      <div className="auth-container">
        {error ? <Alert variant="danger">{error}</Alert> : null}

        <Form onSubmit={handleSubmit}>
          <Form.Group className="add-field" controlId="username">
            <Form.Control
              className="add-input"
              type="text"
              name="username"
              value={form.username}
              onChange={handleChange}
              placeholder="Логин"
            />
          </Form.Group>

          <Form.Group className="add-field" controlId="password">
            <Form.Control
              className="add-input"
              type="password"
              name="password"
              value={form.password}
              onChange={handleChange}
              placeholder="Пароль"
            />
          </Form.Group>

          <Button className="auth-submit" type="submit" disabled={loading}>
            Войти
          </Button>
        </Form>
      </div>
    </>
  )
}
