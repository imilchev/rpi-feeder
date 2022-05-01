import styled from '@emotion/styled'
import { AppBar, Avatar, Badge, Box, IconButton, Toolbar, Tooltip } from '@mui/material'
import MenuIcon from '@mui/icons-material/Menu'
import SearchIcon from '@mui/icons-material/Search'
import { Bell as BellIcon } from '../icons/bell'
import { UserCircle as UserCircleIcon } from '../icons/user-circle'

const DashboardNavbarRoot = styled(AppBar)(({ theme }) => ({
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  backgroundColor: (theme as any).palette.background.paper,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  boxShadow: (theme as any).shadows[3],
}))

interface DashboardNavbarProps {
  onSidebarOpen: () => void
}

export const DashboardNavbar = (props: DashboardNavbarProps) => {
  //const { onSidebarOpen, ...other } = props

  return (
    <>
      <DashboardNavbarRoot
        sx={{
          left: {
            lg: 280,
          },
          width: {
            lg: 'calc(100% - 280px)',
          },
        }}
        // {...other}
      >
        <Toolbar
          disableGutters
          sx={{
            minHeight: 64,
            left: 0,
            px: 2,
          }}
        >
          <IconButton
            onClick={props.onSidebarOpen}
            sx={{
              display: {
                xs: 'inline-flex',
                lg: 'none',
              },
            }}
          >
            <MenuIcon fontSize="small" />
          </IconButton>
          <Tooltip title="Search">
            <IconButton sx={{ ml: 1 }}>
              <SearchIcon fontSize="small" />
            </IconButton>
          </Tooltip>
          <Box sx={{ flexGrow: 1 }} />
          <Tooltip title="Notifications">
            <IconButton sx={{ ml: 1 }}>
              <Badge badgeContent={4} color="primary" variant="dot">
                <BellIcon fontSize="small" />
              </Badge>
            </IconButton>
          </Tooltip>
          <Avatar
            sx={{
              height: 40,
              width: 40,
              ml: 1,
            }}
            src="/static/images/avatars/avatar_1.png"
          >
            <UserCircleIcon fontSize="small" />
          </Avatar>
        </Toolbar>
      </DashboardNavbarRoot>
    </>
  )
}
