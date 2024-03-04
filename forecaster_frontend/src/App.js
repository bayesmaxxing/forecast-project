import React from 'react';
import { Route, Routes, BrowserRouter} from 'react-router-dom';
import ForecastPage from './pages/ForecastPage';
import FilteredForecastPage from './pages/FilteredForecastPage';
import ResolvedForecastPage from './pages/ResolvedForecastPage';
import SpecificForecast from './pages/SpecificForecast';
import HomePage from './pages/HomePage';
import FaqPage from './pages/faqPage';
import Header from './components/header';
import BlogPost from './pages/BlogPage';
import './App.css';

function App() {
  return (
    <BrowserRouter>
      <div className="App">
        <Header />
        <Routes>  
          <Route path="/home" element={<HomePage />}/>
          <Route path="/" element={<ForecastPage />}/>
          <Route path="/forecast/:id" element={<SpecificForecast />}/>
          <Route path="/category/:category" element={<FilteredForecastPage />}/>
          <Route path="/resolved" element={<ResolvedForecastPage />}/>
          <Route path="/faq" element={<FaqPage />}/>
        </Routes>
      </div>
    </BrowserRouter>
  );
}

export default App;
