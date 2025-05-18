import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useSnackbar } from 'notistack';
import authClient from "../api/authClient.js";

function AuthPage({ setIsAuthenticated }) {
    const [isLogin, setIsLogin] = useState(true); // true — вход, false — регистрация
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState(''); // Для регистрации
    const navigate = useNavigate();
    const { enqueueSnackbar } = useSnackbar();

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (isLogin) {
            // Логика входа
            if (!username || !password) {
                enqueueSnackbar('Пожалуйста, заполните все поля', { variant: 'error' });
                return;
            }
            try {
                const response = await authClient.post("/login", {
                    username: username,
                    password: password,
                });
                const { access_token, refresh_token } = response.data; // Предполагаем, что сервер возвращает эти поля

                if (access_token) {
                    localStorage.setItem('accessToken', access_token);
                    if (refresh_token) {
                        localStorage.setItem('refreshToken', refresh_token);
                    }
                    setIsAuthenticated(true);
                    navigate('/dashboard');
                    enqueueSnackbar('Авторизация успешна!', { variant: 'success' });
                } else {
                    enqueueSnackbar('Ошибка авторизации: токен не получен', { variant: 'error' });
                }
            } catch (error) {
                console.error('Ошибка при входе:', error);
                enqueueSnackbar('Ошибка при входе. Проверьте данные и попробуйте снова.', { variant: 'error' });
            }
        } else {
            // Логика регистрации
            if (!username || !password || !confirmPassword) {
                enqueueSnackbar('Пожалуйста, заполните все поля', { variant: 'error' });
                return;
            }
            if (password !== confirmPassword) {
                enqueueSnackbar('Пароли не совпадают', { variant: 'error' });
                return;
            }
            try {
                await authClient.post("/register", {
                    username: username,
                    password: password,
                });
                enqueueSnackbar('Регистрация успешна! Теперь вы можете войти.', { variant: 'success' });
                setIsLogin(true); // Переключаем на форму входа
                setUsername('');
                setPassword('');
                setConfirmPassword('');
            } catch (error) {
                console.error('Ошибка при регистрации:', error);
                enqueueSnackbar('Ошибка при регистрации. Попробуйте снова.', { variant: 'error' });
            }
        }
    };

    return (
        <div className="min-h-screen bg-gray-50">
            <header className="bg-white shadow-sm sticky top-0 z-50">
                <div className="container mx-auto px-4">
                    <div className="flex justify-center items-center py-3">
                        <div className="flex items-center space-x-4">
                            <div className="flex items-center">
                          <span className="text-blue-700 text-2xl font-bold tracking-[45%]">
                            Aletheia
                          </span>
                            </div>
                        </div>
                    </div>
                </div>
            </header>

            {/* Форма */}
            <div className="flex items-center justify-center h-[calc(100vh-80px)]">
                <div className="max-w-md w-full space-y-6 p-8 bg-white rounded-lg shadow-md">
                    <div className="text-center">
                        <h2 className="text-2xl font-bold text-gray-900">
                            {isLogin ? 'Вход' : 'Регистрация'}
                        </h2>
                    </div>
                    <form className="space-y-4" onSubmit={handleSubmit}>
                        <div className="space-y-4">
                            <div>
                                <label htmlFor="username" className="block text-sm font-medium text-gray-700">
                                    Логин
                                </label>
                                <input
                                    id="username"
                                    name="username"
                                    type="text"
                                    required
                                    className="placeholder:text-gray-400 mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                    placeholder="Логин"
                                    value={username}
                                    onChange={(e) => setUsername(e.target.value)}
                                />
                            </div>
                            <div>
                                <label htmlFor="password" className="block text-sm font-medium text-gray-700">
                                    Пароль
                                </label>
                                <input
                                    id="password"
                                    name="password"
                                    type="password"
                                    required
                                    className="placeholder:text-gray-400 mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                    placeholder="Пароль"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                />
                            </div>
                            {!isLogin && (
                                <div>
                                    <label htmlFor="confirm-password"
                                           className="block text-sm font-medium text-gray-700">
                                        Подтвердите пароль
                                    </label>
                                    <input
                                        id="confirm-password"
                                        name="confirm-password"
                                        type="password"
                                        required
                                        className="placeholder:text-gray-400 mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                        placeholder="Подтвердите пароль"
                                        value={confirmPassword}
                                        onChange={(e) => setConfirmPassword(e.target.value)}
                                    />
                                </div>
                            )}
                        </div>

                        <div>
                            <button
                                type="submit"
                                className="w-full flex justify-center py-2 px-4 border border-indigo-300 text-sm font-medium rounded-md text-indigo-600 bg-indigo-50 hover:bg-indigo-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                            >
                                {isLogin ? 'Войти' : 'Зарегистрироваться'}
                            </button>
                        </div>
                    </form>
                    <div className="text-center">
                        <button
                            onClick={() => setIsLogin(!isLogin)}
                            className="font-medium text-indigo-600 hover:text-indigo-500"
                        >
                            {isLogin ? 'Зарегистрироваться' : 'Уже есть аккаунт? Войти'}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default AuthPage;