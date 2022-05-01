import Feeder, { FeederStatus } from '../dto/feeder'

export const feeders: Feeder[] = [
  {
    clientId: 'rpi-feeder',
    friendlyName: 'Oreo feeder',
    softwareVersion: 'dev',
    status: FeederStatus.Online,
  },
  {
    clientId: 'rpi-feeder2',
    friendlyName: 'Snowball feeder',
    softwareVersion: 'dev',
    status: FeederStatus.Offline,
    lastOnline: 1651409787,
  },
]
