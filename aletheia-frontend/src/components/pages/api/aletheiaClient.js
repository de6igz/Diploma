import axios from 'axios';
import authClient from "./authClient.js";

// Создаем экземпляр с базовыми настройками
const aletheiaClient = axios.create({
    baseURL: 'https://de6igz.ru/aletheia/v1/',
    timeout: 5000,
});

// Функция для разлогинивания пользователя
const logout = () => {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    window.location.href = "/";
};

// Переменные для управления процессом обновления токена
let isRefreshing = false;
let refreshSubscribers = [];

// Функция подписки на обновление токена
const subscribeTokenRefresh = (callback) => {
    refreshSubscribers.push(callback);
};

// Функция обработки завершения обновления токена
const onRefreshed = (token) => {
    refreshSubscribers.forEach((callback) => callback(token));
    refreshSubscribers = [];
};

// Интерсептор запросов: добавляем токен
aletheiaClient.interceptors.request.use(
    (config) => {
        // Пропускаем добавление токена для /auth, если такой запрос есть в aletheiaClient
        if (config.url === '/auth') {
            return config;
        }
        const accessToken = localStorage.getItem('accessToken');
        if (accessToken) {
            config.headers['Authorization'] = `Bearer ${accessToken}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
);

// Интерсептор ответа
aletheiaClient.interceptors.response.use(
    (response) => response, // Успешный ответ пропускаем
    async (error) => {
        console.log("Ошибка в Axios:", error);
        const originalRequest = error.config;

        if (error.response) {
            const status = error.response.status;

            // Ошибка 401: токен недействителен, выполняем logout
            if (status === 401) {
                console.log("Токен недействителен, выполняем logout");
                logout();
                return Promise.reject(error);
            }

            // Ошибка 403: токен истек, обновляем его
            if (status === 403 && !originalRequest._retry) {
                originalRequest._retry = true;

                if (!isRefreshing) {
                    isRefreshing = true;
                    const refreshToken = localStorage.getItem("refreshToken");
                    console.log("Начинаем обновление, refreshToken:", refreshToken);

                    if (refreshToken) {
                        try {
                            const response = await authClient.post("/refresh", {
                                refresh_token: refreshToken,
                            });
                            const { access_token, refresh_token } = response.data;
                            console.log("Новые токены:", access_token, refresh_token);

                            // Сохраняем новые токены
                            localStorage.setItem("accessToken", access_token);
                            if (refresh_token) {
                                localStorage.setItem("refreshToken", refresh_token);
                            }

                            isRefreshing = false;
                            onRefreshed(access_token);
                        } catch (refreshError) {
                            console.error("Ошибка при обновлении токена:", refreshError.response?.status);
                            isRefreshing = false;
                            if (refreshError.response?.status === 401) {
                                logout();
                            }
                            onRefreshed(null);
                            return Promise.reject(refreshError);
                        }
                    } else {
                        console.log("refreshToken отсутствует");
                        isRefreshing = false;
                        logout();
                        onRefreshed(null);
                        return Promise.reject(error);
                    }
                }

                // Ожидаем обновления токена и повторяем запрос
                return new Promise((resolve, reject) => {
                    subscribeTokenRefresh((token) => {
                        if (token) {
                            console.log("Повторяем запрос с новым токеном:", token);
                            originalRequest.headers["Authorization"] = `Bearer ${token}`;
                            resolve(aletheiaClient.request(originalRequest)); // Повторяем через aletheiaClient
                        } else {
                            reject(error);
                        }
                    });
                });
            }
        } else if (error.request) {
            console.log("Нет ответа от сервера:", error.request);
            logout()
        } else {
            console.log("Ошибка настройки запроса:", error.message);
        }

        return Promise.reject(error);
    }
);

export default aletheiaClient;