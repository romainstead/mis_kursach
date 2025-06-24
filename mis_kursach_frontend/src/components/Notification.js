import React from 'react';
import { useNotification } from './NotificationContext';
import './Notification.css';

function Notification() {
    const { notification, progress } = useNotification();

    if (!notification) return null;

    return (
        <div className="notification">
            <p>{notification}</p>
            <div className="progress-bar">
                <div className="progress" style={{ width: `${progress}%` }}></div>
            </div>
        </div>
    );
}

export default Notification;