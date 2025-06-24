import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './CreateComplaintForm.css';
import api from '../utils/api';
import { useNotification } from './NotificationContext';
import PropTypes from 'prop-types';

function CreateComplaintForm({ onClose }) {
    const [formData, setFormData] = useState({});
    const [message, setMessage] = useState('');
    const [bookings, setBookings] = useState([]);
    const navigate = useNavigate();
    const { showNotification } = useNotification();

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData({
            ...formData,
            [name]: type === 'checkbox' ? checked : value,
        });
    };

    useEffect(() => {
        api.get('/GetAllBookings').then((res) => setBookings(res.data || []));
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await api.post('/CreateComplaint', {
                ...formData,
                reason: formData.reason,
                booking_id: parseInt(formData.booking_id),
                commentary: formData.commentary,
            });

            if (response.status === 201) {
                showNotification('Жалоба успешно создана'); // Показываем уведомление
                navigate('/complaints'); // Немедленное перенаправление
                if (typeof onClose === 'function') {
                    onClose(); // Закрываем форму, если onClose — функция
                }
            } else {
                setMessage('Ошибка при создании жалобы');
                console.error('Unexpected response:', response);
            }
        } catch (err) {
            console.error('Ошибка при создании жалобы:', err);
            setMessage('Ошибка при создании жалобы');
        }
    };

    return (
        <form onSubmit={handleSubmit} className="form-container">
            <label>Причина</label>
            <input type="text" name="reason" onChange={handleChange} required />
            <label>ID бронирования</label>
            <select name="booking_id" onChange={handleChange} required>
                <option value="">Выберите ID бронирования</option>
                {bookings.map((b) => (
                    <option key={b.id} value={b.id}>{`${b.id}, ${b.guest_name}`}</option>
                ))}
            </select>
            <label>Комментарий</label>
            <input type="text" name="commentary" onChange={handleChange} />
            <button type="submit">Создать жалобу</button>
            {message && <p>{message}</p>}
        </form>
    );
}

CreateComplaintForm.propTypes = {
    onClose: PropTypes.func,
};

export default CreateComplaintForm;