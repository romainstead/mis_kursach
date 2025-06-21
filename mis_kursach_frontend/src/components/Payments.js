import { Link } from 'react-router-dom';
import React, { useState, useEffect } from 'react';
import "./PaymentsTable.css"
import axios from 'axios';

const api = axios.create({
    baseURL: 'http://localhost:8080/api',
});

function Payments() {
    const [payments, setPayments] = useState([]);
    const [openDropdownId, setOpenDropdownId] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const toggleDropdown = (id) => {
        setOpenDropdownId(openDropdownId === id ? null : id);
    };

    useEffect(() => {
        api.get('/GetAllPayments')
            .then(response => {
                setPayments(response.data);
                setLoading(false);
            })
            .catch(err => {
                setError('Failed to fetch payments');
                setLoading(false);
            });
    }, []);

    if (loading) return <div className="text-center mt-8">Loading...</div>;
    if (error) return <div className="text-center mt-8 text-red-500">{error}</div>;

    return (
        <table className="payments-table">
            <thead>
            <tr>
                <th>ID оплаты</th>
                <th>ID брони</th>
                <th>Дата оплаты</th>
                <th>Статус</th>
                <th>Сумма</th>
                <th>Действия</th>
            </tr>
            </thead>
            <tbody>
            {payments.map((c) => (
                <tr key={c.id}>
                    <td>{c.id}</td>
                    <td>{c.booking_id}</td>
                    <td>{new Date(c.pay_date).toLocaleString()}</td>
                    <td>{c.status_name}</td>
                    <td>{c.amount}</td>
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


export default Payments