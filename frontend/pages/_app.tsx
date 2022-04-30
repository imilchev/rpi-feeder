import '../styles/globals.css'
import Head from 'next/head'
import type { ReactElement, ReactNode } from 'react'
import { CacheProvider } from '@emotion/react'
import AdapterDateFns from '@mui/lab/AdapterDateFns'
import LocalizationProvider from '@mui/lab/LocalizationProvider'
import { ThemeProvider } from '@mui/material/styles'
import type { AppProps } from 'next/app'
import { theme } from '../theme'
import { createEmotionCache } from '../utils/create-emotion-cache'
import { CssBaseline } from '@mui/material'
import { EmotionCache } from '@emotion/cache'
import type { NextPage } from 'next'
const clientSideEmotionCache = createEmotionCache()

type NextPageWithLayout = NextPage & {
  getLayout?: (page: ReactElement) => ReactNode
}
interface AppPropsWithLayout extends AppProps {
  Component: NextPageWithLayout
  emotionCache?: EmotionCache
}

const App = ({ Component, pageProps }: AppPropsWithLayout) => {
  pageProps.emotionCache = clientSideEmotionCache
  //const { Component, emotionCache = clientSideEmotionCache, pageProps } = props;

  const getLayout = Component.getLayout ?? ((page) => page)

  return (
    <CacheProvider value={pageProps.emotionCache}>
      <Head>
        <title>Material Kit Pro</title>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </Head>
      <LocalizationProvider dateAdapter={AdapterDateFns}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          {getLayout(<Component {...pageProps.pageProps} />)}
        </ThemeProvider>
      </LocalizationProvider>
    </CacheProvider>
  )
}

export default App
