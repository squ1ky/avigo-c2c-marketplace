import { Link } from 'react-router-dom';
import RegisterForm from '../components/auth/RegisterForm';
import '../styles/auth.css';

function RegisterPage() {
    return (
        <div className="auth-page">
            <div className="auth-card">
                <div className="auth-header">
                    <div className="auth-logo">
                        <span className="logo-avi">Avi</span>
                        <span className="logo-go">Go</span>
                    </div>
                    <h1 className="auth-title">Регистрация</h1>
                    <p className="auth-subtitle">Создайте аккаунт и начните продавать</p>
                </div>

                <RegisterForm />

                <div className="auth-footer">
                    Уже есть аккаунт? <Link to="/auth/login">Войти</Link>
                </div>
            </div>
        </div>
    );
}

export default RegisterPage;
