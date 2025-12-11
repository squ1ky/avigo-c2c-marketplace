import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createListing } from '../services/api';
import ImageUploader from '../components/common/ImageUploader';
import CategorySelect from '../components/common/CategorySelect';
import toast from 'react-hot-toast';
import '../styles/main.css';

function CreateListingPage() {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);

    const [formData, setFormData] = useState({
        title: '',
        description: '',
        price: '',
        currency: 'RUB',
        category_id: '',
        media_ids: []
    });

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

        setLoading(true);

        const payload = {
            ...formData,
            price: parseFloat(formData.price),
        };

        try {
            const created = await createListing(payload);
            toast.success('Объявление опубликовано!');
            navigate(`/listing/${created.id}`);
        } catch (error) {
            console.error(error);
            toast.error(error.message || 'Ошибка создания');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="container" style={{ marginTop: '2rem', maxWidth: '800px' }}>
            <h1 className="section-title">Новое объявление</h1>

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

                <button type="submit" className="btn btn-primary btn-large" disabled={loading}>
                    {loading ? 'Публикация...' : 'Опубликовать'}
                </button>
            </form>
        </div>
    );
}

export default CreateListingPage;
