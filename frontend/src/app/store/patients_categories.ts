import { combineReducers, configureStore } from '@reduxjs/toolkit'

import servicesReducer from '../../entities/service/model/patients_categories'
import userReducer from '../../entities/user/model/patients_categories'
import filterReducer from '../../features/filter/model/patients_categories'

const rootReducer = combineReducers({
  services: servicesReducer,
  filter: filterReducer,
  user: userReducer,
})

// Redux DevTools включается автоматически в dev-режиме.
export const store = configureStore({
  reducer: rootReducer,
  devTools: true,
})

export type RootState = ReturnType<typeof rootReducer>
export type AppDispatch = typeof store.dispatch
