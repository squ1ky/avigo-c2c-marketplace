import { useState } from 'react';
import { changePassword } from '../services/api';
import {
    validatePassword,
    validatePasswordMatch,
} from '../utils/validators';
import toast from 'react-hot-toast';
import '../styles/auth.css';

function ChangePasswordPage() {
    const [form, setForm] = useState({
        current_password: '',
        new_password: '',
        confirm_new_password: '',
    });
    const [saving, setSaving] = useState(false);
    const [errors, setErrors] = useState({});

    const handleChange = (e) => {
        const { name, value } = e.target;
        setForm((prev) => ({ ...prev, [name]: value }));

        if (errors[name]) {
            setErrors((prev) => ({ ...prev, [name]: '' }));
        }
    };

    const validateForm = () => {
        const newErrors = {};

        const passwordError = validatePassword(form.new_password);
        if (passwordError) {
            newErrors.new_password = passwordError;
        }

        const passwordMatchError = validatePasswordMatch(
            form.new_password,
            form.confirm_new_password,
        );
        if (passwordMatchError) {
            newErrors.confirm_new_password = passwordMatchError;
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setErrors({});

        if (!validateForm()) {
            toast.error('Пожалуйста, исправьте ошибки в форме');
            return;
        }

        setSaving(true);
        try {
            await changePassword(
                form.current_password,
                form.new_password,
            );
            toast.success('Пароль успешно изменён');
            setForm({
                current_password: '',
                new_password: '',
                confirm_new_password: '',
            });
        } catch (err) {
            console.error(err);
            const msg = err.message || 'Не удалось изменить пароль';
            setErrors((prev) => ({ ...prev, submit: msg }));
            toast.error(msg);
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="auth-page">
            <div className="auth-card">
                <div className="auth-header">
                    <div className="auth-logo">
                        <span className="logo-avi">Avi</span>
                        <span className="logo-go">Go</span>
                    </div>
                    <h1 className="auth-title">Смена пароля</h1>
                    <p className="auth-subtitle">
                        Введите текущий и новый пароль
                    </p>
                </div>

                <form className="auth-form" onSubmit={handleSubmit}>
                    {errors.submit && (
                        <div className="form-error">{errors.submit}</div>
                    )}

                    <div className="form-group">
                        <label className="form-label">
                            Текущий пароль <span className="required">*</span>
                        </label>
                        <input
                            type="password"
                            name="current_password"
                            className="form-input"
                            value={form.current_password}
                            onChange={handleChange}
                            required
                        />
                    </div>

                    <div className="form-group">
                        <label className="form-label">
                            Новый пароль <span className="required">*</span>
                        </label>
                        <input
                            type="password"
                            name="new_password"
                            className={`form-input ${
                                errors.new_password ? 'error' : ''
                            }`}
                            value={form.new_password}
                            onChange={handleChange}
                            required
                        />
                        {errors.new_password && (
                            <span className="form-error">
                                {errors.new_password}
                            </span>
                        )}
                    </div>

                    <div className="form-group">
                        <label className="form-label">
                            Подтверждение нового пароля{' '}
                            <span className="required">*</span>
                        </label>
                        <input
                            type="password"
                            name="confirm_new_password"
                            className={`form-input ${
                                errors.confirm_new_password ? 'error' : ''
                            }`}
                            value={form.confirm_new_password}
                            onChange={handleChange}
                            required
                        />
                        {errors.confirm_new_password && (
                            <span className="form-error">
                                {errors.confirm_new_password}
                            </span>
                        )}
                    </div>

                    <button
                        type="submit"
                        className="auth-submit"
                        disabled={saving}
                    >
                        {saving ? 'Сохранение...' : 'Изменить пароль'}
                    </button>
                </form>
            </div>
        </div>
    );
}

export default ChangePasswordPage;
