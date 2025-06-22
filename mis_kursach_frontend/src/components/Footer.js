import React from 'react';
import './Footer.css';

function Footer() {
    return (
        <footer className="admin-footer">
            <div className="footer-section">
                © {new Date().getFullYear()} Курсовая работа "Система управления гостиницей". 2025.
            </div>
            <div className="footer-section center">
                Версия 0.0.1
            </div>
            <div className="footer-section right">
                Поддержка: <a href="mailto:rao6@tpu.ru">rao6@tpu.ru</a>
            </div>
        </footer>
    );
}

export default Footer;
