import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    background: {
      default: '#0F172A',
      paper: '#1E293B',
    },
    primary: {
      main: '#38BDF8',
      light: '#BAE6FD',
      dark: '#0F172A',
    },
    secondary: {
      main: '#F472B6',
    },
    text: {
      primary: '#E2E8F0',
      secondary: '#CBD5E1',
    },
    action: {
      hover: '#38BDF820',
      selected: '#38BDF830',
    },
    error: {
      main: '#EF4444',
    },
    success: {
      main: '#10B981',
    },
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          backgroundColor: '#003135',
          margin: 0,
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: '#003135',
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          backgroundColor: '#024950',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          '&:hover': {
            backgroundColor: '#0FA4AF20',
          },
        },
      },
    },
  },
});

export default theme;
