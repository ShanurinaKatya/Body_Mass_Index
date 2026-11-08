import { readFile } from 'node:fs/promises'
import { resolve } from 'node:path'

import { generateApi } from 'swagger-typescript-api'

const SWAGGER_URL =
  process.env.SWAGGER_URL ?? 'http://localhost:8080/swagger/doc.json'
const SWAGGER_FILE = resolve(process.cwd(), '../backend/docs/swagger.json')

async function loadSpec() {
  try {
    const response = await fetch(SWAGGER_URL)
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}`)
    }
    console.log(`Swagger загружен из сервиса: ${SWAGGER_URL}`)
    return await response.json()
  } catch (error) {
    console.warn(`Сервис недоступен (${error.message}), беру файл: ${SWAGGER_FILE}`)
    return JSON.parse(await readFile(SWAGGER_FILE, 'utf8'))
  }
}

const spec = await loadSpec()

await generateApi({
  name: 'Api.ts',
  output: resolve(process.cwd(), './src/shared/api'),
  spec,
  httpClientType: 'axios',
  enumStyle: 'union',
  typeOnlyImports: true,
})
