import Head from 'next/head'
import { useState, SyntheticEvent } from 'react'
import { Container, Typography, Box, Tabs, Tab, ThemeProvider, CssBaseline } from '@mui/material'
import theme from '../styles/theme'
import Dashboard from '../components/Dashboard'
import CRADashboard from '../components/CRADashboard'

// TabPanelProps defines the structure for custom tab content panels.
interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

// CustomTabPanel renders children only when the associated tab is selected.
function CustomTabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ pt: 3 }}>
          {children}
        </Box>
      )}
    </div>
  );
}

// a11yProps generates accessibility attributes for tab elements.
function a11yProps(index: number) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`,
  };
}

// Home serves as the main entry point for the Next.js compliance dashboard.
export default function Home() {
  const [value, setValue] = useState(0);

  const handleChange = (event: SyntheticEvent, newValue: number) => {
    setValue(newValue);
  };

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Head>
        <title>CRA Compliance System</title>
        <meta name="description" content="Cyber Resilience Act Compliance Assessment Dashboard" />
      </Head>
      <Box component="main" sx={{ minHeight: '100vh', backgroundColor: 'background.default' }}>
        <Container maxWidth="xl" sx={{ py: 6 }}>
          <Box sx={{ mb: 6 }}>
            <Typography variant="h4" component="h1" gutterBottom sx={{ mb: 1 }}>
              CRA Compliance System
            </Typography>
            <Typography variant="subtitle1">
              Automated Cyber Resilience Act (CRA) Conformity Assessment
            </Typography>
          </Box>

          <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
            <Tabs
              value={value}
              onChange={handleChange}
              aria-label="cra dashboard tabs"
              textColor="primary"
              indicatorColor="primary"
            >
              <Tab label="Compliance Dashboard" {...a11yProps(0)} />
              <Tab label="Live Agent Logs" {...a11yProps(1)} />
            </Tabs>
          </Box>

          <CustomTabPanel value={value} index={0}>
            <CRADashboard />
          </CustomTabPanel>

          <CustomTabPanel value={value} index={1}>
            <Dashboard />
          </CustomTabPanel>

        </Container>
      </Box>
    </ThemeProvider>
  )
}

