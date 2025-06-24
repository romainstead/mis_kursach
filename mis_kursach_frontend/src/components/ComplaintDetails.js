import React, { useState, useEffect } from 'react';
import {useParams, useNavigate, Link} from 'react-router-dom';
import api from '../utils/api';
import './ComplaintDetails.css';

function ComplaintDetails() {
    const { id } = useParams(); // Получаем ID из URL
    const navigate = useNavigate();
    const [complaint, setComplaint] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchBooking = async () => {
            try {
                const response = await api.get(`/GetComplaintByID/${id}`);
                setComplaint(response.data);
                setLoading(false);
            } catch (err) {
                setError(err.response?.data?.error || 'Ошибка загрузки жалобы');
                setLoading(false);
            }
        };

        fetchBooking();
    }, [id]);

    const handleBack = () => {
        navigate('/complaints'); // Возвращаемся к списку бронирований
    };

    if (loading) {
        return <div className="text-center mt-8">Загрузка...</div>;
    }

    if (error) {
        return <div className="text-center mt-8 text-red-500">{error}</div>;
    }

    if (!complaint) {
        return <div className="text-center mt-8">Жалоба не найдена</div>;
    }

    return (
        <div className="booking-details-wrapper">
            <div className="top-bar">
                <h2>Детали жалобы #{complaint.id}</h2>
                <button className="back-btn" onClick={handleBack}>
                    Назад
                </button>
            </div>
            <div className="booking-card">
                <div className="booking-field">
                    <span className="label">Причина:</span>
                    <span>{complaint.reason}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Комментарий:</span>
                    <span>{complaint.commentary}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Дата подачи жалобы:</span>
                    <span>{new Date(complaint.issue_date).toLocaleDateString()}</span>
                </div>
                <div className="booking-field">
                    <span className="label">ID бронирования:</span>
                    <span>{<Link to={`/bookings/${complaint.booking_id}`}>{complaint.booking_id}</Link>}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Статус:</span>
                    <span>{complaint.status}</span>
                </div>
                <div className="booking-field">
                    <span className="label">Номер комнаты:</span>
                    <span>{complaint.room}</span>
                </div>
            </div>
        </div>
    );
}

export default ComplaintDetails;