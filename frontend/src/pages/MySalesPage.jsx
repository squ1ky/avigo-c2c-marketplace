import React, { useEffect, useState } from 'react';
import { getMySales } from '../services/api';
import '../styles/orders.css';
import { toast } from 'react-hot-toast';

const MySalesPage = () => {
    const [orders, setOrders] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchOrders = async () => {
            try {
                const data = await getMySales();
                setOrders(data || []);
            } catch (error) {
                console.error('Failed to fetch sales:', error);
                toast.error('Не удалось загрузить историю продаж');
            } finally {
                setLoading(false);
            }
        };

        fetchOrders();
    }, []);

    const formatDate = (dateString) => {
        try {
            return new Date(dateString).toLocaleDateString('ru-RU', {
                day: 'numeric',
                month: 'long',
                year: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });
        } catch (e) {
            return dateString;
        }
    };

    const getStatusLabel = (status) => {
        const map = {
            'created': 'Ожидает оплаты',
            'paid': 'Оплачен покупателем',
            'confirmed': 'В работе',
            'completed': 'Успешно продан',
            'cancelled': 'Отменен'
        };
        return map[status] || status;
    };

    if (loading) return <div className="orders-page">Загрузка...</div>;

    return (
        <div className="orders-page">
            <h1 className="orders-title">Мои Продажи 💰</h1>

            {orders.length === 0 ? (
                <div className="empty-state">
                    <h2>У вас пока нет продаж</h2>
                    <p>Разместите больше объявлений, чтобы найти покупателей!</p>
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
                                        + {order.amount ? order.amount.toLocaleString('ru-RU') : 0} {order.currency || 'RUB'}
                                    </span>
                                    <span className="order-date">{formatDate(order.created_at)}</span>
                                </div>

                                <div className="counterparty-info">
                                    <span className="counterparty-label">Покупатель:</span>
                                    <img
                                        src={order.counterparty_img || 'https://via.placeholder.com/32'}
                                        alt="Buyer"
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

export default MySalesPage;
