import { createSlice } from '@reduxjs/toolkit'
import type { PayloadAction } from '@reduxjs/toolkit'

import { MAX_AGE } from '../../../shared/config/patients_categories'

export interface FilterState {
  age: number
}

const initialState: FilterState = { age: MAX_AGE }

const filterSlice = createSlice({
  name: 'filter',
  initialState,
  reducers: {
    setAge(state, action: PayloadAction<number>) {
      state.age = action.payload
    },
    resetFilter(state) {
      state.age = MAX_AGE
    },
  },
})

export const { setAge, resetFilter } = filterSlice.actions
export default filterSlice.reducer
