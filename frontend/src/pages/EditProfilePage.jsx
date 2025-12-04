import { useEffect, useState } from 'react';
import { getProfile, updateProfile } from '../services/api';
import { countries, citiesByCountry } from '../utils/countries';
import toast from 'react-hot-toast';
import '../styles/auth.css';

function EditProfilePage() {
    const [form, setForm] = useState({
        display_name: '',
        phone: '',
        about: '',
        country: '',
        city: '',
        avatar_url: '',
    });
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [error, setError] = useState(null);

    useEffect(() => {
        const load = async () => {
            try {
                const profile = await getProfile();
                setForm({
                    display_name: profile.display_name || '',
                    phone: profile.phone || '',
                    about: profile.about || '',
                    country: profile.country || '',
                    city: profile.city || '',
                    avatar_url: profile.avatar_url || '',
                });
            } catch (err) {
                console.error(err);
                setError(err.message || 'Ошибка загрузки профиля');
            } finally {
                setLoading(false);
            }
        };
        load();
    }, []);

    const handleChange = (e) => {
        const { name, value } = e.target;

        if (name === 'country') {
            // как в RegisterForm: при смене страны сбрасываем город
            setForm((prev) => ({ ...prev, country: value, city: '' }));
        } else {
            setForm((prev) => ({ ...prev, [name]: value }));
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setSaving(true);
        setError(null);
        try {
            const payload = {
                display_name: form.display_name || undefined,
                phone: form.phone || undefined,
                about: form.about || undefined,
                country: form.country || undefined,
                city: form.city || undefined,
                avatar_url: form.avatar_url || undefined,
            };

            await updateProfile(payload);
            toast.success('Профиль обновлён');
        } catch (err) {
            console.error(err);
            setError(err.message || 'Не удалось сохранить профиль');
            toast.error(err.message || 'Не удалось сохранить профиль');
        } finally {
            setSaving(false);
        }
    };

    const availableCities = form.country
        ? citiesByCountry[form.country] || []
        : [];

    if (loading) {
        return (
            <div className="auth-page">
                <div className="auth-card">
                    <div className="auth-header">
                        <h1 className="auth-title">Редактирование профиля</h1>
                        <p className="auth-subtitle">Загрузка данных...</p>
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="auth-page">
            <div className="auth-card">
                <div className="auth-header">
                    <div className="auth-logo">
                        <span className="logo-avi">Avi</span>
                        <span className="logo-go">Go</span>
                    </div>
                    <h1 className="auth-title">Редактирование профиля</h1>
                    <p className="auth-subtitle">
                        Обновите личную информацию
                    </p>
                </div>

                <form className="auth-form" onSubmit={handleSubmit}>
                    {error && <div className="form-error">{error}</div>}

                    <div className="form-group">
                        <label className="form-label">
                            Отображаемое имя <span className="required">*</span>
                        </label>
                        <input
                            type="text"
                            name="display_name"
                            className="form-input"
                            value={form.display_name}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label className="form-label">Телефон (опционально)</label>
                        <input
                            type="tel"
                            name="phone"
                            className="form-input"
                            value={form.phone}
                            onChange={handleChange}
                            placeholder="+7 999 123 45 67"
                        />
                    </div>

                    <div className="form-group">
                        <label className="form-label">О себе</label>
                        <textarea
                            name="about"
                            className="form-input"
                            rows="3"
                            value={form.about}
                            onChange={handleChange}
                        />
                    </div>

                    <div className="form-row">
                        <div className="form-group">
                            <label className="form-label">Страна (опционально)</label>
                            <select
                                name="country"
                                value={form.country}
                                onChange={handleChange}
                                className="form-select"
                            >
                                {countries.map((country) => (
                                    <option
                                        key={country.value}
                                        value={country.value}
                                    >
                                        {country.label}
                                    </option>
                                ))}
                            </select>
                        </div>

                        <div className="form-group">
                            <label className="form-label">Город (опционально)</label>
                            <select
                                name="city"
                                value={form.city}
                                onChange={handleChange}
                                className="form-select"
                                disabled={!form.country}
                            >
                                {availableCities.length > 0 ? (
                                    availableCities.map((city) => (
                                        <option
                                            key={city.value}
                                            value={city.value}
                                        >
                                            {city.label}
                                        </option>
                                    ))
                                ) : (
                                    <option value="">
                                        Сначала выберите страну
                                    </option>
                                )}
                            </select>
                        </div>
                    </div>

                    <div className="form-group">
                        <label className="form-label">URL аватара (опционально)</label>
                        <input
                            type="url"
                            name="avatar_url"
                            className="form-input"
                            value={form.avatar_url}
                            onChange={handleChange}
                            placeholder="https://example.com/avatar.jpg"
                        />
                    </div>

                    <button
                        type="submit"
                        className="auth-submit"
                        disabled={saving}
                    >
                        {saving ? 'Сохранение...' : 'Сохранить изменения'}
                    </button>
                </form>
            </div>
        </div>
    );
}

export default EditProfilePage;
