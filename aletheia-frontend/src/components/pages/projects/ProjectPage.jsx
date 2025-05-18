import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import aletheiaClient from "../api/aletheiaClient.js"; // Предполагается, что apiClient настроен для работы с API

function ProjectPage() {
    const navigate = useNavigate();
    const { id } = useParams(); // Получаем id проекта из URL
    const [services, setServices] = useState([]); // Список сервисов, изначально пустой массив
    const [projectName, setProjectName] = useState(''); // Название проекта
    const [isLoading, setIsLoading] = useState(true); // Состояние загрузки
    const [error, setError] = useState(null); // Состояние ошибки

    // Функция для загрузки данных проекта из API
    const fetchProjectData = async () => {
        try {
            const response = await aletheiaClient.get(`/project/${id}`);
            const projectData = response.data.project.project;
            setProjectName(projectData.projectName);
            setServices(projectData.services || []); // Устанавливаем services, гарантируя, что это массив
            setIsLoading(false);
        } catch (error) {
            console.error('Ошибка при загрузке данных проекта:', error);
            setError('Ошибка при загрузке данных');
            setIsLoading(false);
        }
    };

    // Выполняем запрос при монтировании компонента или изменении id
    useEffect(() => {
        if (id) {
            fetchProjectData();
        }
    }, [id]);

    // Обработка состояния ошибки
    if (error) {
        return <div className="text-red-500">{error}</div>;
    }

    // Обработка состояния загрузки
    if (isLoading) {
        return <div>Загрузка...</div>;
    }

    return (
        <div className="flex-1">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold text-gray-800">{projectName}</h1>
            </div>

            <div className="grid grid-cols-1 gap-4">
                {services.map((service, index) => (
                    <div
                        key={index}
                        className="bg-white rounded-lg shadow-md p-4 border border-gray-200 hover:shadow-lg transition-shadow"
                    >
                        <div className="text-lg font-semibold mb-4">{service.serviceName}</div>

                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <h3 className="font-medium mb-2">События за последние 24 часа</h3>
                                <div className="max-h-48 overflow-y-auto border border-gray-200 rounded p-2 bg-white">
                                    {service.events && service.events.length > 0 ? (
                                        service.events.map((event, idx) => (
                                            <div
                                                key={idx}
                                                className="p-2 border-b border-gray-100 last:border-b-0 cursor-pointer"
                                                onClick={() => navigate(`/events/${event.eventType}`)}
                                            >
                                                <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded text-xs font-medium">
                                                    {event.language}
                                                </span>
                                                <span className="ml-2 text-sm">{event.eventType}</span>
                                                <div className="text-xs text-gray-500">
                                                    Сервис: {event.serviceName}, Количество: {event.eventsCount}
                                                </div>
                                            </div>
                                        ))
                                    ) : (
                                        <p>Нет событий за последние 24 часа</p>
                                    )}
                                </div>
                            </div>

                            <div>
                                <h3 className="font-medium mb-2">Прикрепленные правила</h3>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <h4 className="font-medium mb-2">Error Rules</h4>
                                        <div className="bg-white border border-gray-200 rounded p-2">
                                            <ul className="list-none p-0 m-0">
                                                {service.errorRules && service.errorRules.length > 0 ? (
                                                    service.errorRules.map((rule, idx) => (
                                                        <li
                                                            key={idx}
                                                            className="p-2 border-b border-gray-100 last:border-b-0"
                                                        >
                                                            <div className="flex items-center">
                                                                <i className="fas fa-check-circle text-green-500 mr-2"></i>
                                                                {rule.ruleName}
                                                            </div>
                                                        </li>
                                                    ))
                                                ) : (
                                                    <li>Нет прикрепленных правил</li>
                                                )}
                                            </ul>
                                        </div>
                                    </div>
                                    <div>
                                        <h4 className="font-medium mb-2">Resources Rules</h4>
                                        <div className="bg-white border border-gray-200 rounded p-2">
                                            <ul className="list-none p-0 m-0">
                                                {service.resourceRules && service.resourceRules.length > 0 ? (
                                                    service.resourceRules.map((rule, idx) => (
                                                        <li
                                                            key={idx}
                                                            className="p-2 border-b border-gray-100 last:border-b-0"
                                                        >
                                                            <div className="flex items-center">
                                                                <i className="fas fa-check-circle text-green-500 mr-2"></i>
                                                                {rule.ruleName}
                                                            </div>
                                                        </li>
                                                    ))
                                                ) : (
                                                    <li>Нет прикрепленных правил</li>
                                                )}
                                            </ul>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}

export default ProjectPage;