import axios from 'axios';

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1',
    withCredentials: true,
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
        const status = error.response?.status;
        const msg = (error.response?.data?.error || error.message || '').toLowerCase();

        const isAuthError =
            status === 401 ||
            msg.includes('missing authentication') ||
            msg.includes('token expired') ||
            msg.includes('invalid token');

        if (isAuthError) {
            if (!window.location.pathname.includes('/auth/login')) {
                localStorage.removeItem('user');
                window.location.href = '/auth/login';
                return new Promise(() => {});
            }
        }

        if (error.response?.data?.code) {
            const code = error.response.data.code;
            error.message = errorMessages[code] || error.response.data.error;
        }

        return Promise.reject(error);
    }
);

// === Auth ===

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

export const getPublicProfile = async (userId) => {
    const response = await api.get(`/users/${userId}/profile`);
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

// === Media ===

export const uploadMedia = async (file) => {
    const formData = new FormData();
    formData.append('file', file);
    const response = await api.post('/media/upload', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
};

// === Listings ===

export const createListing = async (userId, data) => {
    if (!userId) throw new Error("User ID is required for creating listing");
    const response = await api.post(`/account/${userId}/listings`, data);
    return response.data;
};

export const getListing = async (userId, listingId) => {
    if (!userId) throw new Error("User ID is required to fetch listing details");
    const response = await api.get(`/account/${userId}/listings/${listingId}`);
    return response.data;
};

export const getUserListings = async (userId, params = {}) => {
    if (!userId) return [];
    const response = await api.get(`/account/${userId}/listings`, { params });
    return response.data;
};

export const getCategories = async () => {
    const response = await api.get('/categories');
    return response.data;
};

export const updateListing = async (userId, listingId, data) => {
    const response = await api.put(`/account/${userId}/listings/${listingId}`, data);
    return response.data;
};

export const deleteListing = async (userId, listingId) => {
    await api.delete(`/account/${userId}/listings/${listingId}`);
};

// === Orders ===

export const getUserPurchases = async (userId) => {
    if (!userId) return [];
    const response = await api.get(`/account/${userId}/orders/purchases`);
    return response.data;
};

export const getUserSales = async (userId) => {
    if (!userId) return [];
    const response = await api.get(`/account/${userId}/orders/sales`);
    return response.data;
};