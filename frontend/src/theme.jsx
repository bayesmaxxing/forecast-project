import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    background: {
      default: '#003135',
      paper: '#024950',
    },
    primary: {
      main: '#0FA4AF',
      light: '#AFDDE5',
      dark: '#003135',
    },
    secondary: {
      main: '#964734',
    },
    text: {
      primary: '#AFDDE5',
      secondary: '#AFDDE5',
    },
    action: {
      hover: '#0FA4AF20',
      selected: '#0FA4AF30',
    },
    error: {
      main: '#964734',
    },
    success: {
      main: '#0FA4AF',
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
