import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Provider } from 'react-redux'
import {
  createBrowserRouter,
  Navigate,
  Outlet,
  RouterProvider,
} from 'react-router-dom'

import 'bootstrap/dist/css/bootstrap.min.css'
import './styles/patients_categories.css'

import { AddPage } from '../pages/add/patients_categories'
import { FeedPage } from '../pages/feed/patients_categories'
import { LoginPage } from '../pages/login/patients_categories'
import { TilePage } from '../pages/cards/patients_categories'
import { BlockingLoader } from '../shared/ui/blocking-loader/patients_categories'
import { TabBar } from '../widgets/tab-bar/patients_categories'
import { store } from './store/patients_categories'

const AppLayout = () => (
  <>
    <Outlet />
    <TabBar />
    <BlockingLoader />
  </>
)

const router = createBrowserRouter([
  {
    path: '/',
    element: <AppLayout />,
    children: [
      { index: true, element: <Navigate to="/feed" replace /> },
      { path: 'feed', element: <FeedPage /> },
      { path: 'cards', element: <TilePage /> },
      { path: 'add', element: <AddPage /> },
      { path: 'login', element: <LoginPage /> },
    ],
  },
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Provider store={store}>
      <RouterProvider router={router} />
    </Provider>
  </StrictMode>,
)
