import React, { useState, useEffect } from 'react';
import './ComplaintsTable.css';
import api from '../utils/api';
import { Link } from 'react-router-dom';
import CreateComplaintForm from './CreateComplaint';

function Complaints() {
    const [openDropdownId, setOpenDropdownId] = useState(null);
    const [complaints, setComplaints] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [isModalOpen, setIsModalOpen] = useState(false);

    const toggleDropdown = (id) => {
        setOpenDropdownId(openDropdownId === id ? null : id);
    };

    const fetchComplaints = async () => {
        try {
            const response = await api.get('/GetAllComplaints');
            setComplaints(response.data || []);
            setLoading(false);
        } catch (error) {
            setError(error.message);
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchComplaints();
    }, []);

    const handleResolveComplaint = async (complaintId) => {
        try {
            await api.post(`/ResolveComplaint?id=${complaintId}&statusCode=3`);
            await fetchComplaints();
            setOpenDropdownId(null);
        } catch (err) {
            console.error('Ошибка при решении жалобы:', err);
            setError('Failed to resolve complaint');
        }
    };

    const handleOpenModal = () => setIsModalOpen(true);
    const handleCloseModal = () => setIsModalOpen(false);

    const handleComplaintCreated = () => {
        fetchComplaints(); // Обновляем список жалоб после создания
        handleCloseModal();
    };

    if (loading) return <div>Загрузка...</div>;
    if (error) return <div>Ошибка: {error}</div>;

    return (
        <div className="complaints-wrapper">
            <div className="top-bar">
                <h2>Жалобы</h2>
                <button className="add-complaint-btn" onClick={handleOpenModal}>
                    Создать жалобу
                </button>
            </div>

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
                                    <li>
                                        <Link to={`/complaints/${c.id}`}>Посмотреть</Link>
                                    </li>
                                    <li onClick={() => handleResolveComplaint(c.id)}>Решить жалобу</li>
                                </ul>
                            )}
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>

            {isModalOpen && (
                <div className="modal-overlay" onClick={handleCloseModal}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <button className="modal-close" onClick={handleCloseModal}>×</button>
                        <CreateComplaintForm onClose={handleComplaintCreated} />
                    </div>
                </div>
            )}
        </div>
    );
}

export default Complaints;