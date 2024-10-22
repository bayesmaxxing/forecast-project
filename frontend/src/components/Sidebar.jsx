import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import './Sidebar.css';

function Sidebar( {onSearchChange} ) {
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    const toggleMenu = () => {
      setIsMenuOpen(!isMenuOpen);
    };

    return (
      <>
        <button
          onClick={toggleMenu}
          className="mobile-menu-button"
          aria-label="Toggle menu"
        >
          <div className={`hamburger ${isMenuOpen ? 'open' : ''}`}>
            <span></span>
            <span></span>
            <span></span>
          </div>
        </button>
        <div className={`sidebarMenu ${isMenuOpen ? 'open' : ''}`}>
          <h2>Categories</h2>
            <ul>
                <li><Link to='/questions'>All Questions</Link></li>
                <li><Link to='/questions/category/ai'>AI</Link></li>
                <li><Link to='/questions/category/sweden'>Sweden</Link></li>
                <li><Link to='/questions/category/economy'>Economy</Link></li>
                <li><Link to='/questions/category/finance'>Finance</Link></li>
                <li><Link to='/questions/category/politics'>Politics</Link></li>
                <li><Link to='/questions/category/world'>World</Link></li>
                <li><Link to='/questions/category/conflict'>Conflict</Link></li>
                <li><Link to='/questions/category/nuclear'>Nuclear</Link></li>
                <li><Link to='/questions/category/x-risk'>X-Risk</Link></li>
                <li><Link to='/questions/category/sports'>Sports</Link></li>
                <li><Link to='/questions/category/personal'>Personal</Link></li>
                <li><Link to='/questions/category/other'>Other</Link></li>
                <li><Link to='/questions/resolved'>Resolved</Link></li>
            </ul> 
            <input
                type="text"
                placeholder="Search forecasts"
                onChange={onSearchChange}
                className="searchInput"
            />
        </div>
        
        {isMenuOpen && (
          <div 
            className="sidebar-overlay"
            onClick={toggleMenu}
          />
        )}
      </>
   );
};
export default Sidebar;
