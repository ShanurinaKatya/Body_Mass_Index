import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, Navigate, Outlet, RouterProvider } from 'react-router-dom'

import 'bootstrap/dist/css/bootstrap.min.css'
import './styles/patients_categories.css'

import { TabBar } from '../widgets/tab-bar/patients_categories'
import { FeedPage } from '../pages/feed/patients_categories'
import { TilePage } from '../pages/cards/patients_categories'

const AppLayout = () => (
  <>
    <Outlet />
    <TabBar />
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
    ],
  },
])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
)
