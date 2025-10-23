import { Link } from 'react-router-dom';
import EmailConfirmForm from '../components/auth/EmailConfirmForm';
import '../styles/auth.css';

function EmailConfirmPage() {
    return (
        <div className="auth-page">
            <div className="auth-card">
                <div className="auth-header">
                    <div className="auth-logo">
                        <span className="logo-avi">Avi</span>
                        <span className="logo-go">Go</span>
                    </div>
                    <h1 className="auth-title">Подтверждение Email</h1>
                    <p className="auth-subtitle">Введите код из письма</p>
                </div>

                <EmailConfirmForm />

                <div className="auth-footer">
                    Вернуться к <Link to="/auth/register">регистрации</Link>
                </div>
            </div>
        </div>
    );
}

export default EmailConfirmPage;
