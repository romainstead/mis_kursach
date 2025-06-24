import React, { createContext, useContext, useState, useEffect } from 'react';

const NotificationContext = createContext();

export function NotificationProvider({ children }) {
    const [notification, setNotification] = useState(null);
    const [progress, setProgress] = useState(100);

    useEffect(() => {
        if (notification) {
            const timer = setInterval(() => {
                setProgress((prev) => {
                    if (prev <= 0) {
                        clearInterval(timer);
                        setNotification(null);
                        return 0;
                    }
                    return prev - 2; // Уменьшаем прогресс каждые 100 мс (5 секунд)
                });
            }, 100);
            return () => clearInterval(timer);
        }
    }, [notification]);

    const showNotification = (message) => {
        setNotification(message);
        setProgress(100);
    };

    return (
        <NotificationContext.Provider value={{ notification, progress, showNotification }}>
            {children}
        </NotificationContext.Provider>
    );
}

export const useNotification = () => useContext(NotificationContext);