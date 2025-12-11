import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { getListing } from '../services/api';
import toast from 'react-hot-toast';

function ListingPage() {

    const { id } = useParams();
    const [listing, setListing] = useState(null);
    const [loading, setLoading] = useState(true);
    const [activeImageIndex, setActiveImageIndex] = useState(0);

    useEffect(() => {
        const fetchListing = async () => {
            try {
                const data = await getListing(id);
                setListing(data);
            } catch (error) {
                toast.error('Не удалось загрузить объявление');
            } finally {
                setLoading(false);
            }
        };
        fetchListing();
    }, [id]);

    // Хелпер для формирования URL
    const getImageUrl = (url) => {
        if (!url) return '/assets/images/placeholder.jpg';
        if (url.startsWith('http')) return url;
        return S3_BASE + url;
    };

    if (loading) return <div className="container">Загрузка...</div>;
    if (!listing) return <div className="container">Объявление не найдено</div>;

    // Вычисляем URL для главной картинки
    const mainImage = listing.media && listing.media.length > 0
        ? getImageUrl(listing.media[activeImageIndex].file_url)
        : '/assets/images/placeholder.jpg';

    return (
        <div className="container" style={{ marginTop: '2rem' }}>
            <div className="listing-layout" style={{ display: 'grid', gridTemplateColumns: '1.5fr 1fr', gap: '2rem' }}>

                <div className="listing-gallery">
                    <div className="main-image-container" style={{
                        width: '100%',
                        height: '400px',
                        backgroundColor: '#f5f5f5',
                        borderRadius: '12px',
                        overflow: 'hidden',
                        marginBottom: '1rem',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center'
                    }}>
                        <img
                            src={mainImage}
                            alt={listing.title}
                            style={{ maxWidth: '100%', maxHeight: '100%', objectFit: 'contain' }}
                        />
                    </div>

                    {listing.media && listing.media.length > 1 && (
                        <div className="thumbnails" style={{ display: 'flex', gap: '0.5rem', overflowX: 'auto' }}>
                            {listing.media.map((item, index) => (
                                <div
                                    key={item.id}
                                    onClick={() => setActiveImageIndex(index)}
                                    style={{
                                        width: '80px',
                                        height: '60px',
                                        cursor: 'pointer',
                                        border: index === activeImageIndex ? '2px solid var(--go-cyan)' : '1px solid #ddd',
                                        borderRadius: '4px',
                                        overflow: 'hidden',
                                        flexShrink: 0 // Чтобы миниатюры не сжимались
                                    }}
                                >
                                    <img
                                        src={getImageUrl(item.file_url)} // Хак здесь тоже
                                        alt=""
                                        style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                                    />
                                </div>
                            ))}
                        </div>
                    )}

                    <div className="listing-description" style={{ marginTop: '2rem' }}>
                        <h3>Описание</h3>
                        <p style={{ whiteSpace: 'pre-wrap', marginTop: '0.5rem', color: '#333' }}>
                            {listing.description}
                        </p>
                    </div>
                </div>

                <div className="listing-sidebar">
                    <div className="listing-card" style={{
                        padding: '1.5rem',
                        border: '1px solid #eee',
                        borderRadius: '12px',
                        boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
                        position: 'sticky',
                        top: '100px'
                    }}>
                        <h1 style={{ fontSize: '1.8rem', marginBottom: '0.5rem' }}>{listing.title}</h1>

                        <div className="listing-price" style={{ fontSize: '2rem', fontWeight: '800', marginBottom: '1.5rem' }}>
                            {listing.price.toLocaleString()} {listing.currency === 'RUB' ? '₽' : listing.currency}
                        </div>

                        <button className="btn btn-primary btn-large" style={{ width: '100%', marginBottom: '1rem' }}>
                            Купить с доставкой
                        </button>

                        <button className="btn btn-secondary btn-large" style={{ width: '100%' }}>
                            Написать продавцу
                        </button>

                        <div className="seller-info" style={{ marginTop: '2rem', paddingTop: '1rem', borderTop: '1px solid #eee' }}>
                            <p style={{ color: '#666', fontSize: '0.9rem' }}>Продавец</p>
                            <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginTop: '0.5rem' }}>
                                <div style={{ width: '40px', height: '40px', background: '#ddd', borderRadius: '50%' }}></div>
                                <span style={{ fontWeight: '600' }}>Пользователь {listing.user_id ? listing.user_id.slice(0, 8) : '...'}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default ListingPage;
