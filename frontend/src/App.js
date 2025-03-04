import React from 'react';
import { ThemeProvider, CssBaseline, Box } from '@mui/material';
import theme from './theme';
import { Route, Routes, BrowserRouter} from 'react-router-dom';
import ForecastPage from './pages/ForecastPage';
import SpecificForecast from './pages/SpecificForecast';
import HomePage from './pages/HomePage';
import FaqPage from './pages/faqPage';
import BlogPage from './pages/BlogPage';
import BlogpostPage from './pages/BlogpostPage';
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
          <div className="App">
            <Header />
              <Routes>  
                <Route path="/" element={<HomePage />}/>
                <Route path="/questions" element={<ForecastPage />}/>
                <Route path="/forecast/:id" element={<SpecificForecast />}/>
                <Route path="/questions/category/:category" element={<ForecastPage />}/>
                <Route path="/questions/resolved" element={<ForecastPage />}/>
                <Route path="/faq" element={<FaqPage />}/>
                <Route path='/blog' element={<BlogPage />}/>
                <Route path='/blog/:slug' element={<BlogpostPage />}/>
                <Route path='/admin' element={<AdminPage />}/>
              </Routes>
          </div>
        </BrowserRouter>
      </Box>
    </ThemeProvider>
  );
}

export default App;
