import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { getListing, updateListing } from '../services/api';
import ImageUploader from '../components/common/ImageUploader';
import CategorySelect from '../components/common/CategorySelect';
import toast from 'react-hot-toast';
import '../styles/main.css';
import {useAuth} from "../context/AuthContext.jsx";

function EditListingPage() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const { user } = useAuth();

    const [formData, setFormData] = useState({
        title: '',
        description: '',
        price: '',
        currency: 'RUB',
        category_id: '',
        media_ids: []
    });

    useEffect(() => {
        const fetchListing = async () => {
            try {
                const response = await getListing(id);
                const listing = response.listing;

                setFormData({
                    title: listing.title,
                    description: listing.description,
                    price: listing.price,
                    currency: listing.currency,
                    category_id: listing.category_id,
                    media_ids: listing.media ? listing.media.map(m => m.id) : []
                });
            } catch (error) {
                console.error("Failed to fetch listing", error);
                toast.error("Не удалось загрузить данные объявления");
                navigate(`/account/${user?.id}/listings`);
            } finally {
                setLoading(false);
            }
        };

        if (id) {
            fetchListing();
        }
    }, [id, navigate]);

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handleCategoryChange = (val) => {
        setFormData(prev => ({ ...prev, category_id: val }));
    };

    const setMediaIds = (ids) => {
        setFormData(prev => ({ ...prev, media_ids: ids }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!formData.category_id) {
            toast.error('Пожалуйста, выберите категорию');
            return;
        }

        setSaving(true);

        const payload = {
            ...formData,
            price: parseFloat(formData.price),
        };

        try {
            await updateListing(id, payload);
            toast.success('Объявление обновлено!');
            navigate(`/listings/${id}`);
        } catch (error) {
            console.error(error);
            toast.error(error.message || 'Ошибка обновления');
        } finally {
            setSaving(false);
        }
    };

    if (loading) {
        return <div className="container" style={{marginTop: '2rem', textAlign: 'center'}}>Загрузка...</div>;
    }

    return (
        <div className="container" style={{ marginTop: '2rem', maxWidth: '800px' }}>
            <h1 className="section-title">Редактирование объявления</h1>

            <form onSubmit={handleSubmit} className="auth-form" style={{ maxWidth: '100%' }}>

                <div className="form-group">
                    <label className="form-label">Категория <span style={{color: 'red'}}>*</span></label>
                    <CategorySelect
                        value={formData.category_id}
                        onChange={handleCategoryChange}
                    />
                </div>

                <div className="form-group">
                    <label className="form-label">Название</label>
                    <input
                        type="text"
                        name="title"
                        className="form-input"
                        placeholder="Например, iPhone 15 Pro, 256GB"
                        value={formData.title}
                        onChange={handleChange}
                        required
                        minLength={5}
                    />
                </div>

                <div className="form-group">
                    <label className="form-label">Фотографии</label>
                    <div style={{marginBottom: '10px', fontSize: '0.9em', color: '#666'}}>
                        Загрузите новые фото, которые будут добавлены к старым.
                    </div>
                    <ImageUploader
                        mediaIds={formData.media_ids}
                        setMediaIds={setMediaIds}
                    />
                </div>

                <div className="form-group">
                    <label className="form-label">Описание</label>
                    <textarea
                        name="description"
                        className="form-input"
                        placeholder="Расскажите о товаре: состояние, дефекты, причина продажи..."
                        value={formData.description}
                        onChange={handleChange}
                        required
                        minLength={10}
                        rows={6}
                        style={{ resize: 'vertical' }}
                    />
                </div>

                <div className="form-group">
                    <label className="form-label">Цена</label>
                    <div style={{ display: 'flex', gap: '1rem' }}>
                        <input
                            type="number"
                            name="price"
                            className="form-input"
                            placeholder="0"
                            value={formData.price}
                            onChange={handleChange}
                            required
                            min={0}
                        />
                        <select
                            name="currency"
                            className="form-input"
                            style={{ width: '100px' }}
                            value={formData.currency}
                            onChange={handleChange}
                        >
                            <option value="RUB">₽</option>
                            <option value="USD">$</option>
                            <option value="EUR">€</option>
                        </select>
                    </div>
                </div>

                <div style={{display: 'flex', gap: '1rem'}}>
                    <button type="submit" className="btn btn-primary btn-large" disabled={saving} style={{flex: 1}}>
                        {saving ? 'Сохранение...' : 'Сохранить изменения'}
                    </button>
                    <button
                        type="button"
                        className="btn btn-secondary btn-large"
                        onClick={() => navigate(-1)}
                        disabled={saving}
                    >
                        Отмена
                    </button>
                </div>
            </form>
        </div>
    );
}

export default EditListingPage;
