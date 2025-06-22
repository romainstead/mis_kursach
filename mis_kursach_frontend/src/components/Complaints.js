// src/components/Complaints.js
import React, { useState, useEffect } from 'react';
import "./ComplaintsTable.css"; // CSS-файл со стилями

function Complaints() {
    const [openDropdownId, setOpenDropdownId] = useState(null);
    const [complaints, setComplaints] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const toggleDropdown = (id) => {
        setOpenDropdownId(openDropdownId === id ? null : id);
    };


    useEffect(() => {
        fetch('http://localhost:8080/api/GetAllComplaints')
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(data => {
                setComplaints(data);
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
        <table className="complaints-table">
            <thead>
            <tr>
                <th>ID</th>
                <th>Дата и время жалобы</th>
                <th>ID брони</th>
                <th>Статус</th>
                <th>Номер</th>
                <th>Действия</th>
            </tr>
            </thead>
            <tbody>
            {complaints.map((c) => (
                <tr key={c.id}>
                    <td>{c.id}</td>
                    <td>{new Date(c.issue_date).toLocaleString()}</td>
                    <td>{c.booking_id}</td>
                    <td>{c.status}</td>
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

export default Complaints;