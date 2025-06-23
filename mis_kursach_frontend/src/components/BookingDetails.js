import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../utils/api';
import './BookingDetails.css';

function BookingDetails() {
    const { id } = useParams(); // Получаем ID из URL
    const navigate = useNavigate();
    const [booking, setBooking] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchBooking = async () => {
            try {
                const response = await api.get(`/GetBookingByID/${id}`);
                setBooking(response.data);
                setLoading(false);
            } catch (err) {
                setError(err.response?.data?.error || 'Ошибка загрузки бронирования');
                setLoading(false);
            }
        };

        fetchBooking();
    }, [id]);

    const handleBack = () => {
        navigate('/bookings'); // Возвращаемся к списку бронирований
    };

    if (loading) {
        return <div className="text-center mt-8">Загрузка...</div>;
    }

    if (error) {
        return <div className="text-center mt-8 text-red-500">{error}</div>;
    }

    if (!booking) {
        return <div className="text-center mt-8">Бронирование не найдено</div>;
    }

    return (
        <div className="booking-details-wrapper">
            <div className="top-bar">
                <h2>Детали бронирования #{booking.id}</h2>
                <button className="back-btn" onClick={handleBack}>
                    Назад
                </button>
            </div>
            <div className="booking-card">
                <div className="booking-field">
                    <span className="label">Статус:</span>
                    <span>{booking.booking_status}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Дата начала:</span>
                    <span>{new Date(booking.start_date).toLocaleString()}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Дата окончания:</span>
                    <span>{new Date(booking.end_date).toLocaleString()}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Заезд:</span>
                    <span>{booking.check_in ? new Date(booking.check_in).toLocaleString() : 'Не указано'}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Выезд:</span>
                    <span>{booking.check_out ? new Date(booking.check_out).toLocaleString() : 'Не указано'}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Номер комнаты:</span>
                    <span>{booking.room}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Детская кроватка:</span>
                    <span>{booking.baby_bed ? 'Да' : 'Нет'}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Сумма бронирования:</span>
                    <span>{booking.booking_sum.toFixed(2)} ₽</span>
                </div>
                <div className="booking-field">
                    <span className="label">Скидка:</span>
                    <span>{booking.discount_amount ? `${booking.discount_amount.toFixed(2)} %` : 'Нет'}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Итоговая сумма:</span>
                    <span>{booking.total_sum.toFixed(2)} ₽</span>
                </div>
            </div>
        </div>
    );
}

export default BookingDetails;