import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../../services/api';
import toast from 'react-hot-toast';

function LoginForm() {
    const navigate = useNavigate();
    const [loading, setLoading] = useState(false);
    const [formData, setFormData] = useState({
        identifier: '',
        password: '',
        rememberMe: false
    });
    const [errors, setErrors] = useState({});

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: type === 'checkbox' ? checked : value
        }));

        if (errors[name]) {
            setErrors(prev => ({ ...prev, [name]: '' }));
        }
    };

    const validateForm = () => {
        const newErrors = {};

        if (!formData.identifier) {
            newErrors.identifier = 'Введите email или username';
        }

        if (!formData.password) {
            newErrors.password = 'Введите пароль';
        } else if (formData.password.length < 8) {
            newErrors.password = 'Пароль должен быть минимум 8 символов';
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!validateForm()) {
            toast.error('Пожалуйста, заполните все поля');
            return;
        }

        setLoading(true);

        try {
            const response = await login(formData.identifier, formData.password);

            // Сохраняем данные пользователя в localStorage
            if (response.user) {
                localStorage.setItem('user', JSON.stringify(response.user));
            }

            toast.success('Вход выполнен успешно!');

            navigate('/');
        } catch (error) {
            toast.error(error.message || 'Ошибка входа');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form className="auth-form" onSubmit={handleSubmit}>
            <div className="form-group">
                <label className="form-label">
                    Email или Username <span className="required">*</span>
                </label>
                <input
                    type="text"
                    name="identifier"
                    value={formData.identifier}
                    onChange={handleChange}
                    className={`form-input ${errors.identifier ? 'error' : ''}`}
                    placeholder="example@mail.com или username"
                    autoComplete="username"
                />
                {errors.identifier && <span className="form-error">{errors.identifier}</span>}
            </div>

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
                    autoComplete="current-password"
                />
                {errors.password && <span className="form-error">{errors.password}</span>}
            </div>

            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', cursor: 'pointer' }}>
                    <input
                        type="checkbox"
                        name="rememberMe"
                        checked={formData.rememberMe}
                        onChange={handleChange}
                        style={{ cursor: 'pointer' }}
                    />
                    <span style={{ fontSize: '0.9rem', color: 'var(--gray-medium)' }}>
            Запомнить меня
          </span>
                </label>

                <a
                    href="#"
                    style={{
                        fontSize: '0.9rem',
                        color: 'var(--go-cyan)',
                        textDecoration: 'none',
                        fontWeight: 500
                    }}
                    onClick={(e) => {
                        e.preventDefault();
                        toast('Функция восстановления пароля в разработке', { icon: '🔧' });
                    }}
                >
                    Забыли пароль?
                </a>
            </div>

            <button type="submit" className="auth-submit" disabled={loading}>
                {loading ? 'Вход...' : 'Войти'}
            </button>
        </form>
    );
}

export default LoginForm;
