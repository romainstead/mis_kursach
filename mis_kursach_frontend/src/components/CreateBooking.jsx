import React, { useEffect, useState } from 'react';
import './CreateBookingForm.css';
import api from "../utils/api";

function CreateBookingForm({ onClose }) {
    const [formData, setFormData] = useState({});
    const [message, setMessage] = useState('');
    const [categories, setCategories] = useState([]);
    const [methods, setMethods] = useState([]);
    const [freeRooms, setFreeRooms] = useState([]); // Ensure initial state is an array

    useEffect(() => {
        api.get("/GetRoomCategories").then(res => setCategories(res.data || []));
        api.get("/GetPaymentMethods").then(res => setMethods(res.data || []));
    }, []);

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData({
            ...formData,
            [name]: type === 'checkbox' ? checked : value
        });
    };

    useEffect(() => {
        if (formData.start_date && formData.end_date && formData.category_code) {
            api.get("/GetFreeRooms", {
                params: {
                    start_date: formData.start_date,
                    end_date: formData.end_date,
                    category_code: formData.category_code
                }
            })
                .then(res => {
                    // Ensure res.data is an array; fallback to empty array if not
                    setFreeRooms(Array.isArray(res.data) ? res.data : []);
                })
                .catch(err => {
                    console.error("Error fetching free rooms:", err);
                    setFreeRooms([]); // Reset to empty array on error
                });
        } else {
            setFreeRooms([]); // Clear free rooms if conditions are not met
        }
    }, [formData.start_date, formData.end_date, formData.category_code]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await api.post("/CreateBooking", {
                ...formData,
                category_code: parseInt(formData.category_code),
                room_number: parseInt(formData.room_number),
                payment_method_code: parseInt(formData.payment_method_code),
                baby_bed: !!formData.baby_bed,
                start_date: formData.start_date,
                end_date: formData.end_date,
                check_in: formData.check_in,
                check_out: formData.check_out,
            });

            if (response.status === 201) {
                setMessage("Бронирование успешно создано")
                window.location.reload();
                onClose();
            } else {
                setMessage("Ошибка при создании бронирования");
                console.error("Unexpected response:", response);
            }
        } catch (err) {
            console.error("Ошибка при создании бронирования:", err);
            setMessage("Ошибка при создании бронирования");
        }
    };

    useEffect(() => {
        if (formData.start_date) {
            const checkIn = `${formData.start_date}T14:00`;
            setFormData(prev => ({ ...prev, check_in: checkIn }));
        }
        if (formData.end_date) {
            const checkOut = `${formData.end_date}T12:00`;
            setFormData(prev => ({ ...prev, check_out: checkOut }));
        }
    }, [formData.start_date, formData.end_date]);

    return (
        <form onSubmit={handleSubmit} className="form-container">
            <label>Дата начала</label>
            <input type="date" name="start_date" onChange={handleChange} required />
            <label>Дата конца</label>
            <input type="date" name="end_date" onChange={handleChange} required />
            <label>Код категории номера</label>
            <select name="category_code" onChange={handleChange} required>
                <option value="">Выберите категорию</option>
                {categories.map(c => (
                    <option key={c.code} value={c.code}>{c.name}</option>
                ))}
            </select>
            <label>Дата и время заезда</label>
            <input
                type="datetime-local"
                name="check_in"
                value={formData.check_in || ''}
                onChange={handleChange}
            />
            <label>Дата и время выезда</label>
            <input
                type="datetime-local"
                name="check_out"
                value={formData.check_out || ''}
                onChange={handleChange}
            />
            <label>Номер комнаты</label>
            {formData.start_date && formData.end_date && formData.category_code ? (
                Array.isArray(freeRooms) && freeRooms.length > 0 ? (
                    <select name="room_number" onChange={handleChange} required>
                        <option value="">Выберите номер</option>
                        {freeRooms.map(r => (
                            <option key={r} value={r}>{r}</option>
                        ))}
                    </select>
                ) : (
                    <div>
                        <select name="room_number" disabled>
                            <option value="">Нет свободных номеров</option>
                        </select>
                        <p className="error-message">Свободные номера не найдены для выбранных дат и категории.</p>
                    </div>
                )
            ) : (
                <select name="room_number" disabled>
                    <option value="">Выберите даты и категорию</option>
                </select>
            )}
            <label className="checkbox-label">
                <input type="checkbox" name="baby_bed" onChange={handleChange} />
                Детская кроватка
            </label>
            <label>ФИО гостя</label>
            <input type="text" name="guest_name" placeholder="Иванов Иван Иванович" onChange={handleChange} required />
            <label>Номер паспорта гостя</label>
            <input type="text" name="guest_passport_number" placeholder="1111 111111" onChange={handleChange} required />
            <label>Номер телефона гостя</label>
            <input type="text" name="guest_phone_number" placeholder="Телефон" onChange={handleChange} required />
            <label>Метод оплаты</label>
            <select name="payment_method_code" onChange={handleChange} required>
                <option value="">Выберите метод оплаты</option>
                {methods.map(m => (
                    <option key={m.code} value={m.code}>{m.name}</option>
                ))}
            </select>
            <button type="submit">Создать бронирование</button>
            {message && <p>{message}</p>}
        </form>
    );
}

export default CreateBookingForm;