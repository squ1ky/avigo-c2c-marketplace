import { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { confirmEmail } from '../../services/api';
import toast from 'react-hot-toast';

function EmailConfirmForm() {
    const navigate = useNavigate();
    const [code, setCode] = useState(['', '', '', '', '', '']);
    const [loading, setLoading] = useState(false);
    const [resendTimer, setResendTimer] = useState(60);
    const [canResend, setCanResend] = useState(false);
    const [isReady, setIsReady] = useState(false);
    const inputRefs = useRef([]);

    const userId = localStorage.getItem('pendingUserId');
    const userEmail = localStorage.getItem('pendingUserEmail');

    // Проверка наличия данных
    useEffect(() => {
        if (!userId || !userEmail) {
            toast.error('Сначала зарегистрируйтесь');
            navigate('/auth/register');
        } else {
            setIsReady(true);
        }
    }, [userId, userEmail, navigate]);

    useEffect(() => {
        if (resendTimer > 0) {
            const timer = setTimeout(() => setResendTimer(resendTimer - 1), 1000);
            return () => clearTimeout(timer);
        } else {
            setCanResend(true);
        }
    }, [resendTimer]);

    const handleChange = (index, value) => {
        if (!/^\d*$/.test(value)) return;

        const newCode = [...code];
        newCode[index] = value;
        setCode(newCode);

        if (value && index < 5) {
            inputRefs.current[index + 1]?.focus();
        }

        if (index === 5 && value) {
            const fullCode = [...newCode.slice(0, 5), value].join('');
            if (fullCode.length === 6) {
                handleSubmit(fullCode);
            }
        }
    };

    // Backspace
    const handleKeyDown = (index, e) => {
        if (e.key === 'Backspace' && !code[index] && index > 0) {
            inputRefs.current[index - 1]?.focus();
        }
    };

    const handlePaste = (e) => {
        e.preventDefault();
        const pastedData = e.clipboardData.getData('text').slice(0, 6);

        if (!/^\d+$/.test(pastedData)) return;

        const newCode = pastedData.split('');
        setCode([...newCode, ...Array(6 - newCode.length).fill('')]);

        const lastIndex = Math.min(pastedData.length, 5);
        inputRefs.current[lastIndex]?.focus();

        // Автоотправка если вставили 6 цифр
        if (pastedData.length === 6) {
            handleSubmit(pastedData);
        }
    };

    const handleSubmit = async (fullCode = null) => {
        const codeToSend = fullCode || code.join('');

        if (codeToSend.length !== 6) {
            toast.error('Введите 6-значный код');
            return;
        }

        setLoading(true);

        try {
            await confirmEmail(userId, codeToSend);

            toast.success('Email подтверждён! Теперь войдите в аккаунт');

            // Очищаем localStorage
            localStorage.removeItem('pendingUserId');
            localStorage.removeItem('pendingUserEmail');

            // Переход на страницу логина
            navigate('/auth/login');
        } catch (error) {
            toast.error(error.message || 'Неверный код');
            setCode(['', '', '', '', '', '']);
            inputRefs.current[0]?.focus();
        } finally {
            setLoading(false);
        }
    };

    const handleResend = async () => {
        if (!canResend) return;

        try {
            // TODO: Добавить API метод для повторной отправки кода
            // await resendConfirmationCode(userId);

            toast.success('Код отправлен повторно');
            setResendTimer(60);
            setCanResend(false);
        } catch (error) {
            toast.error('Ошибка отправки кода');
        }
    };

    // Не показываем форму пока нет данных
    if (!isReady) {
        return null;
    }

    return (
        <div>
            <p style={{ textAlign: 'center', marginBottom: '1rem', color: 'var(--gray-medium)' }}>
                Мы отправили код подтверждения на<br />
                <strong>{userEmail}</strong>
            </p>

            <div className="code-input-group" onPaste={handlePaste}>
                {code.map((digit, index) => (
                    <input
                        key={index}
                        ref={(el) => (inputRefs.current[index] = el)}
                        type="text"
                        maxLength={1}
                        value={digit}
                        onChange={(e) => handleChange(index, e.target.value)}
                        onKeyDown={(e) => handleKeyDown(index, e)}
                        className={`code-input ${digit ? 'filled' : ''}`}
                        disabled={loading}
                    />
                ))}
            </div>

            <button
                type="button"
                className="auth-submit"
                onClick={() => handleSubmit()}
                disabled={loading || code.join('').length !== 6}
            >
                {loading ? 'Проверка...' : 'Подтвердить'}
            </button>

            <div className="resend-timer">
                {canResend ? (
                    <button className="resend-button" onClick={handleResend}>
                        Отправить код повторно
                    </button>
                ) : (
                    <span>Отправить код повторно через {resendTimer} сек</span>
                )}
            </div>
        </div>
    );
}

export default EmailConfirmForm;
