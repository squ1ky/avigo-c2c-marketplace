import { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { login } from '../../services/api';
import { useAuth } from '../../context/AuthContext';
import toast from 'react-hot-toast';

function LoginForm() {
    const navigate = useNavigate();
    const location = useLocation();
    const { setUser } = useAuth();

    const [formData, setFormData] = useState({
        email: '',
        password: ''
    });
    const [loading, setLoading] = useState(false);

    const savedPath = location.state?.from?.pathname;

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);

        try {
            const data = await login(formData.email, formData.password);
            setUser(data.user);
            localStorage.setItem('user', JSON.stringify(data.user));
            toast.success('Вход выполнен');

            const targetPath = savedPath || `/account/${data.user.id}`;
            navigate(targetPath, { replace: true });

        } catch (error) {
            console.error(error);
            toast.error(error.message || 'Ошибка входа');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form className="auth-form" onSubmit={handleSubmit}>
            <div className="form-group">
                <label className="form-label">Email</label>
                <input
                    type="email"
                    name="email"
                    className="form-input"
                    value={formData.email}
                    onChange={handleChange}
                    required
                />
            </div>
            <div className="form-group">
                <label className="form-label">Пароль</label>
                <input
                    type="password"
                    name="password"
                    className="form-input"
                    value={formData.password}
                    onChange={handleChange}
                    required
                />
            </div>
            <button type="submit" className="auth-submit" disabled={loading}>
                {loading ? 'Вход...' : 'Войти'}
            </button>
        </form>
    );
}

export default LoginForm;
