import React, { useState } from 'react';
import { Link, NavLink } from 'react-router-dom';
import './Navbar.css';

const Navbar = () => {
    const [isOpen, setIsOpen] = useState(false);

    const toggleMenu = () => {
        setIsOpen(!isOpen);
    };

    return (
        <nav className="navbar">
            <div className="navbar-brand">
                <Link to="/">Logo </Link>
            </div>
            <button className="navbar-toggle" onClick={toggleMenu}>
                ☰
            </button>
            <ul className={`navbar-menu ${isOpen ? 'open' : ''}`}>
                <li>
                    <NavLink to="/" exact activeClassName="active">
                        Главная
                    </NavLink>
                </li>
                <li>
                    <NavLink to="/complaints" activeClassName="active">
                        Жалобы
                    </NavLink>
                </li>
                <li>
                    <NavLink to="/bookings" activeClassName="active">
                        Бронирования
                    </NavLink>
                </li>
                <li>
                    <NavLink to="/rooms" activeClassName="active">
                        Номера
                    </NavLink>
                </li>
                <li>
                    <NavLink to="/payments" activeClassName="active">
                        Платежи
                    </NavLink>
                </li>
                <li>
                    <NavLink to="/user" activeClassName="active">
                        demouser
                    </NavLink>
                </li>
                <li>
                    <NavLink to="/logout" activeClassName="active">
                        Выйти
                    </NavLink>
                </li>
            </ul>
        </nav>
    );
};

export default Navbar;