import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getUserPurchases } from '../services/api';
import '../styles/orders.css';
import { toast } from 'react-hot-toast';

const UserPurchasesPage = () => {
    const { userId: id } = useParams();
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
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

        fetchOrders();
    }, [id]);

    const formatDate = (dateString) => {
        try {
            return new Date(dateString).toLocaleDateString('ru-RU', {
                day: 'numeric', month: 'long', year: 'numeric',
                hour: '2-digit', minute: '2-digit'
            });
        } catch (e) {
            return dateString;
        }
    };

    const getStatusLabel = (status) => {
        const map = {
            'created': 'Создан',
            'paid': 'Оплачен',
            'confirmed': 'Подтвержден',
            'completed': 'Завершен',
            'cancelled': 'Отменен'
        };
        return map[status] || status;
    };

    if (loading) return <div className="orders-page">Загрузка...</div>;

    return (
        <div className="orders-page">
            <h1 className="orders-title">Мои Покупки 🛍️</h1>

            {orders.length === 0 ? (
                <div className="empty-state">
                    <h2>У вас пока нет покупок</h2>
                    <p>Самое время найти что-нибудь интересное!</p>
                </div>
            ) : (
                <div className="orders-list">
                    {orders.map((order) => (
                        <div key={order.id} className="order-card">
                            <div className="order-header">
                                <span className="order-id">Заказ #{order.id.slice(0, 8)}</span>
                                <span className={`status-badge status-${order.status}`}>
                                    {getStatusLabel(order.status)}
                                </span>
                            </div>

                            <div className="order-body">
                                <div className="order-info">
                                    <h3>Товар {order.listing_title || 'из объявления'}</h3>
                                    <span className="order-price">
                                        {order.amount ? order.amount.toLocaleString('ru-RU') : 0} {order.currency || 'RUB'}
                                    </span>
                                    <span className="order-date">{formatDate(order.created_at)}</span>
                                </div>

                                <div className="counterparty-info">
                                    <span className="counterparty-label">Продавец:</span>
                                    <img
                                        src={order.counterparty_img || 'https://via.placeholder.com/32'}
                                        alt="Seller"
                                        className="counterparty-avatar"
                                    />
                                    <span className="counterparty-name">
                                        {order.counterparty_name || 'Неизвестный'}
                                    </span>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default UserPurchasesPage;
