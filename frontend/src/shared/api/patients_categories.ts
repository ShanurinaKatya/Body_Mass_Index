import { getToken } from '../lib/storage/patients_categories'
import { Api } from './Api'

// Относительный baseURL: запросы уходят на тот же origin и попадают
// в обратный прокси Vite (/api -> http://localhost:8080).
export const api = new Api({ baseURL: '/api' })

// JWT кладём в заголовок Authorization для всех методов,
// включая помеченные в Swagger как @Security BearerAuth.
api.instance.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})
