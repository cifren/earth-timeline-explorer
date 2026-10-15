import { createBrowserRouter } from 'react-router'

import { AppLayout } from './app-layout'

export const router = createBrowserRouter([
  {
    Component: AppLayout,
    children: [
      {
        path: '/',
        lazy: async () => ({ Component: (await import('./explorer-page')).ExplorerPage }),
      },
    ],
  },
])
