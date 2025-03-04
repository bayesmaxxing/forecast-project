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
        backgroundColor: theme.palette.background.default,
        borderBottom: `1px solid ${theme.palette.primary.main}`,
        zIndex: theme.zIndex.drawer + 1,
      }}
    >
      <Toolbar>
        <Typography
          variant="h6"
          component="div"
          sx={{ 
            flexGrow: 1, 
            fontWeight: 'bold',
            color: theme.palette.primary.light, 
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
                  color: theme.palette.primary.light,
                  '&:hover': {
                    backgroundColor: theme.palette.action.hover,
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
