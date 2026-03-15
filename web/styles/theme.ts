import { createTheme } from '@mui/material/styles';

const theme = createTheme({
  palette: {
    primary: {
      main: '#1a73e8', // Google Blue
      light: '#e8f0fe',
      dark: '#174ea6',
    },
    secondary: {
      main: '#34a853', // Google Green
    },
    error: {
      main: '#d93025', // Google Red
    },
    warning: {
      main: '#f9ab00', // Google Yellow
    },
    background: {
      default: '#f8f9fa',
      paper: '#ffffff',
    },
    text: {
      primary: '#202124',
      secondary: '#5f6368',
    },
    divider: '#dadce0',
  },
  typography: {
    fontFamily: '"Inter", "Roboto", "Helvetica", "Arial", sans-serif',
    h4: {
      fontWeight: 600,
      color: '#202124',
      letterSpacing: '-0.02em',
    },
    h5: {
      fontWeight: 600,
      color: '#202124',
    },
    h6: {
      fontWeight: 600,
      color: '#202124',
    },
    subtitle1: {
      fontWeight: 500,
      color: '#5f6368',
    },
    body1: {
      color: '#3c4043',
    },
    button: {
      textTransform: 'none',
      fontWeight: 500,
    },
  },
  shape: {
    borderRadius: 8,
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          padding: '8px 24px',
          boxShadow: 'none',
          '&:hover': {
            boxShadow: '0 1px 2px 0 rgba(60,64,67,0.302), 0 1px 3px 1px rgba(60,64,67,0.149)',
          },
        },
        containedPrimary: {
          backgroundColor: '#1a73e8',
          '&:hover': {
            backgroundColor: '#174ea6',
          },
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          boxShadow: '0 1px 2px 0 rgba(60,64,67,0.302), 0 1px 3px 1px rgba(60,64,67,0.149)',
          border: '1px solid #dadce0',
        },
        elevation1: {
          boxShadow: '0 1px 2px 0 rgba(60,64,67,0.302), 0 1px 3px 1px rgba(60,64,67,0.149)',
        },
      },
    },
    MuiTab: {
      styleOverrides: {
        root: {
          fontWeight: 500,
          '&.Mui-selected': {
            color: '#1a73e8',
          },
        },
      },
    },
    MuiTableCell: {
      styleOverrides: {
        head: {
          fontWeight: 600,
          color: '#5f6368',
          backgroundColor: '#f8f9fa',
        },
        root: {
          padding: '12px 16px',
        },
      },
    },
  },
});

export default theme;
