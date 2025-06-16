// src/components/Bookings.js
import React, { useState, useEffect } from 'react';

function Bookings() {
    const [bookings, setBookings] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetch('http://localhost:8080/api/GetAllBookings')
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                setBookings(data);
                setLoading(false);
            })
            .catch(error => {
                setError(error.message);
                setLoading(false);
            });
    }, []);

    if (loading) return <div>Загрузка...</div>;
    if (error) return <div>Ошибка: {error}</div>;

    return (
        <div>
            <h2>Бронирования</h2>
            {bookings.length === 0 ? (
                <p>Бронирования не найдены</p>
            ) : (
                <ul>
                    {bookings.map(booking => (
                        <li key={booking.id}>
                            {/* Замени поля на реальные из твоей структуры данных */}
                            ID: {booking.id}, Начало брони: {booking.start_date}, Конец брони: {booking.end_date}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default Bookings;