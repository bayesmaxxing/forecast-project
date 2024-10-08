import React from 'react';
import { Route, Routes, BrowserRouter} from 'react-router-dom';
import ForecastPage from './pages/ForecastPage';
import FilteredForecastPage from './pages/FilteredForecastPage';
import ResolvedForecastPage from './pages/ResolvedForecastPage';
import SpecificForecast from './pages/SpecificForecast';
import HomePage from './pages/HomePage';
import FaqPage from './pages/faqPage';
import BlogPage from './pages/BlogPage';
import BlogpostPage from './pages/BlogpostPage';
import Header from './components/header';
import AdminPage from './pages/AdminPage';
import './App.css';

function App() {
  return (
    <BrowserRouter>
      <div className="App">
        <Header />
        <Routes>  
          <Route path="/" element={<HomePage />}/>
          <Route path="/questions" element={<ForecastPage />}/>
          <Route path="/forecast/:id" element={<SpecificForecast />}/>
          <Route path="/questions/category/:category" element={<FilteredForecastPage />}/>
          <Route path="/questions/resolved" element={<ResolvedForecastPage />}/>
          <Route path="/faq" element={<FaqPage />}/>
          <Route path='/blog' element={<BlogPage />}/>
          <Route path='/blog/:slug' element={<BlogpostPage />}/>
          <Route path='/admin' element={<AdminPage />}/>
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
