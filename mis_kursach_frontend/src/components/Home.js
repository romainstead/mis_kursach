// src/components/Home.js
import React from 'react';
import { Link } from 'react-router-dom';

function Home() {
    return (
        <div>
            <h2>Главная страница</h2>
            <nav>
                <ul>
                    <li><Link to="/bookings">Бронирования</Link></li>
                    <li><Link to="/complaints">Жалобы</Link></li>
                    <li><Link to="/payments">Платежи</Link></li>
                    <li><Link to="/rooms">Комнаты</Link></li>
                </ul>
            </nav>
        </div>
    );
}

export default Home;