import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getProfile, getPublicProfile } from '../services/api';
import { useAuth } from '../context/AuthContext';
import '../styles/auth.css';
import '../styles/profile.css';
import toast from 'react-hot-toast';

function AccountPage() {
    const { id } = useParams();
    const { user, logout } = useAuth();
    const navigate = useNavigate();

    const [profile, setProfile] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const isMyProfile = user && (user.id === id || !id);

    useEffect(() => {
        const loadProfile = async () => {
            try {
                let data;
                if (isMyProfile) {
                    data = await getProfile();
                } else {
                    data = await getPublicProfile(id);
                }
                setProfile(data);
            } catch (err) {
                console.error('Failed to load profile', err);
                setError('Не удалось загрузить профиль');
            } finally {
                setLoading(false);
            }
        };

        loadProfile();
    }, [id, isMyProfile]);

    const handleLogout = async () => {
        try {
            await logout();
            navigate('/auth/login');
        } catch (error) {
            console.error('Logout error', error);
        }
    };

    const goTo = (path) => navigate(path);

    const formatLastLogin = () => {
        if (!profile || !profile.last_login_at) {
            return '';
        }
        try {
            return new Date(profile.last_login_at).toLocaleString('ru-RU');
        } catch {
            return profile.last_login_at;
        }
    };

    if (loading) return <div className="container" style={{marginTop: '2rem'}}>Загрузка...</div>;
    if (error) return <div className="container" style={{marginTop: '2rem'}}>{error}</div>;

    return (
        <div className="auth-page">
            <div className="profile-card-wide">
                <div className="profile-main">

                    <div className="profile-avatar-column">
                        <div className="profile-avatar-wrapper">
                            {profile?.avatar_url ? (
                                <img src={profile.avatar_url} alt="Avatar" className="profile-avatar-img" />
                            ) : (
                                <div className="profile-avatar-placeholder">
                                    {(profile?.display_name || profile?.username || '?')[0].toUpperCase()}
                                </div>
                            )}
                        </div>

                        <div className="profile-basic-info">
                            <div className="profile-name">
                                {profile.display_name || profile.username || 'Пользователь'}
                            </div>
                            <div className="profile-location">
                                {[profile.city, profile.country].filter(Boolean).join(', ') || 'Адрес не указан'}
                            </div>
                            {isMyProfile && profile.last_login_at && (
                                <div className="profile-last-login">
                                    Вход: {formatLastLogin()}
                                </div>
                            )}
                        </div>

                        {isMyProfile ? (
                            <div className="profile-actions-vertical">
                                <button className="btn btn-primary profile-btn" onClick={() => goTo(`/account/${user.id}/edit`)}>
                                    Обновить профиль
                                </button>
                                <button className="btn btn-secondary profile-btn" onClick={() => goTo(`/account/${user.id}/change-password`)}>
                                    Сменить пароль
                                </button>
                                <button
                                    className="btn profile-btn"
                                    onClick={handleLogout}
                                    style={{
                                        marginTop: '0.5rem',
                                        backgroundColor: 'white',
                                        border: '1px solid #dc3545',
                                        color: '#dc3545',
                                        fontWeight: '500'
                                    }}
                                >
                                    Выйти
                                </button>
                            </div>
                        ) : (
                            <div className="profile-actions-vertical" style={{marginTop: '1.5rem'}}>
                                <button className="btn btn-primary profile-btn" onClick={() => toast.success('Чат скоро будет!')}>
                                    Написать сообщение
                                </button>
                            </div>
                        )}
                    </div>

                    <div className="profile-details-column">
                        <h2 className="profile-section-title">Информация</h2>

                        <div className="profile-fields-grid">
                            <div className="profile-field">
                                <div className="profile-field-label">О себе</div>
                                <div className="profile-field-value">
                                    {profile.about || (isMyProfile ? 'Расскажите о себе' : 'Нет информации')}
                                </div>
                            </div>
                            {(isMyProfile || profile.phone) && (
                                <div className="profile-field">
                                    <div className="profile-field-label">Телефон</div>
                                    <div className="profile-field-value">
                                        {profile.phone || 'Не указан'}
                                    </div>
                                </div>
                            )}
                        </div>

                        <div className="profile-orders-section" style={{marginTop: '2.5rem', borderTop: '1px solid #eee', paddingTop: '1.5rem'}}>
                            <h2 className="profile-section-title">
                                {isMyProfile ? 'Моя активность' : 'Активность пользователя'}
                            </h2>

                            <div className="orders-buttons-grid" style={{display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '15px'}}>

                                <div
                                    onClick={() => goTo(isMyProfile ? `/account/${user.id}/listings` : `/account/${id}/listings`)}
                                    className="order-card-btn"
                                    style={{
                                        display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
                                        padding: '20px', backgroundColor: '#f8f9fa', border: '1px solid #e9ecef', borderRadius: '12px', cursor: 'pointer'
                                    }}
                                >
                                    <span style={{fontSize: '28px', marginBottom: '8px'}}>📋</span>
                                    <span style={{fontWeight: '600', color: '#333'}}>
                                        {isMyProfile ? 'Мои объявления' : 'Объявления'}
                                    </span>
                                </div>

                                {isMyProfile && (
                                    <div
                                        onClick={() => goTo(`/account/${user.id}/purchases`)}
                                        className="order-card-btn"
                                        style={{
                                            display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
                                            padding: '20px', backgroundColor: '#f8f9fa', border: '1px solid #e9ecef', borderRadius: '12px', cursor: 'pointer'
                                        }}
                                    >
                                        <span style={{fontSize: '28px', marginBottom: '8px'}}>🛍️</span>
                                        <span style={{fontWeight: '600', color: '#333'}}>Мои Покупки</span>
                                    </div>
                                )}

                                {isMyProfile && (
                                    <div
                                        onClick={() => goTo(`/account/${user.id}/sales`)}
                                        className="order-card-btn"
                                        style={{
                                            display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center',
                                            padding: '20px', backgroundColor: '#f8f9fa', border: '1px solid #e9ecef', borderRadius: '12px', cursor: 'pointer'
                                        }}
                                    >
                                        <span style={{fontSize: '28px', marginBottom: '8px'}}>💰</span>
                                        <span style={{fontWeight: '600', color: '#333'}}>Мои Продажи</span>
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default AccountPage;
