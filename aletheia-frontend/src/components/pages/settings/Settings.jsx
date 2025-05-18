import React, { useState } from "react";
import TokenDialog from './TokenDialog'; // Убедитесь, что путь правильный
import authClient from '../api/authClient'; // Убедитесь, что путь правильный
import { jwtDecode } from 'jwt-decode'; // Импортируем jwt-decode

function Settings() {
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [token, setToken] = useState(''); // Исходный токен
    const [formattedToken, setFormattedToken] = useState(''); // Отформатированный токен
    const [expirationDate, setExpirationDate] = useState(''); // Дата действия

    const handleGetToken = async () => {
        try {
            const response = await authClient.post('sdk-token');
            const { access_token } = response.data;

            if (access_token) {
                // Устанавливаем исходный токен
                setToken(access_token);

                // Декодируем токен для получения даты истечения
                const decodedToken = jwtDecode(access_token);
                const expDate = new Date(decodedToken.exp * 1000).toLocaleDateString();
                setExpirationDate(expDate);

                // Разбиваем токен на части и форматируем с точками
                const tokenParts = access_token.split('.');
                const formatted = tokenParts
                    .map((part, index) => (index < tokenParts.length - 1 ? `${part}.` : part))
                    .join('\n');
                setFormattedToken(formatted);

                // Открываем диалоговое окно
                setIsDialogOpen(true);
            } else {
                console.error('Токен не получен в ответе');
            }
        } catch (error) {
            console.error('Ошибка при получении токена:', error);
        }
    };

    const handleCloseDialog = () => {
        setIsDialogOpen(false);
    };

    return (
        <section id="settings" className="mb-12">
            <header className="mb-6 flex justify-between items-center">
                <span>
                    <h1 className="text-2xl font-bold">Settings</h1>
                    <p className="text-gray-600">Get your token for SDK</p>
                </span>
            </header>
            <div className="space-y-4">
                <button
                    onClick={handleGetToken}
                    className="px-4 py-2  rounded-md hover:bg-[#2044B8] hover:text-[#DBE9FE] bg-[#DBE9FE] text-[#2044B8]"
                >
                    Получить токен
                </button>
                {isDialogOpen && (
                    <TokenDialog
                        token={formattedToken} // Отформатированный токен для отображения
                        originalToken={token} // Исходный токен для копирования
                        expirationDate={expirationDate}
                        onClose={handleCloseDialog}
                    />
                )}
            </div>
        </section>
    );
}

export default Settings;