import { Link } from 'react-router-dom';
import { useContext } from 'react';
import { AuthContext } from '../../context/AuthContext';
import '../../styles/header.css';

function Header() {
    const { user } = useContext(AuthContext);

    const profileLink = user ? `/account/${user.id}` : '/auth/login';
    const createListingLink = user ? `/account/${user.id}/listings/create` : '/auth/login';

    return (
        <header className="header">
            <nav className="navbar">
                <Link to="/" className="logo">
                <span className="logo-avi">Avi</span>
                    <span className="logo-go">Go</span>
                    <img
                        src="/assets/images/gopher.png"
                        alt="Go Gopher"
                        className="gopher"
                    />
                </Link>

                <div className="search-bar">
                    <svg
                        className="search-icon"
                        width="20"
                        height="20"
                        viewBox="0 0 24 24"
                        fill="none"
                    >
                        <path
                            d="M21 21L16.65 16.65M19 11C19 15.4183 15.4183 19 11 19C6.58172 19 3 15.4183 3 11C3 6.58172 6.58172 3 11 3C15.4183 3 19 6.58172 19 11Z"
                            stroke="currentColor"
                            strokeWidth="2"
                            strokeLinecap="round"
                        />
                    </svg>
                    <input
                        type="text"
                        placeholder="Поиск товаров..."
                        className="search-input"
                    />
                    <button className="search-button">Найти</button>
                </div>

                <ul className="nav-menu">
                    <li>
                        <Link to="/catalog" className="nav-link">
                            <svg
                                width="20"
                                height="20"
                                viewBox="0 0 24 24"
                                fill="none"
                            >
                                <path
                                    d="M3 7H21M3 12H21M3 17H21"
                                    stroke="currentColor"
                                    strokeWidth="2"
                                    strokeLinecap="round"
                                />
                            </svg>
                            Каталог
                        </Link>
                    </li>
                    <li>
                        <Link to="/seller" className="nav-link">
                            Продавцам
                        </Link>
                    </li>
                    <li>
                        <Link to="/help" className="nav-link">
                            Помощь
                        </Link>
                    </li>
                    {/* по желанию можно добавить:
                    <li>
                        <Link to="/account" className="nav-link">
                            Личный кабинет
                        </Link>
                    </li>
                    */}
                </ul>

                <div className="nav-buttons">
                    <button className="btn-icon" aria-label="Избранное">
                        <svg
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="none"
                        >
                            <path
                                d="M12 21.35L10.55 20.03C5.4 15.36 2 12.28 2 8.5C2 5.42 4.42 3 7.5 3C9.24 3 10.91 3.81 12 5.09C13.09 3.81 14.76 3 16.5 3C19.58 3 22 5.42 22 8.5C22 12.28 18.6 15.36 13.45 20.03L12 21.35Z"
                                stroke="currentColor"
                                strokeWidth="2"
                                fill="none"
                            />
                        </svg>
                        <span className="badge">3</span>
                    </button>

                    <button className="btn-icon" aria-label="Сообщения">
                        <svg
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="none"
                        >
                            <path
                                d="M21 15C21 15.5304 20.7893 16.0391 20.4142 16.4142C20.0391 16.7893 19.5304 17 19 17H7L3 21V5C3 4.46957 3.21071 3.96086 3.58579 3.58579C3.96086 3.21071 4.46957 3 5 3H19C19.5304 3 20.0391 3.21071 20.4142 3.58579C20.7893 3.96086 21 4.46957 21 5V15Z"
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                            />
                        </svg>
                        <span className="badge">1</span>
                    </button>

                    <Link
                        to={profileLink}
                        className="btn-icon user-menu"
                        aria-label="Профиль"
                    >
                        <svg
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="none"
                        >
                            <path
                                d="M20 21V19C20 17.9391 19.5786 16.9217 18.8284 16.1716C18.0783 15.4214 17.0609 15 16 15H8C6.93913 15 5.92172 15.4214 5.17157 16.1716C4.42143 16.9217 4 17.9391 4 19V21M16 7C16 9.20914 14.2091 11 12 11C9.79086 11 8 9.20914 8 7C8 4.79086 9.79086 3 12 3C14.2091 3 16 4.79086 16 7Z"
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                            />
                        </svg>
                    </Link>

                    <Link
                        to={createListingLink}
                        className={"btn btn-primary"}
                        state={!user ? { from: { pathname: '/auth/login' } } : null}>
                        + Разместить объявление
                    </Link>
                </div>

                <button className="hamburger" aria-label="Меню">
                    <span></span>
                    <span></span>
                    <span></span>
                </button>
            </nav>
        </header>
    );
}

export default Header;
