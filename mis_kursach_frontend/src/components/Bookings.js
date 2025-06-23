import React, { useState, useEffect } from 'react';
import "./BookingsTable.css";
import api from "../utils/api";
import CreateBookingForm from "./CreateBooking";
import {Link} from "react-router-dom";

function Bookings() {
    const [bookings, setBookings] = useState([]);
    const [openDropdownId, setOpenDropdownId] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [confirming, setConfirming] = useState(false);

    const toggleDropdown = (id) => {
        setOpenDropdownId(openDropdownId === id ? null : id);
    };

    // Загрузка бронирований
    const fetchBookings = async () => {
        try {
            const response = await api.get('http://127.0.0.1:8080/api/GetAllBookings');
            setBookings(response.data || []);
            setLoading(false);
        } catch (err) {
            setError('Failed to fetch bookings');
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchBookings();
    }, []);

    // Обработка подтверждения бронирования
    const handleConfirmBooking = async (bookingId) => {
        try {
            setConfirming(true);
            await api.post(`/ConfirmBooking?id=${bookingId}`);
            await fetchBookings();
            setOpenDropdownId(null);
        } catch (err) {
            console.error('Ошибка при подтверждении бронирования:', err);
            setError('Failed to confirm booking');
        } finally {
            setConfirming(false);
        }
    };

// В JSX для <li>:


    const handleOpenModal = () => setIsModalOpen(true);
    const handleCloseModal = () => setIsModalOpen(false);

    // Поддержка создания нового бронирования
    const handleBookingCreated = (newBooking) => {
        setBookings(prevBookings => [...prevBookings, newBooking]);
    };

    return (
        <div className="bookings-wrapper">
            <div className="top-bar">
                <h2>Бронирования</h2>
                <button className="add-booking-btn" onClick={handleOpenModal}>
                    Добавить бронирование
                </button>
            </div>

            {loading && <div className="text-center mt-8">Loading...</div>}
            {error && <div className="text-center mt-8 text-red-500">{error}</div>}

            {!loading && !error && (
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
                                        <li>
                                            <Link to={`/bookings/${c.id}`}>Посмотреть</Link>
                                        </li>
                                        <li onClick={() => handleConfirmBooking(c.id)}>{confirming ? 'Заселяется...' : 'Заселить...'}</li>
                                        <li>Удалить</li>
                                    </ul>
                                )}
                            </td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            )}

            {isModalOpen && (
                <div className="modal-overlay" onClick={handleCloseModal}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <button className="modal-close" onClick={handleCloseModal}>×</button>
                        <CreateBookingForm onClose={handleCloseModal} onBookingCreated={handleBookingCreated} />
                    </div>
                </div>
            )}
        </div>
    );
}

export default Bookings;
