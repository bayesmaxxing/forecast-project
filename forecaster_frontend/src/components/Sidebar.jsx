import React from 'react';
import { Link } from 'react-router-dom';
import './Sidebar.css';

function Sidebar() {
    return (
        <div className='sidebarMenu'>
            <h2>Categories</h2>
            <ul>
                <li><Link to='/'>All Questions</Link></li>
                <li><Link to='/category/ai'>AI</Link></li>
                <li><Link to='/category/sweden'>Sweden</Link></li>
                <li><Link to='/category/economy'>Economy</Link></li>
                <li><Link to='/category/finance'>Finance</Link></li>
                <li><Link to='/category/politics'>Politics</Link></li>
                <li><Link to='/category/world'>World</Link></li>
                <li><Link to='/category/conflict'>Conflict</Link></li>
                <li><Link to='/category/nuclear'>Nuclear</Link></li>
                <li><Link to='/category/x-risk'>X-Risk</Link></li>
                <li><Link to='/category/sports'>Sports</Link></li>
                <li><Link to='/category/personal'>Personal</Link></li>
                <li><Link to='/category/other'>Other</Link></li>
                <li><Link to='/resolved'>Resolved</Link></li>
            </ul>
        </div>
    );
};
export default Sidebar;