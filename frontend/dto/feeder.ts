export enum FeederStatus {
  Online = 'online',
  Offline = 'offline',
}

interface Feeder {
  ClientId: string
  FriendlyName?: string
  SoftwareVersion: string
  Status: FeederStatus
  LastOnline?: number
}

export default Feeder
