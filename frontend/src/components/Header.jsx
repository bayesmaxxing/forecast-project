import React from 'react';
import { Link } from 'react-router-dom';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  IconButton,
  Menu,
  MenuItem,
  useTheme,
  useMediaQuery
} from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';

const Header = () => {
  const [anchorEl, setAnchorEl] = React.useState(null);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const handleMenu = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const navItems = [
    { text: 'Home', path: '/' },
    { text: 'Questions', path: '/questions' },
    { text: 'Blog', path: '/blog' },
    { text: 'FAQ', path: '/faq' }
  ];

  return (
    <AppBar
      position="fixed"
      sx={{
        backgroundColor: 'rgba(10, 10, 10, 0.8)',
        backdropFilter: 'blur(12px)',
        borderBottom: `1px solid ${theme.palette.divider}`,
        zIndex: theme.zIndex.drawer + 1,
        boxShadow: '0 1px 3px 0 rgba(0, 0, 0, 0.3)',
      }}
    >
      <Toolbar>
        <Typography
          variant="h6"
          component="div"
          sx={{
            flexGrow: 1,
            fontWeight: 700,
            background: 'linear-gradient(135deg, #FF6B6B 0%, #FFA07A 100%)',
            backgroundClip: 'text',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            letterSpacing: '-0.02em',
          }}
        >
          Samuel's Forecasts
        </Typography>

        {isMobile ? (
          <>
            <IconButton
              size="large"
              edge="end"
              color="inherit"
              aria-label="menu"
              onClick={handleMenu}
            >
              <MenuIcon />
            </IconButton>
            <Menu
              anchorEl={anchorEl}
              open={Boolean(anchorEl)}
              onClose={handleClose}
              sx={{
                '& .MuiPaper-root': {
                  width: '100%',
                  maxWidth: '100%',
                  left: '0 !important',
                  right: '0',
                  backgroundColor: theme.palette.background.default,
                  color: theme.palette.primary.light,
                },
                '& .MuiMenuItem-root': {
                  justifyContent: 'center',
                  padding: '16px',
                  '&:hover': {
                    backgroundColor: theme.palette.action.hover,
                  },
                },
              }}
              PaperProps={{
                style: {
                  marginTop: '8px',
                }
              }}
            >
              {navItems.map((item) => (
                <MenuItem 
                  key={item.text} 
                  onClick={handleClose}
                  component={Link}
                  to={item.path}
                >
                  {item.text}
                </MenuItem>
              ))}
            </Menu>
          </>
        ) : (
          <Box sx={{ display: 'flex', gap: 2 }}>
            {navItems.map((item) => (
              <Button
                key={item.text}
                component={Link}
                to={item.path}
                sx={{
                  color: 'text.primary',
                  fontWeight: 500,
                  position: 'relative',
                  '&::after': {
                    content: '""',
                    position: 'absolute',
                    bottom: 0,
                    left: '50%',
                    transform: 'translateX(-50%)',
                    width: 0,
                    height: '2px',
                    background: 'linear-gradient(90deg, #FF6B6B, #FFA07A)',
                    transition: 'width 0.3s ease-in-out',
                  },
                  '&:hover': {
                    backgroundColor: 'transparent',
                    color: 'primary.main',
                    '&::after': {
                      width: '80%',
                    },
                  },
                }}
              >
                {item.text}
              </Button>
            ))}
          </Box>
        )}
      </Toolbar>
    </AppBar>
  );
};

export default Header;
