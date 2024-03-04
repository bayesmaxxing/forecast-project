import React from 'react';
import { Link } from 'react-router-dom';
import './Header.css'; 


function Header () {
    return (
        <header className='header'>
            <h1 className='header-title'>Samuel's Forecasts</h1>
            <nav className='navbar'>
                <ul className='nav-links'>
                    <li><Link to='/home'>Home</Link></li>
                    <li><Link to='/'>Questions</Link></li>
                    <li><Link to='/faq'>FAQ</Link></li>
                </ul>
            </nav>
        </header>
    )
}
export default Header;