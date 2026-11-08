const TOKEN_KEY = 'patients_categories.jwt'

export const getToken = (): string | null => localStorage.getItem(TOKEN_KEY)

export const setToken = (token: string): void => {
  localStorage.setItem(TOKEN_KEY, token)
}

export const removeToken = (): void => {
  localStorage.removeItem(TOKEN_KEY)
}
