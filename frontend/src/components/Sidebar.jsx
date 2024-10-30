import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  Drawer,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Box,
  useTheme,
  useMediaQuery,
  MenuItem,
  Select,
  FormControl,
  InputLabel,
  Divider,
  Typography
} from '@mui/material';
import {
  Home,
  Code,
  AttachMoney,
  Language,
  Balance,
  WarningAmber,
  QueryStats,
  CrisisAlert,
  Category,
  Info,
  SportsBasketball,
  Dangerous,
  Timeline,
  Add
} from '@mui/icons-material';

// Constants
const DRAWER_WIDTH = 240;

const Sidebar = () => {
  const theme = useTheme();
  const navigate = useNavigate();
  const location = useLocation();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const [mobileValue, setMobileValue] = useState(location.pathname);

  const menuItems = [
    { text: 'All questions', path: '/questions', icon: <Category /> },
    { text: 'AI', path: '/questions/category/ai', icon: <Code /> },
    { text: 'Politics', path: '/questions/category/politics', icon: <Balance /> },
    { text: 'Economy', path: '/questions/category/economy', icon: <Timeline /> },
    { text: 'Sweden', path: '/questions/category/sweden', icon: <Category /> },
    { text: 'Finance', path: '/questions/category/finance', icon: <AttachMoney /> },
    { text: 'World', path: '/questions/category/world', icon: <Language /> },
    { text: 'Conflict', path: '/questions/category/conflict', icon: <WarningAmber /> },
    { text: 'Nuclear', path: '/questions/category/nuclear', icon: <CrisisAlert /> },
    { text: 'X-risk', path: '/questions/category/x-risk', icon: <Dangerous /> },
    { text: 'Sports', path: '/questions/category/sports', icon: <SportsBasketball /> },
    { text: 'Personal', path: '/questions/category/personal', icon: <Category /> },
    { text: 'Other', path: '/questions/category/other', icon: <Category /> },
    { text: 'Resolved', path: '/questions/resolved', icon: <Timeline /> } 
  ];


  const handleNavigation = (path) => {
    navigate(path);
    if (isMobile) {
      setMobileValue(path);
    }
  };

  const handleMobileChange = (event) => {
    const path = event.target.value;
    setMobileValue(path);
    navigate(path);
  };

  if (isMobile) {
    return (
      <Box
        sx={{
          position: 'absolute',
          top: 64, // AppBar height
          left: 0,
          right: 0,
          zIndex: 1100,
          backgroundColor: 'background.paper',
          px: 2,
          py: 1,
          borderBottom: 1,
          borderColor: 'divider'
        }}
      >
        <FormControl fullWidth size="small">
          <Select
            value={mobileValue}
            onChange={handleMobileChange}
            displayEmpty
            sx={{ 
              '& .MuiSelect-select': { 
                display: 'flex', 
                alignItems: 'center',
                gap: 1
              }
            }}
          >
            {menuItems.map((item) => (
              <MenuItem key={item.path} value={item.path}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  {item.icon}
                  <Typography>{item.text}</Typography>
                </Box>
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>
    );
  }

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: DRAWER_WIDTH,
        flexShrink: 0,
        '& .MuiDrawer-paper': {
          width: DRAWER_WIDTH,
          boxSizing: 'border-box',
          mt: '64px',
          height: `calc(100% - 64px)`,
          borderRight: `1px solid ${theme.palette.divider}`
        },
      }}
    >
      <List>
        {menuItems.map((item) => (
          <ListItem
            button
            key={item.path}
            onClick={() => handleNavigation(item.path)}
            selected={location.pathname === item.path}
            sx={{
              '&.Mui-selected': {
                backgroundColor: theme.palette.primary.light,
                color: theme.palette.primary.contrastText,
                '&:hover': {
                  backgroundColor: theme.palette.primary.main,
                },
                '& .MuiListItemIcon-root': {
                  color: theme.palette.primary.contrastText,
                }
              }
            }}
          >
            <ListItemIcon sx={{ color: location.pathname === item.path ? 'inherit' : theme.palette.text.primary }}>
              {item.icon}
            </ListItemIcon>
            <ListItemText primary={item.text} />
          </ListItem>
        ))}
      </List>
    </Drawer>
  );
};

export default Sidebar;
