import React, { useState } from 'react';
import { Link, NavLink } from 'react-router-dom';
import './Navbar.css';
import { FaUserCircle } from 'react-icons/fa';

const Navbar = () => {
    const [isOpen, setIsOpen] = useState(false);
    const toggleMenu = () => setIsOpen(!isOpen);

    const username = localStorage.getItem("username");

    return (
        <nav className="navbar">
            <div className="navbar-brand">
                <Link to="/main">Logo</Link>
            </div>

            <button className="navbar-toggle" onClick={toggleMenu}>
                ☰
            </button>

            <ul className={`navbar-menu ${isOpen ? 'open' : ''}`}>
                <li>
                    <NavLink to="/main" className={({ isActive }) => isActive ? "active" : ""}>Главная</NavLink>
                </li>
                <li>
                    <NavLink to="/complaints" className={({ isActive }) => isActive ? "active" : ""}>Жалобы</NavLink>
                </li>
                <li>
                    <NavLink to="/bookings" className={({ isActive }) => isActive ? "active" : ""}>Бронирования</NavLink>
                </li>
                <li>
                    <NavLink to="/rooms" className={({ isActive }) => isActive ? "active" : ""}>Номера</NavLink>
                </li>
                <li>
                    <NavLink to="/payments" className={({ isActive }) => isActive ? "active" : ""}>Платежи</NavLink>
                </li>

                {username && (
                    <li>
                        <NavLink to="/user" className={({ isActive }) => isActive ? "active user-link" : "user-link"}>
                            <FaUserCircle style={{ fontSize: '20px', marginRight: '5px', verticalAlign: 'middle' }} />
                            {username}
                        </NavLink>
                    </li>
                )}

                <li>
                    <NavLink to="/logout" className={({ isActive }) => isActive ? "active" : ""}>Выйти</NavLink>
                </li>
            </ul>
        </nav>
    );
};

export default Navbar;
