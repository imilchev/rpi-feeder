import type { ReactElement } from 'react'
import Head from 'next/head'
import { Box, Container } from '@mui/material'
import FeederListResults from '../../components/feeders/feeder-list-results'
import DashboardLayout from '../../components/dashboard-layout'
import { feeders } from '../../__mocks__/feeders'
import { useRouter } from 'next/router'
import { GetServerSideProps } from 'next'

const Feeder = () => {
  const router = useRouter()
  const clientId = router.query.clientid
  console.log(clientId)
  return (
    <>
      <Head>
        <title>Feeder {clientId} | RPi Feeder</title>
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
            <FeederListResults Feeders={feeders} />
          </Box>
        </Container>
      </Box>
    </>
  )
}

export const getServerSideProps: GetServerSideProps = async () => {
  return {
    props: {}, // will be passed to the page component as props
  }
}

Feeder.getLayout = (page: ReactElement) => <DashboardLayout>{page}</DashboardLayout>

export default Feeder
