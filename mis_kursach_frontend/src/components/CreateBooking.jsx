import React, {useEffect, useState} from 'react';
import './CreateBookingForm.css';
import api from "../utils/api"; // создадим этот файл

function CreateBookingForm({ onClose }) {
    const [formData, setFormData] = useState({});
    const [message, setMessage] = useState('');
    const [categories, setCategories] = useState([]);
    const [methods, setMethods] = useState([]);
    const [freeRooms, setFreeRooms] = useState([]);

    useEffect(() => {
        api.get("/GetRoomCategories").then(res => setCategories(res.data));
        api.get("/GetPaymentMethods").then(res => setMethods(res.data));
    }, []);
    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData({
            ...formData,
            [name]: type === 'checkbox' ? checked : value
        });
    };

    useEffect(() => {
        if (formData.start_date && formData.end_date) {
            api.get("/GetFreeRooms", {
                params: {
                    start_date: formData.start_date,
                    end_date: formData.end_date
                }
            }).then(res => setFreeRooms(res.data));
        }
    }, [formData.start_date, formData.end_date]);


    const handleSubmit = (e) => {
        e.preventDefault();
        // Здесь ваш запрос на сервер
        setMessage('Бронирование успешно создано');
        onClose();
    };

    return (
        <form onSubmit={handleSubmit} className="form-container">
            Дата начала
            <input type="date" name="start_date"  onChange={handleChange} required />
            Дата конца
            <input type="date" name="end_date" onChange={handleChange} required />
            Код категории номера
            <select name="category_code" onChange={handleChange} required>
                <option value="">Выберите категорию</option>
                {categories.map(c => (
                    <option key={c.code} value={c.code}>{c.name}</option>
                ))}
            </select>
            Дата и время заезда
            <input type="datetime-local" name="check_in" onChange={handleChange}/>
            Дата и время выезда
            <input type="datetime-local" name="check_out" onChange={handleChange} />
            {formData.start_date && formData.end_date && (
                <select name="room_number" onChange={handleChange} required>
                    <option value="">Выберите номер</option>
                    {freeRooms.map(r => (
                        <option key={r} value={r}>{r}</option>
                    ))}
                </select>
            )}
            <label className="checkbox-label">
                <input type="checkbox" name="baby_bed" onChange={handleChange} />
                Детская кроватка
            </label>
            ФИО гостя
            <input type="text" name="guest_name" placeholder="Иванов Иван Иванович" onChange={handleChange} required />
            Номер паспорта гостя
            <input type="text" name="guest_passport_number" placeholder="1111 111111" onChange={handleChange} required />
            Номер телефона гостя
            <input type="text" name="guest_phone_number" placeholder="Телефон" onChange={handleChange} required />
            Метод оплаты
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
