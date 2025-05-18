import React, { useState, useEffect } from 'react';
import authClient from "../pages/api/authClient.js";

function Header({ onToggleSidebar }) {
    const [user, setUser] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    // Получение данных пользователя из API
    useEffect(() => {
        const fetchUserData = async () => {
            try {
                const response = await authClient.get('/me');
                setUser(response.data); // Ожидаемый ответ: { "id": "1", "userName": "12345" }
                setIsLoading(false);
            } catch (error) {
                console.error('Ошибка при загрузке данных пользователя:', error);
                setIsLoading(false);
            }
        };
        fetchUserData();
    }, []);

    // Функция для получения инициалов из имени
    const getInitials = (name) => {
        if (!name) return '';
        const names = name.split(' ');
        return names
            .map((n) => n[0]?.toUpperCase() || '')
            .join('')
            .slice(0, 2); // Максимум 2 буквы
    };

    // Функция для генерации цвета на основе имени
    const getColorFromName = (name) => {
        if (!name) return '#ccc'; // Цвет по умолчанию, если имя отсутствует
        let hash = 0;
        for (let i = 0; i < name.length; i++) {
            hash = name.charCodeAt(i) + ((hash << 5) - hash);
        }
        const hue = Math.abs(hash) % 360;
        return `hsl(${hue}, 70%, 50%)`;
    };

    // Функция для разлогина
    const handleLogout = () => {
        // Примерная логика разлогина: очистка токенов и перенаправление
        localStorage.removeItem('accessToken');
        localStorage.removeItem('refreshToken');
        window.location.href = '/';
    };

    // Отображение загрузки
    if (isLoading) {
        return <div>Загрузка...</div>;
    }

    // Вычисление инициалов и цвета для текущего пользователя
    const initials = getInitials(user?.userName);
    const backgroundColor = getColorFromName(user?.userName);

    return (
        <header className="bg-white shadow-sm sticky top-0 z-50">
            <div className="container mx-auto px-4">
                <div className="flex justify-between items-center py-3">
                    <div className="flex items-center space-x-4">
                        <button
                            onClick={onToggleSidebar}
                            className="md:hidden text-gray-600 hover:text-blue-600"
                        >
                            <i className="fas fa-bars text-xl"></i>
                        </button>
                        <div className="flex items-center">
                            <span className="text-blue-700 text-2xl font-bold tracking-[45%]">
                                Aletheia
                            </span>
                        </div>
                    </div>
                    <div className="flex items-center space-x-6">
                        <div className="relative">
                            <input
                                type="text"
                                placeholder="Search anything..."
                                className="bg-gray-100 rounded-full py-2 pl-10 pr-4 w-64 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:bg-white transition"
                            />
                            <i className="fas fa-search absolute left-3 top-3 text-gray-500"></i>
                        </div>
                        <div className="flex items-center space-x-4">
                            <button className="text-gray-600 hover:text-blue-600 relative">
                                <i className="fas fa-bell text-xl"></i>
                                <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center">
                                    3
                                </span>
                            </button>
                            <div className="relative">
                                <div
                                    className="flex items-center space-x-2 cursor-pointer"
                                    onClick={() => setIsMenuOpen(!isMenuOpen)}
                                >
                                    <div
                                        className="rounded-full h-8 w-8 flex items-center justify-center text-white font-semibold"
                                        style={{ backgroundColor }}
                                    >
                                        {initials}
                                    </div>
                                    <span className="text-sm font-medium hidden md:inline-block">
                                        {user?.userName || 'Гость'}
                                    </span>
                                    <i className="fas fa-chevron-down text-xs text-gray-500"></i>
                                </div>
                                {isMenuOpen && (
                                    <div className="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-10">
                                        <button
                                            onClick={handleLogout}
                                            className="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                                        >
                                            Выйти
                                        </button>
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </header>
    );
}

export default Header;