import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getUserPurchases, completeOrder, cancelOrder } from '../services/api';
import '../styles/orders.css';
import { toast } from 'react-hot-toast';

const UserPurchasesPage = () => {
    const { userId: id } = useParams();
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);

    const fetchOrders = async () => {
        if (!id) return;
        try {
            const data = await getUserPurchases(id);
            setOrders(data || []);
        } catch (error) {
            console.error('Failed to fetch purchases:', error);
            toast.error('Не удалось загрузить историю покупок');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchOrders();
    }, [id]);

    const handleComplete = async (orderId) => {
        try {
            await completeOrder(orderId);
            toast.success('Заказ завершен!');
            fetchOrders();
        } catch (error) {
            toast.error('Ошибка при завершении заказа');
        }
    };

    const handleCancel = async (orderId) => {
        if (!window.confirm('Вы уверены, что хотите отменить заказ?')) return;
        try {
            await cancelOrder(orderId, "Cancelled by buyer");
            toast.success('Заказ отменен');
            fetchOrders();
        } catch (error) {
            toast.error('Ошибка при отмене заказа');
        }
    };

    const formatDate = (dateString) => {
        try { return new Date(dateString).toLocaleDateString('ru-RU', { day: 'numeric', month: 'long', year: 'numeric' }); } catch { return dateString; }
    };

    const getStatusLabel = (status) => {
        const map = {
            'created': 'Ожидает подтверждения',
            'confirmed': 'Отправлен / В пути',
            'completed': 'Получен',
            'cancelled': 'Отменен'
        };
        return map[status] || status;
    };

    if (loading) return <div className="orders-page"><div className="container">Загрузка...</div></div>;

    return (
        <div className="orders-page container" style={{ marginTop: '2rem' }}>
            <h1 className="orders-title" style={{ marginBottom: '1.5rem' }}>Мои Покупки 🛍️</h1>

            {orders.length === 0 ? (
                <div className="empty-state">
                    <h2>У вас пока нет покупок</h2>
                </div>
            ) : (
                <div className="orders-list" style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    {orders.map((order) => {
                        const { listing, counterparty } = order;
                        const price = listing?.price || 0;
                        const currency = listing?.currency === 'RUB' ? '₽' : listing?.currency;
                        const mainImage = listing?.media?.[0]?.file_url || '/assets/images/placeholder.jpg';
                        const listingLink = `/listings/${listing?.id}`;
                        const userLink = `/account/${counterparty?.id}`;

                        return (
                            <div key={order.id} className="order-card" style={{ border: '1px solid #eee', padding: '1.5rem', borderRadius: '12px', background: '#fff' }}>
                                <div className="order-header" style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1rem' }}>
                                    <span className="order-id" style={{ color: '#888' }}>Заказ #{order.id.slice(0, 8)}</span>
                                    <span className={`status-badge status-${order.status}`} style={{ fontWeight: 'bold', color: order.status === 'cancelled' ? 'red' : order.status === 'completed' ? 'green' : '#f39c12' }}>
                                        {getStatusLabel(order.status)}
                                    </span>
                                </div>

                                <div className="order-body" style={{ display: 'flex', gap: '1.5rem' }}>
                                    <div style={{ width: '100px', height: '100px', flexShrink: 0, borderRadius: '8px', overflow: 'hidden', border: '1px solid #eee' }}>
                                        <img src={mainImage} alt="" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                                    </div>

                                    <div style={{ flex: 1 }}>
                                        <h3 style={{ marginTop: 0, marginBottom: '0.5rem' }}>
                                            <Link to={listingLink} style={{ textDecoration: 'none', color: '#333' }}>
                                                {listing?.title || 'Товар недоступен'}
                                            </Link>
                                        </h3>
                                        <div style={{ fontSize: '1.2rem', fontWeight: 'bold', marginBottom: '0.5rem' }}>
                                            {price.toLocaleString('ru-RU')} {currency}
                                        </div>
                                        <div style={{ color: '#666', fontSize: '0.9rem' }}>
                                            Дата: {formatDate(order.created_at)}
                                        </div>

                                        <div className="order-actions" style={{ marginTop: '1rem', display: 'flex', gap: '10px' }}>
                                            {order.status === 'confirmed' && (
                                                <button className="btn btn-success"
                                                        style={{ padding: '8px 16px', background: '#27ae60', color: 'white', border: 'none', borderRadius: '6px', cursor: 'pointer' }}
                                                        onClick={() => handleComplete(order.id)}>
                                                    Я получил товар
                                                </button>
                                            )}
                                            {(order.status === 'created' || order.status === 'confirmed') && (
                                                <button className="btn btn-danger"
                                                        style={{ padding: '8px 16px', background: 'transparent', border: '1px solid #e74c3c', color: '#e74c3c', borderRadius: '6px', cursor: 'pointer' }}
                                                        onClick={() => handleCancel(order.id)}>
                                                    Отменить
                                                </button>
                                            )}
                                        </div>
                                    </div>

                                    <div style={{ minWidth: '150px', borderLeft: '1px solid #eee', paddingLeft: '1.5rem', display: 'flex', flexDirection: 'column', justifyContent: 'center' }}>
                                        <div style={{ fontSize: '0.9rem', color: '#666', marginBottom: '0.5rem' }}>Продавец</div>
                                        <Link to={userLink} style={{ display: 'flex', alignItems: 'center', gap: '10px', textDecoration: 'none', color: 'inherit' }}>
                                            <div style={{
                                                width: '32px',
                                                height: '32px',
                                                borderRadius: '50%',
                                                background: '#eee',
                                                overflow: 'hidden',
                                                flexShrink: 0
                                            }}>
                                                {counterparty?.profile?.avatar_url ? (
                                                    <img src={counterparty.profile.avatar_url} alt=""
                                                         style={{width: '100%', height: '100%', objectFit: 'cover'}}/>
                                                ) : (
                                                    <div style={{
                                                        width: '100%',
                                                        height: '100%',
                                                        display: 'flex',
                                                        alignItems: 'center',
                                                        justifyContent: 'center'
                                                    }}>
                                                        {(counterparty?.display_name || '?')[0]}
                                                    </div>
                                                )}
                                            </div>
                                            <span
                                                style={{fontWeight: '500'}}>{counterparty?.display_name || 'User'}</span>
                                        </Link>
                                    </div>
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
};

export default UserPurchasesPage;
