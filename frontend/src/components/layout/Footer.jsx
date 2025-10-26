import { Link } from 'react-router-dom';
import '../../styles/footer.css';

function Footer() {
    return (
        <footer className="footer">
            <div className="footer-container">
                {/* Секция подписки */}
                <div className="footer-subscribe">
                    <div className="subscribe-content">
                        <h3>Будьте в курсе лучших предложений</h3>
                        <p>Подпишитесь на рассылку и получайте уведомления о новых товарах</p>
                    </div>
                    <form className="subscribe-form">
                        <input type="email" placeholder="Введите ваш email" className="subscribe-input" required />
                        <button type="submit" className="btn btn-primary">Подписаться</button>
                    </form>
                </div>

                {/* Основная часть футера */}
                <div className="footer-content">
                    {/* О компании */}
                    <div className="footer-section footer-about">
                        <div className="logo">
                            <span className="logo-avi">Avito</span>
                            <span className="logo-go">Go</span>
                        </div>
                        <p className="footer-description">
                            Современная C2C платформа для покупки и продажи товаров.
                            Быстро, надежно, удобно.
                        </p>
                        <div className="footer-social">
                            <a href="#" className="social-link" aria-label="VK">
                                <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M15.07 2H8.93C3.33 2 2 3.33 2 8.93v6.14C2 20.67 3.33 22 8.93 22h6.14c5.6 0 6.93-1.33 6.93-6.93V8.93C22 3.33 20.67 2 15.07 2zm3.15 14.63h-1.49c-.55 0-.72-.45-1.71-1.44-.86-.83-1.24-.94-1.45-.94-.3 0-.38.08-.38.46v1.31c0 .36-.11.57-1.07.57-1.58 0-3.33-.96-4.56-2.75-1.86-2.66-2.37-4.66-2.37-5.07 0-.21.08-.41.46-.41h1.49c.35 0 .48.16.61.54.71 2.05 1.89 3.84 2.37 3.84.18 0 .27-.08.27-.54v-2.09c-.06-.98-.58-1.06-.58-1.41 0-.17.14-.34.37-.34h2.33c.29 0 .4.16.4.5v2.81c0 .29.13.4.21.4.18 0 .33-.11.67-.45 1.04-1.17 1.78-2.98 1.78-2.98.1-.21.26-.41.61-.41h1.49c.42 0 .51.22.42.52-.16.71-1.98 3.56-1.98 3.56-.15.24-.18.35 0 .62.13.19.54.53.82.85.51.51 1.04 1.03 1.16 1.35.13.33-.07.5-.4.5z"/>
                                </svg>
                            </a>
                            <a href="#" className="social-link" aria-label="Telegram">
                                <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm4.64 6.8c-.15 1.58-.8 5.42-1.13 7.19-.14.75-.42 1-.68 1.03-.58.05-1.02-.38-1.58-.75-.88-.58-1.38-.94-2.23-1.5-.99-.65-.35-1.01.22-1.59.15-.15 2.71-2.48 2.76-2.69a.2.2 0 00-.05-.18c-.06-.05-.14-.03-.21-.02-.09.02-1.49.95-4.22 2.79-.4.27-.76.41-1.08.4-.36-.01-1.04-.2-1.55-.37-.63-.2-1.12-.31-1.08-.66.02-.18.27-.36.74-.55 2.92-1.27 4.86-2.11 5.83-2.51 2.78-1.16 3.35-1.36 3.73-1.36.08 0 .27.02.39.12.1.08.13.19.14.27-.01.06.01.24 0 .38z"/>
                                </svg>
                            </a>
                        </div>
                    </div>

                    {/* Покупателям */}
                    <div className="footer-section">
                        <h4>Покупателям</h4>
                        <ul>
                            <li><Link to="/catalog">Каталог товаров</Link></li>
                            <li><Link to="/how-to-buy">Как покупать</Link></li>
                            <li><Link to="/delivery">Доставка</Link></li>
                            <li><Link to="/payment">Оплата</Link></li>
                        </ul>
                    </div>

                    {/* Продавцам */}
                    <div className="footer-section">
                        <h4>Продавцам</h4>
                        <ul>
                            <li><Link to="/start-selling">Начать продавать</Link></li>
                            <li><Link to="/seller-guide">Руководство продавца</Link></li>
                            <li><Link to="/pricing">Тарифы</Link></li>
                        </ul>
                    </div>

                    {/* О компании */}
                    <div className="footer-section">
                        <h4>О компании</h4>
                        <ul>
                            <li><Link to="/about">О нас</Link></li>
                            <li><Link to="/blog">Блог</Link></li>
                            <li><Link to="/contacts">Контакты</Link></li>
                        </ul>
                    </div>

                    {/* Помощь */}
                    <div className="footer-section">
                        <h4>Помощь</h4>
                        <ul>
                            <li><Link to="/faq">Частые вопросы</Link></li>
                            <li><Link to="/support">Служба поддержки</Link></li>
                        </ul>
                        <div className="footer-contact">
                            <p><strong>8 (800) 555-35-35</strong></p>
                            <p>support@avitogo.ru</p>
                        </div>
                    </div>
                </div>

                {/* Нижняя часть */}
                <div className="footer-bottom">
                    <div className="footer-legal">
                        <p>&copy; 2025 AvitoGo. Все права защищены.</p>
                        <div className="footer-legal-links">
                            <Link to="/privacy">Политика конфиденциальности</Link>
                            <span>•</span>
                            <Link to="/terms">Пользовательское соглашение</Link>
                        </div>
                    </div>
                    <div className="footer-badges">
                        <span className="badge-item">🔒 Безопасные платежи</span>
                        <span className="badge-item">✓ Проверенные продавцы</span>
                        <span className="badge-item">⚡ На Go</span>
                    </div>
                </div>
            </div>
        </footer>
    );
}

export default Footer;