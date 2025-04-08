import React from 'react';
import { ThemeProvider, CssBaseline, Box } from '@mui/material';
import theme from './theme';
import { Route, Routes, BrowserRouter} from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Login from './components/Login';
import Register from './components/Register';
import ProtectedRoute from './components/ProtectedRoute';
import ForecastPage from './pages/ForecastPage';
import SpecificForecast from './pages/SpecificForecast';
import HomePage from './pages/HomePage';
import FaqPage from './pages/faqPage';
import Header from './components/Header';
import AdminPage from './pages/AdminPage';

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{
        bgcolor: 'background.default',
        minHeight: '100vh'
      }}>
        <BrowserRouter>
          <AuthProvider>
            <div className="App">
              <Header />
              <Routes>  
                <Route path="/" element={<HomePage />}/> 
                <Route path="/questions" element={<ForecastPage />}/>
                <Route path="/forecast/:id" element={<SpecificForecast />}/>
                <Route path="/questions/category/:category" element={<ForecastPage />}/>
                <Route path="/questions/resolved" element={<ForecastPage />}/>
                <Route path="/faq" element={<FaqPage />}/>
                <Route path='/admin' element={<AdminPage />}/>
                <Route path='/login' element={<Login />}/>
                <Route path='/register' element={<Register />}/>
              </Routes>
            </div>
          </AuthProvider>
        </BrowserRouter>
      </Box>
    </ThemeProvider>
  );
}

export default App;
