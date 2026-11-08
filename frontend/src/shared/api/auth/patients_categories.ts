import axios from 'axios'

import type { HandlerLoginResponse } from '../Api'

// Авторизация и деавторизация — только axios, без кодогенерации и thunk.
// Состояние в Redux передаётся вручную из страницы/панели.
export async function loginRequest(
  username: string,
  password: string,
): Promise<string> {
  const { data } = await axios.post<HandlerLoginResponse>('/api/users/auth', {
    username,
    password,
  })
  if (!data.access_token) {
    throw new Error('Сервис не вернул токен')
  }
  return data.access_token
}

export async function logoutRequest(token: string): Promise<void> {
  await axios.post('/api/users/logout', null, {
    headers: { Authorization: `Bearer ${token}` },
  })
}
