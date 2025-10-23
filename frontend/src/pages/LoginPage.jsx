import { Link } from 'react-router-dom';
import LoginForm from '../components/auth/LoginForm';
import '../styles/auth.css';

function LoginPage() {
    return (
        <div className="auth-page">
            <div className="auth-card">
                <div className="auth-header">
                    <div className="auth-logo">
                        <span className="logo-avi">Avi</span>
                        <span className="logo-go">Go</span>
                    </div>
                    <h1 className="auth-title">Вход</h1>
                    <p className="auth-subtitle">Войдите в свой аккаунт</p>
                </div>

                <LoginForm />

                <div className="auth-footer">
                    Нет аккаунта? <Link to="/auth/register">Зарегистрироваться</Link>
                </div>
            </div>
        </div>
    );
}

export default LoginPage;
