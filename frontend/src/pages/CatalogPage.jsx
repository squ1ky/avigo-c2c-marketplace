import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { getCatalogListings } from '../services/api';
import '../styles/catalog.css';

const PAGE_SIZE = 50;

function CatalogPage() {
    const navigate = useNavigate();
    const [listings, setListings] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadListings = async () => {
            try {
                const data = await getCatalogListings({ limit: PAGE_SIZE });
                setListings(data || []);
            } catch (error) {
                console.error('Failed to fetch catalog listings:', error);
                toast.error('Не удалось загрузить каталог');
            } finally {
                setLoading(false);
            }
        };

        loadListings();
    }, []);

    const openListing = (listing) => {
        navigate(`/account/${listing.user_id}/listings/${listing.id}`);
    };

    if (loading) {
        return (
            <div className="catalog-page">
                <div className="catalog-status">Загрузка...</div>
            </div>
        );
    }

    return (
        <div className="catalog-page">
            <div className="catalog-header">
                <h1>Каталог</h1>
                <div className="catalog-count">
                    {listings.length} объявлений
                </div>
            </div>

            {listings.length === 0 ? (
                <div className="catalog-empty">
                    Пока нет активных объявлений
                </div>
            ) : (
                <div className="catalog-grid">
                    {listings.map((listing) => {
                        const imageUrl = listing.media?.[0]?.file_url || '/assets/images/placeholder.jpg';

                        return (
                            <button
                                key={listing.id}
                                type="button"
                                className="catalog-card"
                                onClick={() => openListing(listing)}
                            >
                                <div className="catalog-card-image">
                                    <img src={imageUrl} alt={listing.title} />
                                </div>
                                <div className="catalog-card-body">
                                    <div className="catalog-card-title">{listing.title}</div>
                                    <div className="catalog-card-price">
                                        {Number(listing.price || 0).toLocaleString('ru-RU')} {listing.currency === 'RUB' ? '₽' : listing.currency}
                                    </div>
                                    <div className="catalog-card-meta">
                                        Просмотров: {listing.views_count || 0}
                                    </div>
                                </div>
                            </button>
                        );
                    })}
                </div>
            )}
        </div>
    );
}

export default CatalogPage;
