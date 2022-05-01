import PerfectScrollbar from 'react-perfect-scrollbar'
import { formatDistanceToNow } from 'date-fns'
import fromUnixTime from 'date-fns/fromUnixTime'
import {
  Chip,
  Box,
  Card,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Typography,
  useTheme,
} from '@mui/material'
import FeederDto, { FeederStatus } from '../../dto/feeder'
import React, { useEffect, useState } from 'react'
import { useRouter } from 'next/router'

class Feeder {
  public readonly clientId: string
  public readonly friendlyName: string
  public readonly softwareVersion: string
  public readonly status: string
  public readonly statusColor: 'success' | 'error'
  public lastOnlineDistance: string

  private readonly lastOnline?: Date

  constructor(f: FeederDto) {
    this.clientId = f.ClientId
    this.friendlyName = f.FriendlyName ? f.FriendlyName : f.ClientId
    this.softwareVersion = f.SoftwareVersion
    this.status = f.Status

    switch (f.Status) {
      case FeederStatus.Online:
        this.statusColor = 'success'
        break
      case FeederStatus.Offline:
        this.statusColor = 'error'
        break
      default:
        throw Error('Unknown status ' + f.Status)
    }

    if (f.LastOnline) {
      this.lastOnline = fromUnixTime(f.LastOnline)
    }
    this.lastOnlineDistance = ''
    this.setLastOnlineDistance()
  }

  public setLastOnlineDistance() {
    if (this.lastOnline) {
      this.lastOnlineDistance = formatDistanceToNow(this.lastOnline, { addSuffix: true })
    }
  }
}

interface Props {
  Feeders: FeederDto[]
}

const FeederListResults = (props: Props) => {
  const theme = useTheme()
  console.log(theme)
  const [feeders, setFeeders] = useState(props.Feeders.map((f) => new Feeder(f)))
  const router = useRouter()
  const tick = () => {
    const newFeeders = feeders

    newFeeders.forEach((f) => f.setLastOnlineDistance())
    setFeeders(newFeeders)
  }

  // Refresh the data every 10secs
  useEffect(() => {
    const ticker = setInterval(() => tick(), 10000)
    return clearInterval(ticker)
  })

  return (
    <Card>
      <PerfectScrollbar>
        <Box sx={{ minWidth: 1050 }}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Software version</TableCell>
                <TableCell>Status</TableCell>
                <TableCell>Last online</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {feeders.map((f) => (
                <TableRow
                  hover
                  onClick={() => router.push('feeders/' + f.clientId)}
                  key={f.clientId}
                >
                  <TableCell>
                    <Box
                      sx={{
                        alignItems: 'center',
                        display: 'flex',
                      }}
                    >
                      <Typography color="textPrimary" variant="body1">
                        {f.friendlyName}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell>{f.softwareVersion}</TableCell>
                  <TableCell>
                    <Chip label={f.status} color={f.statusColor} />
                  </TableCell>
                  <TableCell>{f.lastOnlineDistance}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Box>
      </PerfectScrollbar>
    </Card>
  )
}

export default FeederListResults
