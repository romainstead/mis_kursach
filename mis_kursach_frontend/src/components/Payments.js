// src/components/Payments.js
import React, { useState, useEffect } from 'react';

function Payments() {
    const [payments, setPayments] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        fetch('http://localhost:8080/api/GetAllPayments')
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                setPayments(data);
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
            <h2>Жалобы</h2>
            {payments.length === 0 ? (
                <p>Жалобы не найдены</p>
            ) : (
                <ul>
                    {payments.map(payment => (
                        <li key={payment.id}>
                            ID: {payment.id}, ID бронирования: {payment.booking_id}, Дата оплаты: {payment.pay_date}, Сумма: {payment.amount}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default Payments;