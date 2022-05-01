import { useEffect } from 'react'
import NextLink from 'next/link'
import { useRouter } from 'next/router'
import { Box, Button, Divider, Drawer, Typography, useMediaQuery } from '@mui/material'
import OpenInNewIcon from '@mui/icons-material/OpenInNew'
import { ChartBar as ChartBarIcon } from '../icons/chart-bar'
import { Cog as CogIcon } from '../icons/cog'
import { Lock as LockIcon } from '../icons/lock'
import { ShoppingBag as ShoppingBagIcon } from '../icons/shopping-bag'
import { User as UserIcon } from '../icons/user'
import { UserAdd as UserAddIcon } from '../icons/user-add'
import { XCircle as XCircleIcon } from '../icons/x-circle'
import Logo from './logo'
import { NavItem } from './nav-item'
import SmartToyRoundedIcon from '@mui/icons-material/SmartToyRounded'

const items = [
  {
    href: '/',
    icon: <ChartBarIcon fontSize="small" />,
    title: 'Dashboard',
  },
  {
    href: '/feeders',
    icon: <SmartToyRoundedIcon fontSize="small" />,
    title: 'Feeders',
  },
  {
    href: '/products',
    icon: <ShoppingBagIcon fontSize="small" />,
    title: 'Products',
  },
  {
    href: '/account',
    icon: <UserIcon fontSize="small" />,
    title: 'Account',
  },
  {
    href: '/settings',
    icon: <CogIcon fontSize="small" />,
    title: 'Settings',
  },
  {
    href: '/login',
    icon: <LockIcon fontSize="small" />,
    title: 'Login',
  },
  {
    href: '/register',
    icon: <UserAddIcon fontSize="small" />,
    title: 'Register',
  },
  {
    href: '/404',
    icon: <XCircleIcon fontSize="small" />,
    title: 'Error',
  },
]

interface DashboardSidebarProps {
  onClose?: () => void
  open: boolean
}

export const DashboardSidebar = (props: DashboardSidebarProps) => {
  const router = useRouter()
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const lgUp = useMediaQuery((theme: any) => theme.breakpoints.up('lg'), {
    defaultMatches: true,
    noSsr: false,
  })

  useEffect(
    () => {
      if (!router.isReady) {
        return
      }

      if (props.open) {
        props.onClose?.()
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [router.asPath],
  )

  const content = (
    <>
      <Box
        sx={{
          display: 'flex',
          flexDirection: 'column',
          height: '100%',
        }}
      >
        <div>
          <Box sx={{ p: 3 }}>
            <NextLink href="/" passHref>
              <a>
                <Logo
                  sx={{
                    height: 42,
                    width: 42,
                  }}
                />
              </a>
            </NextLink>
          </Box>
        </div>
        <Divider
          sx={{
            borderColor: '#2D3748',
            my: 3,
          }}
        />
        <Box sx={{ flexGrow: 1 }}>
          {items.map((item) => (
            <NavItem key={item.title} icon={item.icon} href={item.href} title={item.title} />
          ))}
        </Box>
        <Divider sx={{ borderColor: '#2D3748' }} />
        <Box
          sx={{
            px: 2,
            py: 3,
          }}
        >
          <Typography color="neutral.100" variant="subtitle2">
            Need more features?
          </Typography>
          <Typography color="neutral.500" variant="body2">
            Check out our Pro solution template.
          </Typography>
          <NextLink href="https://material-kit-pro-react.devias.io/" passHref>
            <Button
              color="secondary"
              component="a"
              endIcon={<OpenInNewIcon />}
              fullWidth
              sx={{ mt: 2 }}
              variant="contained"
            >
              Pro Live Preview
            </Button>
          </NextLink>
        </Box>
      </Box>
    </>
  )

  if (lgUp) {
    return (
      <Drawer
        anchor="left"
        open
        PaperProps={{
          sx: {
            backgroundColor: 'neutral.900',
            color: '#FFFFFF',
            width: 280,
          },
        }}
        variant="permanent"
      >
        {content}
      </Drawer>
    )
  }

  return (
    <Drawer
      anchor="left"
      onClose={props.onClose}
      open={props.open}
      PaperProps={{
        sx: {
          backgroundColor: 'neutral.900',
          color: '#FFFFFF',
          width: 280,
        },
      }}
      sx={{ zIndex: (theme) => theme.zIndex.appBar + 100 }}
      variant="temporary"
    >
      {content}
    </Drawer>
  )
}
