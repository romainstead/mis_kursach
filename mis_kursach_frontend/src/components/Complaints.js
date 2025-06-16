// src/components/Complaints.js
import React, { useState, useEffect } from 'react';

function Complaints() {
    const [complaints, setComplaints] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

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
        <div>
            <h2>Жалобы</h2>
            {complaints.length === 0 ? (
                <p>Жалобы не найдены</p>
            ) : (
                <ul>
                    {complaints.map(complaint => (
                        <li key={complaint.id}>
                            ID: {complaint.id}, Причина: {complaint.reason}, Комментарий: {complaint.commentary}, Дата подачи: {complaint.issue_date}, ID бронирования: {complaint.booking_id}, Статус: {complaint.status_code}
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}

export default Complaints;