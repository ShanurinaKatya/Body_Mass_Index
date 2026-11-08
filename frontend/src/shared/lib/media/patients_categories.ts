const MINIO_PREFIX = 'http://localhost:9000/body-mass-index'

// Адрес объекта в MinIO заменяем на второй маршрут обратного прокси Vite,
// чтобы браузер не обращался к другому домену напрямую (CORS).
export const toMediaUrl = (url: string | null | undefined): string => {
  if (!url) {
    return ''
  }
  return url.replace(MINIO_PREFIX, '/s3')
}
