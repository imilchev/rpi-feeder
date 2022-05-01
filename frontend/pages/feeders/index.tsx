import type { ReactElement } from 'react'
import Head from 'next/head'
import { Box, Container } from '@mui/material'
import FeederListResults from '../../components/feeders/feeder-list-results'
import DashboardLayout from '../../components/dashboard-layout'
import Feeder from '../../dto/feeder'
import { GetServerSideProps } from 'next'

interface Props {
  Feeders: Feeder[]
}

const Feeders = (props: Props) => {
  return (
    <>
      <Head>
        <title>Feeders | RPi Feeder</title>
      </Head>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          py: 8,
        }}
      >
        <Container maxWidth={false}>
          <Box sx={{ mt: 3 }}>
            <FeederListResults Feeders={props.Feeders} />
          </Box>
        </Container>
      </Box>
    </>
  )
}

export const getServerSideProps: GetServerSideProps = async () => {
  const feeders = await fetch('http://localhost:1234/v1/feeders')
    .then((response) => {
      if (!response.ok) {
        console.error(response.body)
      }
      return response.json() as Promise<Feeder[]>
    })
    .catch((e) => console.error(e))
  return {
    props: {
      pageProps: {
        Feeders: feeders,
      },
    }, // will be passed to the page component as props
  }
}

Feeders.getLayout = (page: ReactElement) => <DashboardLayout>{page}</DashboardLayout>

export default Feeders
