import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-hot-toast';
import { getMyListings, deleteListing } from '../services/api';
import '../styles/orders.css';

const MyListingsPage = () => {
    const [listings, setListings] = useState([]);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    const fetchMyListings = async () => {
        try {
            const data = await getMyListings();
            setListings(data || []);
        } catch (error) {
            console.error('Failed to fetch listings:', error);
            toast.error('Не удалось загрузить ваши объявления');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchMyListings();
    }, []);

    const handleDelete = async (id, e) => {
        e.stopPropagation();
        if (!window.confirm('Вы действительно хотите удалить это объявление?')) {
            return;
        }

        try {
            await deleteListing(id);
            toast.success('Объявление удалено');
            setListings(prev => prev.filter(item => item.id !== id));
        } catch (error) {
            console.error('Delete error:', error);
            toast.error('Ошибка при удалении');
        }
    };

    const handleEdit = (id, e) => {
        e.stopPropagation();
        navigate(`/listing/edit/${id}`);
    };

    const formatDate = (dateString) => {
        try {
            return new Date(dateString).toLocaleDateString('ru-RU', {
                day: 'numeric', month: 'long', year: 'numeric'
            });
        } catch {
            return dateString;
        }
    };

    if (loading) return <div className="orders-page">Загрузка...</div>;

    return (
        <div className="orders-page">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '30px' }}>
                <h1 className="orders-title" style={{ margin: 0 }}>Мои Объявления 📋</h1>
                <button
                    className="btn btn-primary"
                    onClick={() => navigate('/create')}
                    style={{ padding: '10px 20px', borderRadius: '8px' }}
                >
                    + Создать
                </button>
            </div>

            {listings.length === 0 ? (
                <div className="empty-state">
                    <h2>У вас нет активных объявлений</h2>
                    <p>Самое время что-нибудь продать!</p>
                    <button
                        className="btn btn-primary"
                        onClick={() => navigate('/create')}
                        style={{ marginTop: '20px' }}
                    >
                        Разместить объявление
                    </button>
                </div>
            ) : (
                <div className="orders-list">
                    {listings.map((item) => (
                        <div
                            key={item.id}
                            className="order-card"
                            onClick={() => navigate(`/listing/${item.id}`)}
                            style={{ cursor: 'pointer' }}
                        >
                            <div className="order-header">
                                <span className="order-id">
                                    {formatDate(item.created_at)}
                                </span>
                                <span className="status-badge status-active" style={{ background: '#e3f2fd', color: '#0d47a1' }}>
                                    Активно
                                </span>
                            </div>

                            <div className="order-body">
                                <div style={{ display: 'flex', gap: '15px', alignItems: 'center' }}>
                                    <img
                                        src={
                                            (item.media && item.media.length > 0)
                                                ? item.media[0].file_url
                                                : '/assets/images/placeholder.png'
                                        }
                                        alt={item.title}
                                        style={{
                                            width: '80px',
                                            height: '80px',
                                            objectFit: 'cover',
                                            borderRadius: '8px',
                                            backgroundColor: '#f0f0f0'
                                        }}
                                    />
                                    <div className="order-info">
                                        <h3 style={{ margin: '0 0 5px 0', fontSize: '1.1rem' }}>{item.title}</h3>
                                        <span className="order-price">
                                            {item.price ? item.price.toLocaleString('ru-RU') : 0} {item.currency || 'RUB'}
                                        </span>
                                        <span style={{ fontSize: '0.85rem', color: '#777' }}>
                                            👀 Просмотров: {item.views_count || 0}
                                        </span>
                                    </div>
                                </div>

                                <div className="listing-actions" style={{ display: 'flex', gap: '10px' }}>
                                    <button
                                        onClick={(e) => handleEdit(item.id, e)}
                                        className="btn btn-secondary"
                                        style={{ padding: '8px 15px', fontSize: '0.9rem' }}
                                    >
                                        Изменить
                                    </button>
                                    <button
                                        onClick={(e) => handleDelete(item.id, e)}
                                        className="btn btn-danger"
                                        style={{
                                            padding: '8px 15px',
                                            fontSize: '0.9rem',
                                            background: '#ffebee',
                                            color: '#c62828',
                                            border: '1px solid #ffcdd2'
                                        }}
                                    >
                                        Удалить
                                    </button>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default MyListingsPage;
