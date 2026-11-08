import { createAsyncThunk, createSlice } from '@reduxjs/toolkit'

import type { RootState } from '../../../app/store/patients_categories'
import { api } from '../../../shared/api/patients_categories'
import type {
  DsPatientsCategoriesFeedSerializer,
  DsPatientsCategoriesSerializer,
} from '../../../shared/api/Api'
import { MAX_AGE } from '../../../shared/config/patients_categories'
import type { PatientsCategorie } from '../../../shared/model/patients_categories'
import {
  getMockFeed,
  getMockTiles,
} from '../../../shared/model/patients_categories'

export interface ServicesState {
  tiles: PatientsCategorie[]
  feed: PatientsCategorie | null
  draft: PatientsCategorie | null
  loading: boolean
  error: string | null
  usingMock: boolean
}

const initialState: ServicesState = {
  tiles: [],
  feed: null,
  draft: null,
  loading: false,
  error: null,
  usingMock: false,
}

type RawService = DsPatientsCategoriesSerializer &
  Partial<DsPatientsCategoriesFeedSerializer> &
  Partial<{ is_creator: number }>

const normalize = (raw: RawService | undefined): PatientsCategorie => ({
  id: raw?.id ?? 0,
  name: raw?.name ?? '',
  description: raw?.description ?? '',
  image_url: raw?.image_url ?? '',
  video_url: raw?.video_url ?? '',
  age: raw?.age ?? 0,
  weight: raw?.weight ?? 0,
  likes: raw?.likes ?? 0,
  creator_id: raw?.creator_id ?? 0,
  created_at: raw?.created_at ?? '',
  updated_at: raw?.updated_at ?? '',
  is_creator: raw?.is_creator ?? 0,
  is_liked: raw?.is_liked ?? 0,
})

// Плитка: фильтр по возрасту хранится в Redux и уходит в запрос к бэкенду.
// Гость при недоступном сервисе получает mock-коллекцию.
export const fetchTiles = createAsyncThunk<
  { items: PatientsCategorie[]; usingMock: boolean },
  void,
  { state: RootState }
>('services/fetchTiles', async (_, { getState, rejectWithValue }) => {
  const { filter, user } = getState()
  const query = filter.age < MAX_AGE ? { age: filter.age } : undefined

  try {
    const response =
      await api.patientsCategories.patientsCategoriesList(query)
    return {
      items: (response.data ?? []).map((item) => normalize(item)),
      usingMock: false,
    }
  } catch {
    if (!user.isAuthenticated) {
      return { items: getMockTiles(filter.age), usingMock: true }
    }
    return rejectWithValue('Не удалось загрузить плитку услуг')
  }
})

export const fetchFeed = createAsyncThunk<
  { item: PatientsCategorie; usingMock: boolean },
  { id?: number; next?: boolean },
  { state: RootState }
>('services/fetchFeed', async (args, { getState, rejectWithValue }) => {
  try {
    const response = await api.patientsCategories.feedList(args)
    return { item: normalize(response.data), usingMock: false }
  } catch {
    if (!getState().user.isAuthenticated) {
      return { item: getMockFeed(args.id, args.next), usingMock: true }
    }
    return rejectWithValue('Не удалось загрузить ленту')
  }
})

export const fetchDraft = createAsyncThunk<
  PatientsCategorie | null,
  void,
  { state: RootState }
>('services/fetchDraft', async (_, { rejectWithValue }) => {
  try {
    const response = await api.patientsCategories.draftList()
    const draft = normalize(response.data)
    return draft.id === 0 ? null : draft
  } catch {
    return rejectWithValue('Не удалось получить черновик')
  }
})

export const createService = createAsyncThunk<
  PatientsCategorie,
  { name: string; image?: File; video?: File }
