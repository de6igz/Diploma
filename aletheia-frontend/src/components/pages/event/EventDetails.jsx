import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import aletheiaClient from "../api/aletheiaClient.js"; // Предполагается, что apiClient настроен для работы с API

function EventDetailsPage() {
    const navigate = useNavigate();
    const { id } = useParams(); // Получаем ID события из URL
    const [eventData, setEventData] = useState(null); // Состояние для данных события

    // Функция для загрузки данных события
    const fetchEventDetails = async () => {
        try {
            const response = await aletheiaClient.get(`event/?eventId=${id}`);
            setEventData(response.data.resp.event);
        } catch (error) {
            console.error('Ошибка при загрузке деталей события:', error);
            // Здесь можно добавить уведомление для пользователя
        }
    };

    // Загружаем данные при монтировании компонента
    useEffect(() => {
        fetchEventDetails();
    }, [id]);

    // Если данные еще не загружены, показываем индикатор загрузки
    if (!eventData) {
        return <div>Загрузка...</div>;
    }

    // Определяем иконки для типов действий
    const actionIcons = {
        TELEGRAM: 'fab fa-telegram text-blue-500',
        DISCORD: 'fab fa-discord text-indigo-500',
        EMAIL: 'fas fa-envelope text-orange-500',
    };

    return (
        <div className="">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold text-gray-900">Event Details</h1>
            </div>

            <div className="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200">
                <div className="bg-gray-50 p-4 flex justify-between items-center border-b border-gray-200">
                    <div className="flex items-center">
                        <span className="status-indicator status-error mr-2"></span>
                        <span className="font-semibold">{eventData.eventType} ({eventData.language})</span>
                    </div>
                    <div>
                        <button className="text-gray-500 hover:text-gray-700 mr-2">
                            <i className="fas fa-copy"></i>
                        </button>
                        <button className="text-gray-500 hover:text-gray-700">
                            <i className="fas fa-ellipsis-v"></i>
                        </button>
                    </div>
                </div>
                <div className="p-5">
                    <div className="mb-4">
                        <h3 className="text-gray-500 text-sm mb-2">Последнее событие</h3>
                        <div className="code-block">{eventData.log}</div>
                    </div>

                    <div className="grid grid-cols-2 gap-4 mb-4">
                        <div>
                            <p className="text-sm text-gray-500 mb-1">Дата:</p>
                            <p className="font-medium">{eventData.timestamp}</p>
                        </div>
                        <div>
                            <p className="text-sm text-gray-500 mb-1">Имя сервиса:</p>
                            <p className="font-medium">{eventData.serviceName}</p>
                        </div>
                    </div>

                    <div className="mb-4">
                        <p className="text-sm text-gray-500 mb-1">Примененные правила:</p>
                        <ul className="list-disc pl-5">
                            {eventData.usedRules.map((rule, index) => (
                                <li key={index} className="font-medium">{rule.rule_name}</li>
                            ))}
                        </ul>
                    </div>

                    <div className="mb-4">
                        <p className="text-sm text-gray-500 mb-1">ID лога:</p>
                        <p className="font-medium">{eventData.id}</p>
                    </div>

                    <div className="mb-4">
                        <h3 className="text-gray-700 font-semibold mb-3">Примененные действия:</h3>
                        {eventData.usedActions.map((action, index) => (
                            <div key={index} className="user-action mb-3">
                                <div className="user-action-header">
                                    <div className="flex items-center">
                                        <i className={`${actionIcons[action.type] || 'fas fa-question'} mr-2`}></i>
                                        <span className="font-medium">{action.type}</span>
                                    </div>
                                    <span className="badge badge-success">Отправлено</span>
                                </div>
                                <div className="user-action-content">
                                    <div className="grid grid-cols-2 gap-2">
                                        {Object.entries(action.params).map(([key, value]) => (
                                            <div key={key}>
                                                <p className="text-xs text-gray-500 mb-1">{key}:</p>
                                                <p className="text-sm">{value}</p>
                                            </div>
                                        ))}
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>

                    <div className="flex justify-end">
                        <button className="mr-2 px-4 py-2 text-gray-600 border border-gray-300 rounded-md text-sm">
                            Игнорировать
                        </button>
                        <button
                            className="px-4 py-2 text-gray-600 border border-gray-300 rounded-md text-sm"
                            onClick={() => navigate(`/events`)}
                        >
                            Посмотреть все события
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default EventDetailsPage;