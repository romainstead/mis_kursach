import { Link } from 'react-router-dom';
import React, { useState, useEffect } from 'react';
import "./BookingsTable.css"
import axios from 'axios';

const api = axios.create({
    baseURL: 'http://localhost:8080/api',
});

function Bookings() {
    const [bookings, setBookings] = useState([]);
    const [openDropdownId, setOpenDropdownId] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const toggleDropdown = (id) => {
        setOpenDropdownId(openDropdownId === id ? null : id);
    };

    useEffect(() => {
        api.get('/GetAllBookings')
            .then(response => {
                setBookings(response.data);
                setLoading(false);
            })
            .catch(err => {
                setError('Failed to fetch bookings');
                setLoading(false);
            });
    }, []);

    if (loading) return <div className="text-center mt-8">Loading...</div>;
    if (error) return <div className="text-center mt-8 text-red-500">{error}</div>;

    return (
        <table className="bookings-table">
            <thead>
            <tr>
                <th>ID</th>
                <th>Дата начала брони</th>
                <th>Дата конца брони</th>
                <th>Статус</th>
                <th>Комната</th>
                <th>Действия</th>
            </tr>
            </thead>
            <tbody>
            {bookings.map((c) => (
                <tr key={c.id}>
                    <td>{c.id}</td>
                    <td>{new Date(c.start_date).toLocaleString()}</td>
                    <td>{new Date(c.end_date).toLocaleString()}</td>
                    <td>{c.booking_status}</td>
                    <td>{c.room}</td>
                    <td className="dropdown-cell">
                        <button className="dropdown-toggle" onClick={() => toggleDropdown(c.id)}>
                            ⋮
                        </button>
                        {openDropdownId === c.id && (
                            <ul className="dropdown-menu">
                                <li>Посмотреть</li>
                                <li>Изменить статус</li>
                                <li>Удалить</li>
                            </ul>
                        )}
                    </td>
                </tr>
            ))}
            </tbody>
        </table>
    );
}


export default Bookings