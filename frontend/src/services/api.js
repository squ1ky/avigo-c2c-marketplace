import axios from 'axios';

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
    withCredentials: true,  // Send cookies with request
    headers: {
        'Content-Type': 'application/json',
    },
});

const errorMessages = {
    'USER_NOT_FOUND': 'Пользователь не найден',
    'USER_ALREADY_EXISTS': 'Пользователь с таким email уже существует',
    'INVALID_CREDENTIALS': 'Неверный email или пароль',
    'EMAIL_NOT_VERIFIED': 'Email не подтверждён. Проверьте почту',
    'INVALID_CONFIRMATION': 'Неверный код подтверждения',
    'CONFIRMATION_EXPIRED': 'Код подтверждения истёк. Запросите новый',
    'USER_ALREADY_VERIFIED': 'Email уже подтверждён',
    'VALIDATION_FAILED': 'Ошибка валидации данных',
    'INTERNAL_ERROR': 'Внутренняя ошибка сервера',
};

api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.data?.code) {
            const code = error.response.data.code;
            error.message = errorMessages[code] || error.response.data.error;
        }
        return Promise.reject(error);
    }
);

// === API ===

export const register = async (data) => {
    const response = await api.post('/auth/register', data);
    return response.data.data;
};

export const confirmEmail = async (userId, code) => {
    const response = await api.post('/auth/confirm-email', {
        user_id: userId,
        code: code,
    });
    return response.data;
};

export const login = async (identifier, password) => {
    const response = await api.post('/auth/login', {
        identifier,
        password,
    });
    return response.data;
};

export const logout = async () => {
    const response = await api.post('/auth/logout');
    return response.data;
};

// === User ===

export const getCurrentUser = async () => {
    const response = await api.get('/auth/me');
    return response.data.user;
};

export const getProfile = async () => {
    const response = await api.get('/users/me/profile');
    return response.data.profile;
};

export const updateProfile = async (data) => {
    const response = await api.patch('/users/me/profile', data);
    return response.data.profile;
};

export const changePassword = async (currentPassword, newPassword) => {
    const response = await api.post('/auth/change-password', {
        current_password: currentPassword,
        new_password: newPassword,
    });
    return response.data;
};