import React from 'react';
import { Link } from 'react-router-dom';
import './Sidebar.css';

function Sidebar( {onSearchChange} ) {
    return (
        <div className='sidebarMenu'>
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
                onChange={onSearchChange} // Invoke the passed function on input change
                className="searchInput"
            />
        </div>
    );
};
export default Sidebar;