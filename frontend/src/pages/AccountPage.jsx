import { useEffect, useState } from 'react';
import { getProfile } from '../services/api';
import '../styles/auth.css';
import '../styles/profile.css';
import { useNavigate } from "react-router-dom";

function AccountPage() {
    const [profile, setProfile] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const navigate = useNavigate();

    const handleEditProfile = () => {
        navigate('/account/edit');
    };

    const handleChangePassword = () => {
        navigate('/account/change-password');
    };

    const handleMyPurchases = () => {
        navigate('/account/purchases');
    };

    const handleMySales = () => {
        navigate('/account/sales');
    };

    const handleMyListings = () => {
        navigate('/account/listings');
    };

    useEffect(() => {
        const loadProfile = async () => {
            try {
                const data = await getProfile();
                setProfile(data);
            } catch (err) {
                console.error('Failed to load profile', err);
                setError(err.message || 'Ошибка загрузки профиля');
            } finally {
                setLoading(false);
            }
        };

        loadProfile();
    }, []);

    const formatLastLogin = () => {
        if (!profile || !profile.last_login_at) {
            return 'Нет данных';
        }
        try {
            return new Date(profile.last_login_at).toLocaleString('ru-RU');
        } catch {
            return profile.last_login_at;
        }
    };

    if (loading) {
        return (
            <div className="auth-page">
                <div className="auth-card">
                    <div className="auth-header">
                        <h1 className="auth-title">Личный кабинет</h1>
                        <p className="auth-subtitle">Загрузка профиля...</p>
                    </div>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="auth-page">
                <div className="auth-card">
                    <div className="auth-header">
                        <h1 className="auth-title">Личный кабинет</h1>
                        <p className="auth-subtitle">Ошибка загрузки профиля</p>
                    </div>
                    <p className="form-error">{error}</p>
                </div>
            </div>
        );
    }

    return (
        <div className="auth-page">
            <div className="profile-card-wide">
                <div className="profile-main">
                    {/* Левая колонка: аватар и имя */}
                    <div className="profile-avatar-column">
                        <div className="profile-avatar-wrapper">
                            {profile?.avatar_url ? (
                                <img
                                    src={profile.avatar_url}
                                    alt="Аватар"
                                    className="profile-avatar-img"
                                />
                            ) : (
                                <div className="profile-avatar-placeholder">
                                    {profile?.display_name
                                        ? profile.display_name[0].toUpperCase()
                                        : 'U'}
                                </div>
                            )}
                        </div>
                        <div className="profile-basic-info">
                            <div className="profile-name">
                                {profile.display_name || 'Без имени'}
                            </div>
                            <div className="profile-location">
                                {[profile.city, profile.country]
                                    .filter(Boolean)
                                    .join(', ') || 'Город и страна не указаны'}
                            </div>
                            <div className="profile-last-login">
                                Последний вход: {formatLastLogin()}
                            </div>
                        </div>
                        <div className="profile-actions-vertical">
                            <button className="btn btn-primary profile-btn" onClick={handleEditProfile}>
                                Обновить профиль
                            </button>
                            <button className="btn btn-secondary profile-btn" onClick={handleChangePassword}>
                                Сменить пароль
                            </button>
                        </div>
                    </div>

                    <div className="profile-details-column">
                        <h2 className="profile-section-title">Информация профиля</h2>

                        <div className="profile-fields-grid">
                            <div className="profile-field">
                                <div className="profile-field-label">Телефон</div>
                                <div className="profile-field-value">
                                    {profile.phone || 'Не указан'}
                                </div>
                            </div>
                            <div className="profile-field">
                                <div className="profile-field-label">О себе</div>
                                <div className="profile-field-value">
                                    {profile.about || 'Расскажите о себе'}
                                </div>
                            </div>
                            <div className="profile-field">
                                <div className="profile-field-label">Страна</div>
                                <div className="profile-field-value">
                                    {profile.country || 'Не указана'}
                                </div>
                            </div>
                            <div className="profile-field">
                                <div className="profile-field-label">Город</div>
                                <div className="profile-field-value">
                                    {profile.city || 'Не указан'}
                                </div>
                            </div>
                        </div>

                        <div className="profile-orders-section"
                             style={{marginTop: '2.5rem', borderTop: '1px solid #eee', paddingTop: '1.5rem'}}>
                            <h2 className="profile-section-title">Активность</h2>

                            <div className="orders-buttons-grid"
                                 style={{display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '15px'}}>

                                {/* Кнопка "Мои Объявления" */}
                                <div
                                    onClick={handleMyListings}
                                    style={{
                                        display: 'flex',
                                        flexDirection: 'column',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        padding: '20px',
                                        backgroundColor: '#f8f9fa',
                                        border: '1px solid #e9ecef',
                                        borderRadius: '12px',
                                        cursor: 'pointer',
                                        transition: 'all 0.2s',
                                        textAlign: 'center'
                                    }}
                                    className="order-card-btn"
                                >
                                    <span style={{fontSize: '28px', marginBottom: '8px'}}>📋</span>
                                    <span style={{fontWeight: '600', color: '#333'}}>Мои Объявления</span>
                                </div>

                                {/* Кнопка "Мои Покупки" */}
                                <div
                                    onClick={handleMyPurchases}
                                    style={{
                                        display: 'flex',
                                        flexDirection: 'column',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        padding: '20px',
                                        backgroundColor: '#f8f9fa',
                                        border: '1px solid #e9ecef',
                                        borderRadius: '12px',
                                        cursor: 'pointer',
                                        transition: 'all 0.2s',
                                        textAlign: 'center'
                                    }}
                                    className="order-card-btn"
                                >
                                    <span style={{fontSize: '28px', marginBottom: '8px'}}>🛍️</span>
                                    <span style={{fontWeight: '600', color: '#333'}}>Мои Покупки</span>
                                </div>

                                {/* Кнопка "Мои Продажи" */}
                                <div
                                    onClick={handleMySales}
                                    style={{
                                        display: 'flex',
                                        flexDirection: 'column',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        padding: '20px',
                                        backgroundColor: '#f8f9fa',
                                        border: '1px solid #e9ecef',
                                        borderRadius: '12px',
                                        cursor: 'pointer',
                                        transition: 'all 0.2s',
                                        textAlign: 'center'
                                    }}
                                    className="order-card-btn"
                                >
                                    <span style={{fontSize: '28px', marginBottom: '8px'}}>💰</span>
                                    <span style={{fontWeight: '600', color: '#333'}}>Мои Продажи</span>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default AccountPage;
