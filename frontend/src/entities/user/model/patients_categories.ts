import { createSlice } from '@reduxjs/toolkit'
import type { PayloadAction } from '@reduxjs/toolkit'

export interface UserState {
  username: string | null
  isAuthenticated: boolean
  loading: boolean
  error: string | null
}

// Начальное состояние — гость. Токен в localStorage, но Redux в память
// не сохраняется, поэтому F5 возвращает интерфейс гостя.
const initialState: UserState = {
  username: null,
  isAuthenticated: false,
  loading: false,
  error: null,
}

const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    setLoading(state, action: PayloadAction<boolean>) {
      state.loading = action.payload
    },
    setError(state, action: PayloadAction<string | null>) {
      state.error = action.payload
    },
    setUser(state, action: PayloadAction<string>) {
      state.username = action.payload
      state.isAuthenticated = true
      state.error = null
    },
    clearUser() {
      return initialState
    },
  },
})

export const { setLoading, setError, setUser, clearUser } = userSlice.actions
export default userSlice.reducer
