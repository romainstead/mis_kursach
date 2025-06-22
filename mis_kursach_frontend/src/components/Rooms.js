import { Link } from 'react-router-dom';
import React, { useState, useEffect } from 'react';
import "./RoomsTable.css"
import api from "../utils/api";

function Rooms() {
    const [rooms, setRooms] = useState([]);
    const [openDropdownId, setOpenDropdownId] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const toggleDropdown = (id) => {
        setOpenDropdownId(openDropdownId === id ? null : id);
    };

    useEffect(() => {
        api.get('/GetAllRooms')
            .then(response => {
                setRooms(response.data);
                setLoading(false);
            })
            .catch(err => {
                setError('Failed to fetch bookings');
                setLoading(false);
            });
    }, []);

    if (loading) return <div className="text-center mt-8">Loading...</div>;
    if (error) return <div className="text-center mt-8 text-red-500">{error}</div>;

    return (
        <table className="rooms-table">
            <thead>
            <tr>
                <th>№ комнаты</th>
                <th>Категория</th>
                <th>Состояние</th>
                <th>Вместимость</th>
                <th>Действия</th>
            </tr>
            </thead>
            <tbody>
            {rooms.map((r) => (
                <tr key={r.number}>
                    <td>{r.number}</td>
                    <td>{r.category_name}</td>
                    <td>{r.state_name}</td>
                    <td>{r.capacity}</td>
                    <td className="dropdown-cell">
                        <button className="dropdown-toggle" onClick={() => toggleDropdown(r.number)}>
                            ⋮
                        </button>
                        {openDropdownId === r.id && (
                            <ul className="dropdown-menu">
                                <li>Редактировать</li>
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


export default Rooms;