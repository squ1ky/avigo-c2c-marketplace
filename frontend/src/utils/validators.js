export const validatePassword = (password) => {
    if (!password) {
        return 'Пароль обязателен';
    }

    if (password.length < 8) {
        return 'Пароль должен быть минимум 8 символов';
    }

    const hasUpper = /[A-Z]/.test(password);
    const hasLower = /[a-z]/.test(password);
    const hasNumber = /[0-9]/.test(password);
    const hasSpecial = /[!@#$%^&*(),.?":{}|<>_\-]/.test(password);

    if (!hasUpper) {
        return 'Пароль должен содержать заглавную букву';
    }
    if (!hasLower) {
        return 'Пароль должен содержать строчную букву';
    }
    if (!hasNumber) {
        return 'Пароль должен содержать цифру';
    }
    if (!hasSpecial) {
        return 'Пароль должен содержать специальный символ';
    }

    return null;
};

export const validateUsername = (username) => {
    if (!username) {
        return 'Username обязателен';
    }

    if (username.length < 3 || username.length > 64) {
        return 'Username должен быть 3-64 символа';
    }

    if (!/^[a-zA-Z][a-zA-Z0-9_-]*$/.test(username)) {
        return 'Username должен начинаться с буквы и содержать только буквы, цифры, _ или -';
    }

    return null;
};

export const validateDisplayName = (displayName) => {
    if (!displayName) {
        return 'Имя обязательно';
    }

    if (displayName.length < 1 || displayName.length > 64) {
        return 'Имя должно быть 1-64 символа';
    }

    return null;
};

export const validateEmail = (email) => {
    if (!email) {
        return 'Email обязателен';
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
        return 'Введите корректный email';
    }

    if (email.length > 255) {
        return 'Email слишком длинный';
    }

    return null;
};

export const validatePhone = (phone) => {
    if (!phone) return null;

    if (phone.length > 32) {
        return 'Телефон слишком длинный';
    }

    return null;
};

export const validatePasswordMatch = (password, confirmPassword) => {
    if (password !== confirmPassword) {
        return 'Пароли не совпадают';
    }
    return null;
};

export const getPasswordStrength = (password) => {
    if (!password) return { strength: 0, label: '' };

    let strength = 0;

    if (password.length >= 8) strength++;
    if (password.length >= 12) strength++;
    if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++;
    if (/[0-9]/.test(password)) strength++;
    if (/[!@#$%^&*(),.?":{}|<>_\-]/.test(password)) strength++;

    const labels = ['', 'Слабый', 'Средний', 'Хороший', 'Сильный', 'Очень сильный'];
    const colors = ['', '#ff4444', '#ffaa00', '#44aa44', '#00aa00', '#007700'];

    return {
        strength,
        label: labels[strength] || '',
        color: colors[strength] || '',
        percent: (strength / 5) * 100
    };
};