import React, { useState } from 'react';
import './Login.css';
import api from '../utils/api';

function Login() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');
    const [isPopupOpen, setIsPopupOpen] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setMessage('');

        try {
            const response = await api.post('/login', { username, password });

            if (!response.data) {
                setMessage('Попытка входа провалена.');
                return;
            }

            localStorage.setItem('token', response.data.token);
            localStorage.setItem("username", response.data.username);
            setMessage('Успешный вход, перенаправление...');
            setTimeout(() => window.location.href = '/main', 1000);
        } catch (error) {
            console.error('Login error:', error);
            setMessage(error.response?.data?.error || 'Error: ' + error.message);
        }
    };

    const handleOpenPopup = () => setIsPopupOpen(true);
    const handleClosePopup = () => setIsPopupOpen(false);

    return (
        <div className="login-container">
            <h1>LOGO</h1>
            <h2>Вход в систему</h2>
            <form className="login-form" onSubmit={handleSubmit}>
                <input
                    type="text"
                    placeholder="Логин"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                />
                <input
                    type="password"
                    placeholder="Пароль"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                />
                <div className="forgot-password" onClick={handleOpenPopup}>Забыли пароль?</div>
                <button type="submit">Войти</button>
                {message && <div className={`message ${message.includes("успеш") ? "success" : "error"}`}>{message}</div>}
            </form>

            {isPopupOpen && (
                <div className="modal-overlay" onClick={handleClosePopup}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <button className="modal-close" onClick={handleClosePopup}>×</button>
                        <p>Для смены пароля обратитесь к администратору системы или обратитесь на почту rao6@tpu.ru.</p>
                    </div>
                </div>
            )}
        </div>
    );
}

export default Login;