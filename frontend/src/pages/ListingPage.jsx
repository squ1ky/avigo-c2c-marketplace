import { useEffect, useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { getListing, createOrder } from '../services/api';
import { useAuth } from '../context/AuthContext';
import toast from 'react-hot-toast';

function ListingPage() {
    const { userId, listingId } = useParams();
    const { user } = useAuth();
    const navigate = useNavigate();

    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [isBuying, setIsBuying] = useState(false);
    const [activeImageIndex, setActiveImageIndex] = useState(0);
    const [showPhone, setShowPhone] = useState(false);

    useEffect(() => {
        const fetchListing = async () => {
            if (!userId || !listingId) return;
            try {
                const response = await getListing(userId, listingId);
                setData(response);
            } catch (error) {
                console.error(error);
                toast.error('Не удалось загрузить объявление');
            } finally {
                setLoading(false);
            }
        };
        fetchListing();
    }, [userId, listingId]);

    const handleBuy = async () => {
        if (!user) {
            toast.error('Войдите, чтобы оформить заказ');
            navigate('/auth/login');
            return;
        }

        if (data.listing.user_id === user.id) {
            toast.error('Вы не можете купить свой собственный товар');
            return;
        }

        try {
            setIsBuying(true);
            await createOrder(data.listing.id);
            toast.success('Заказ успешно создан!');
            navigate(`/account/${user.id}/purchases`);
        } catch (error) {
            console.error(error);
            toast.error(error.message || 'Не удалось создать заказ');
        } finally {
            setIsBuying(false);
        }
    };

    const getImageUrl = (url) => {
        if (!url) return '/assets/images/placeholder.jpg';
        return url;
    };

    if (loading) return <div className="container" style={{ marginTop: '2rem' }}>Загрузка...</div>;
    if (!data) return <div className="container" style={{ marginTop: '2rem' }}>Объявление не найдено</div>;

    const { listing, user: seller } = data;
    const mainImage = listing.media && listing.media.length > 0
        ? getImageUrl(listing.media[activeImageIndex].file_url)
        : '/assets/images/placeholder.jpg';

    const phoneRaw = seller.profile?.phone;
    const hasPhone = !!phoneRaw;
    const canBuy = listing.status === 'active' && !listing.is_sold;

    return (
        <div className="container" style={{ marginTop: '2rem' }}>
            <div className="listing-layout" style={{ display: 'grid', gridTemplateColumns: '1.5fr 1fr', gap: '2rem' }}>
                <div className="listing-gallery">
                    <div className="main-image-container" style={{
                        width: '100%', height: '400px', backgroundColor: '#f5f5f5',
                        borderRadius: '12px', overflow: 'hidden', marginBottom: '1rem',
                        display: 'flex', alignItems: 'center', justifyContent: 'center'
                    }}>
                        <img src={mainImage} alt={listing.title} style={{ maxWidth: '100%', maxHeight: '100%', objectFit: 'contain' }} />
                    </div>
                    {listing.media && listing.media.length > 1 && (
                        <div className="thumbnails" style={{ display: 'flex', gap: '0.5rem', overflowX: 'auto' }}>
                            {listing.media.map((item, index) => (
                                <div key={item.id} onClick={() => setActiveImageIndex(index)}
                                     style={{
                                         width: '80px', height: '60px', cursor: 'pointer',
                                         border: index === activeImageIndex ? '2px solid var(--go-cyan)' : '1px solid #ddd',
                                         borderRadius: '4px', overflow: 'hidden', flexShrink: 0
                                     }}
                                >
                                    <img src={getImageUrl(item.file_url)} alt="" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                                </div>
                            ))}
                        </div>
                    )}
                    <div className="listing-description" style={{ marginTop: '2rem' }}>
                        <h3>Описание</h3>
                        <p style={{ whiteSpace: 'pre-wrap', marginTop: '0.5rem', color: '#333' }}>{listing.description}</p>
                    </div>
                </div>

                <div className="listing-sidebar">
                    <div className="listing-card" style={{
                        padding: '1.5rem', border: '1px solid #eee', borderRadius: '12px',
                        boxShadow: '0 4px 12px rgba(0,0,0,0.05)', position: 'sticky', top: '100px'
                    }}>
                        <h1 style={{ fontSize: '1.8rem', marginBottom: '0.5rem' }}>{listing.title}</h1>
                        <div className="listing-price" style={{ fontSize: '2rem', fontWeight: '800', marginBottom: '1.5rem' }}>
                            {listing.price.toLocaleString()} {listing.currency === 'RUB' ? '₽' : listing.currency}
                        </div>

                        {hasPhone && (
                            !showPhone ? (
                                <button className="btn btn-primary btn-large"
                                        style={{ width: '100%', marginBottom: '0.5rem', backgroundColor: '#2ecc71', borderColor: '#27ae60' }}
                                        onClick={() => setShowPhone(true)}
                                >
                                    Показать телефон
                                </button>
                            ) : (
                                <a href={`tel:${phoneRaw}`} className="btn btn-large"
                                   style={{
                                       width: '100%', marginBottom: '0.5rem', backgroundColor: '#fff',
                                       border: '2px solid #2ecc71', color: '#2ecc71', fontSize: '1.2rem',
                                       fontWeight: 'bold', display: 'flex', justifyContent: 'center', alignItems: 'center', textDecoration: 'none'
                                   }}
                                >
                                    {phoneRaw}
                                </a>
                            )
                        )}

                        <button className="btn btn-secondary btn-large" style={{ width: '100%', marginBottom: '1rem' }}>Написать сообщение</button>
                        <button className="btn btn-primary btn-large" style={{ width: '100%', opacity: canBuy ? 1 : 0.6 }}
                                onClick={handleBuy} disabled={isBuying || !canBuy}
                        >
                            {isBuying ? 'Обработка...' : (listing.is_sold ? 'Товар продан' : 'Купить с доставкой')}
                        </button>

                        <div className="seller-info" style={{ marginTop: '2rem', paddingTop: '1rem', borderTop: '1px solid #eee' }}>
                            <p style={{ color: '#666', fontSize: '0.9rem' }}>Продавец</p>
                            <Link to={`/account/${seller.id}`} style={{ textDecoration: 'none', color: 'inherit' }}>
                                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginTop: '0.5rem' }}>
                                    <div style={{ width: '40px', height: '40px', background: '#ddd', borderRadius: '50%', overflow: 'hidden' }}>
                                        {seller.profile?.avatar_url ? (
                                            <img src={seller.profile.avatar_url} alt="" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                                        ) : (
                                            <div style={{ width: '100%', height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                                                {(seller.display_name || seller.username || '?')[0].toUpperCase()}
                                            </div>
                                        )}
                                    </div>
                                    <div>
                                        <div style={{ fontWeight: '600' }}>{seller.display_name || seller.username}</div>
                                        <div style={{ fontSize: '0.8rem', color: '#888' }}>
                                            {[seller.profile?.city, seller.profile?.country].filter(Boolean).join(', ')}
                                        </div>
                                    </div>
                                </div>
                            </Link>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default ListingPage;