>('services/createService', async (data, { rejectWithValue }) => {
  try {
    // Не передаём незаполненные файлы: сгенерированный FormData
    // иначе отправит строку "undefined" вместо файла.
    const payload: { name: string; image?: File; video?: File } = {
      name: data.name,
    }
    if (data.image) {
      payload.image = data.image
    }
    if (data.video) {
      payload.video = data.video
    }
    const response =
      await api.patientsCategories.patientsCategoriesCreate(payload)
    return normalize(response.data)
  } catch {
    return rejectWithValue('Не удалось создать услугу')
  }
})

export const publishService = createAsyncThunk<
  PatientsCategorie,
  { id: number; description: string; age: number; weight: number }
>('services/publishService', async (data, { rejectWithValue }) => {
  const { id, ...body } = data
  try {
    const response = await api.patientsCategories.publishUpdate(id, body)
    return normalize(response.data)
  } catch {
    return rejectWithValue('Не удалось опубликовать услугу')
  }
})

export const deleteService = createAsyncThunk<
  number,
  number,
  { state: RootState }
>('services/deleteService', async (id, { rejectWithValue }) => {
  try {
    await api.patientsCategories.patientsCategoriesDelete(id)
    return id
  } catch {
    return rejectWithValue('Не удалось удалить услугу')
  }
})

export const likeService = createAsyncThunk<
  { id: number; likes: number; liked: boolean },
  { id: number; like: number }
>('services/likeService', async ({ id, like }, { rejectWithValue }) => {
  try {
    const response = await api.patientsCategories.likeCreate(id, { like })
    return {
      id,
      likes: response.data.likes ?? 0,
      liked: like === 1,
    }
  } catch {
    return rejectWithValue('Не удалось изменить лайк')
  }
})

const servicesSlice = createSlice({
  name: 'services',
  initialState,
  reducers: {
    clearServicesError(state) {
      state.error = null
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchTiles.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(fetchTiles.fulfilled, (state, action) => {
        state.loading = false
        state.tiles = action.payload.items
        state.usingMock = action.payload.usingMock
      })
      .addCase(fetchTiles.rejected, (state, action) => {
        state.loading = false
        state.tiles = []
        state.error = action.payload as string
      })
      .addCase(fetchFeed.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(fetchFeed.fulfilled, (state, action) => {
        state.loading = false
        state.feed = action.payload.item
        state.usingMock = action.payload.usingMock
      })
      .addCase(fetchFeed.rejected, (state, action) => {
        state.loading = false
        state.feed = null
        state.error = action.payload as string
      })
      .addCase(fetchDraft.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(fetchDraft.fulfilled, (state, action) => {
        state.loading = false
        state.draft = action.payload
      })
      .addCase(fetchDraft.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string
      })
      .addCase(createService.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(createService.fulfilled, (state, action) => {
        state.loading = false
        state.draft = action.payload
      })
      .addCase(createService.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string
      })
      .addCase(publishService.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(publishService.fulfilled, (state, action) => {
        state.loading = false
        state.draft = null
        state.feed = action.payload
      })
      .addCase(publishService.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string
      })
      .addCase(deleteService.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(deleteService.fulfilled, (state, action) => {
        state.loading = false
        state.tiles = state.tiles.filter(
          (service) => service.id !== action.payload,
        )
      })
      .addCase(deleteService.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string
      })
      .addCase(likeService.pending, (state) => {
        state.loading = true
        state.error = null
      })
      .addCase(likeService.fulfilled, (state, action) => {
        state.loading = false
        const { id, likes, liked } = action.payload
        if (state.feed?.id === id) {
          state.feed.likes = likes
          state.feed.is_liked = liked ? 1 : 0
        }
        const tile = state.tiles.find((service) => service.id === id)
        if (tile) {
          tile.likes = likes
          tile.is_liked = liked ? 1 : 0
        }
      })
      .addCase(likeService.rejected, (state, action) => {
        state.loading = false
        state.error = action.payload as string
      })
  },
})

export const { clearServicesError } = servicesSlice.actions
export default servicesSlice.reducer
