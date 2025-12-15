import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { parsePhoneNumber } from 'libphonenumber-js';
import { register } from '../../services/api';
import {
    validateEmail,
    validateUsername,
    validateDisplayName,
    validatePassword,
    validatePasswordMatch,
    getPasswordStrength
} from '../../utils/validators';
import { countries, citiesByCountry } from '../../utils/countries';
import toast from 'react-hot-toast';

function RegisterForm() {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
        confirmPassword: '',
        display_name: '',
        phone: '',
        country: '',
        city: ''
    });
    const [errors, setErrors] = useState({});

    const handleChange = (e) => {
        const { name, value } = e.target;

        if (name === 'country') {
            setFormData(prev => ({ ...prev, [name]: value, city: '' }));
        } else {
            setFormData(prev => ({ ...prev, [name]: value }));
        }

        if (errors[name]) {
            setErrors(prev => ({ ...prev, [name]: '' }));
        }
    };

    const handlePhoneBlur = () => {
        if (formData.phone) {
            try {
                const phoneNumber = parsePhoneNumber(formData.phone, 'RU');
                if (phoneNumber && phoneNumber.isValid()) {
                    setFormData(prev => ({
                        ...prev,
                        phone: phoneNumber.formatInternational()
                    }));
                }
            } catch (error) {
                console.log('Phone parsing error:', error);
            }
        }
    };

    const validateForm = () => {
        const newErrors = {};

        const usernameError = validateUsername(formData.username);
        if (usernameError) newErrors.username = usernameError;

        const emailError = validateEmail(formData.email);
        if (emailError) newErrors.email = emailError;

        const displayNameError = validateDisplayName(formData.display_name);
        if (displayNameError) newErrors.display_name = displayNameError;

        const passwordError = validatePassword(formData.password);
        if (passwordError) newErrors.password = passwordError;

        const passwordMatchError = validatePasswordMatch(formData.password, formData.confirmPassword);
        if (passwordMatchError) newErrors.confirmPassword = passwordMatchError;

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!validateForm()) {
            toast.error('Пожалуйста, исправьте ошибки в форме');
            return;
        }

        setLoading(true);

        try {
            const { confirmPassword, ...dataToSend } = formData;

            if (dataToSend.phone) {
                try {
                    const phoneNumber = parsePhoneNumber(dataToSend.phone, 'RU');
                    if (phoneNumber && phoneNumber.isValid()) {
                        dataToSend.phone = phoneNumber.format('E.164');
                    }
                } catch (error) {
                    console.log('Phone formatting error:', error);
                }
            }

            if (!dataToSend.phone) delete dataToSend.phone;
            if (!dataToSend.country) delete dataToSend.country;
            if (!dataToSend.city) delete dataToSend.city;

            const response = await register(dataToSend);

            toast.success('Регистрация успешна! Проверьте email для подтверждения');

            localStorage.setItem('pendingUserId', response.user_id);
            localStorage.setItem('pendingUserEmail', response.email);

            navigate('/auth/confirm-email');
        } catch (error) {
            toast.error(error.message || 'Ошибка регистрации');
        } finally {
            setLoading(false);
        }
    };

    const passwordStrength = getPasswordStrength(formData.password);
    const availableCities = formData.country ? citiesByCountry[formData.country] || [] : [];

    return (
        <form className="auth-form" onSubmit={handleSubmit}>
            <div className="form-row">
                <div className="form-group">
                    <label className="form-label">
                        Username <span className="required">*</span>
                    </label>
                    <input
                        type="text"
                        name="username"
                        value={formData.username}
                        onChange={handleChange}
                        className={`form-input ${errors.username ? 'error' : ''}`}
                        placeholder="username123"
                    />
                    {errors.username && <span className="form-error">{errors.username}</span>}
                </div>

                <div className="form-group">
                    <label className="form-label">
                        Email <span className="required">*</span>
                    </label>
                    <input
                        type="email"
                        name="email"
                        value={formData.email}
                        onChange={handleChange}
                        className={`form-input ${errors.email ? 'error' : ''}`}
                        placeholder="example@mail.com"
                    />
                    {errors.email && <span className="form-error">{errors.email}</span>}
                </div>
            </div>

            <div className="form-group">
                <label className="form-label">
                    Отображаемое имя <span className="required">*</span>
                </label>
                <input
                    type="text"
                    name="display_name"
                    value={formData.display_name}
                    onChange={handleChange}
                    className={`form-input ${errors.display_name ? 'error' : ''}`}
                    placeholder="Иван Иванов"
                />
                {errors.display_name && <span className="form-error">{errors.display_name}</span>}
            </div>

            <div className="form-row">
                <div className="form-group">
                    <label className="form-label">
                        Пароль <span className="required">*</span>
                    </label>
                    <input
                        type="password"
                        name="password"
                        value={formData.password}
                        onChange={handleChange}
                        className={`form-input ${errors.password ? 'error' : ''}`}
                        placeholder="••••••••"
                    />
                    {errors.password && <span className="form-error">{errors.password}</span>}

                    {formData.password && (
                        <div className="password-strength">
                            <div className="password-strength-bar">
                                <div
                                    className="password-strength-fill"
                                    style={{
                                        width: `${passwordStrength.percent}%`,
                                        backgroundColor: passwordStrength.color
                                    }}
                                />
                            </div>
                            <div className="password-strength-label" style={{ color: passwordStrength.color }}>
                                {passwordStrength.label}
                            </div>
                        </div>
                    )}
                </div>

                <div className="form-group">
                    <label className="form-label">
                        Подтвердите пароль <span className="required">*</span>
                    </label>
                    <input
                        type="password"
                        name="confirmPassword"
                        value={formData.confirmPassword}
                        onChange={handleChange}
                        className={`form-input ${errors.confirmPassword ? 'error' : ''}`}
                        placeholder="••••••••"
                    />
                    {errors.confirmPassword && <span className="form-error">{errors.confirmPassword}</span>}
                </div>
            </div>

            <div className="form-group">
                <label className="form-label">Телефон (опционально)</label>
                <input
                    type="tel"
                    name="phone"
                    value={formData.phone}
                    onChange={handleChange}
                    onBlur={handlePhoneBlur}
                    className="form-input"
                    placeholder="+7 999 123 45 67"
                />
            </div>

            <div className="form-row">
                <div className="form-group">
                    <label className="form-label">Страна (опционально)</label>
                    <select
                        name="country"
                        value={formData.country}
                        onChange={handleChange}
                        className="form-select"
                    >
                        {countries.map(country => (
                            <option key={country.value} value={country.value}>
                                {country.label}
                            </option>
                        ))}
                    </select>
                </div>

                <div className="form-group">
                    <label className="form-label">Город (опционально)</label>
                    <select
                        name="city"
                        value={formData.city}
                        onChange={handleChange}
                        className="form-select"
                        disabled={!formData.country}
                    >
                        {availableCities.length > 0 ? (
                            availableCities.map(city => (
                                <option key={city.value} value={city.value}>
                                    {city.label}
                                </option>
                            ))
                        ) : (
                            <option value="">Сначала выберите страну</option>
                        )}
                    </select>
                </div>
            </div>

            <button type="submit" className="auth-submit" disabled={loading}>
                {loading ? 'Регистрация...' : 'Зарегистрироваться'}
            </button>
        </form>
    );
}

export default RegisterForm;
