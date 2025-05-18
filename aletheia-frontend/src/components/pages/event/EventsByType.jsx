import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import aletheiaClient from "../api/aletheiaClient.js"; // Предполагается, что apiClient настроен для работы с API
//TODO: сделать фильт по времени, удалить фильтр по уровням
function EventsByType() {
    // Состояние для фильтров и данных
    const [serviceFilter, setServiceFilter] = useState('Все сервисы');
    const [levelFilter, setLevelFilter] = useState('Все уровни');
    const [periodFilter, setPeriodFilter] = useState('Последние 24 часа');
    const [events, setEvents] = useState([]);
    const [filteredEvents, setFilteredEvents] = useState([]);
    const [servicesList, setServicesList] = useState([]); // Динамический список сервисов
    const navigate = useNavigate();
    const { eventType } = useParams();

    // Загрузка данных событий из API
    const fetchEvents = async () => {
        try {
            const response = await aletheiaClient.get(`/events_by_event_type?eventType=${eventType}`);
            const eventsData = response.data.resp.events;
            setEvents(eventsData);
            setFilteredEvents(eventsData);

            // Формируем список уникальных сервисов для фильтра
            const uniqueServices = [...new Set(eventsData.map(event => event.serviceName))];
            setServicesList(uniqueServices);
        } catch (error) {
            console.error('Ошибка при загрузке событий:', error);
        }
    };

    useEffect(() => {
        fetchEvents();
    }, [eventType]);

    // Фильтрация событий при изменении фильтров
    useEffect(() => {
        let filtered = events;
        if (serviceFilter !== 'Все сервисы') {
            filtered = filtered.filter((event) => event.serviceName === serviceFilter);
        }
        if (levelFilter !== 'Все уровни') {
            filtered = filtered.filter((event) => event.level === levelFilter);
        }
        // Фильтрация по периоду (пока заглушка, можно доработать позже)
        setFilteredEvents(filtered);
    }, [serviceFilter, levelFilter, periodFilter, events]);

    // Вычисление сводки
    const summary = {
        total: filteredEvents.length,
        info: filteredEvents.filter((e) => e.level === 'INFO').length,
        warn: filteredEvents.filter((e) => e.level === 'WARN').length,
        error: filteredEvents.filter((e) => e.level === 'ERROR').length,
    };

    return (
        <div className="flex-1">
            {/* Заголовок */}
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold">
                    События типа "{eventType}" за последние 24 часа
                </h1>
                <div className="space-x-2">
                    <button className="btn btn-gray" onClick={() => navigate(`/projects`)}>
                        <i className="fas fa-arrow-left mr-2"></i> Назад
                    </button>
                    <button className="btn btn-orange">
                        <i className="fas fa-filter mr-2"></i> Фильтр по сервисам
                    </button>
                </div>
            </div>

            {/* Фильтры */}
            <div className="filter-bar flex flex-wrap gap-3 items-center mb-6">
                <div>
                    <label className="text-sm text-gray-500 block mb-1">Сервис</label>
                    <select
                        className="filter-select min-w-[150px]"
                        value={serviceFilter}
                        onChange={(e) => setServiceFilter(e.target.value)}
                    >
                        <option>Все сервисы</option>
                        {servicesList.map((service, index) => (
                            <option key={index} value={service}>
                                {service}
                            </option>
                        ))}
                    </select>
                </div>
                <div>
                    <label className="text-sm text-gray-500 block mb-1">Уровень</label>
                    <select
                        className="filter-select min-w-[150px]"
                        value={levelFilter}
                        onChange={(e) => setLevelFilter(e.target.value)}
                    >
                        <option>Все уровни</option>
                        <option>INFO</option>
                        <option>WARN</option>
                        <option>ERROR</option>
                    </select>
                </div>
                <div>
                    <label className="text-sm text-gray-500 block mb-1">Период</label>
                    <select
                        className="filter-select min-w-[150px]"
                        value={periodFilter}
                        onChange={(e) => setPeriodFilter(e.target.value)}
                    >
                        <option>Последние 24 часа</option>
                        <option>Последние 12 часов</option>
                        <option>Последние 6 часов</option>
                        <option>Последний час</option>
                    </select>
                </div>
                <div className="ml-auto">
                    <button className="btn btn-orange">
                        <i className="fas fa-search mr-2"></i> Применить фильтры
                    </button>
                </div>
            </div>

            {/* Сводка */}
            <div className="bg-white p-4 mb-6 border border-gray-200 rounded-md">
                <h3 className="text-lg font-semibold mb-3">Сводка</h3>
                <div className="grid grid-cols-4 gap-4">
                    <div className="text-center p-3 bg-gray-50 rounded-md">
                        <div className="text-2xl font-bold">{summary.total}</div>
                        <div className="text-sm text-gray-500">Всего событий</div>
                    </div>
                    <div className="text-center p-3 bg-blue-50 rounded-md">
                        <div className="text-2xl font-bold text-blue-600">{summary.info}</div>
                        <div className="text-sm text-gray-500">INFO</div>
                    </div>
                    <div className="text-center p-3 bg-yellow-50 rounded-md">
                        <div className="text-2xl font-bold text-yellow-600">{summary.warn}</div>
                        <div className="text-sm text-gray-500">WARN</div>
                    </div>
                    <div className="text-center p-3 bg-red-50 rounded-md">
                        <div className="text-2xl font-bold text-red-600">{summary.error}</div>
                        <div className="text-sm text-gray-500">ERROR</div>
                    </div>
                </div>
            </div>

            {/* Список событий */}
            <h3 className="text-lg font-semibold mb-3">События ({filteredEvents.length})</h3>
            {filteredEvents.map((event) => {
                // Парсим log и выводим только поля до timestamp
                let logData;
                try {
                    const parsedLog = JSON.parse(event.log);
                    logData = {
                        os: parsedLog.os,
                        arch: parsedLog.arch,
                        tags: parsedLog.tags,
                        level: parsedLog.level,
                        user_id: parsedLog.user_id,
                        version: parsedLog.version,
                        language: parsedLog.language,
                    };
                } catch (error) {
                    logData = {};
                    console.error('Ошибка парсинга log:', error);
                }

                return (
                    <div key={event.id} className="event-card p-4 bg-white mb-4">
                        <div className="flex justify-between items-start mb-2">
                            <div>
                                <span className="bg-blue-100 text-blue-800 px-2 py-1 rounded text-xs font-medium">
                                    {event.language}
                                </span>
                                <span className="ml-2 font-semibold">{event.eventType}</span>
                            </div>
                            <div className="text-gray-500 text-sm">
                                <i className="fas fa-clock mr-1"></i>
                                {event.timestamp}
                            </div>
                        </div>
                        <div className="text-sm text-gray-600 mb-3">Сервис: {event.serviceName}</div>
                        <div className="text-sm bg-gray-50 p-3 rounded mb-3 font-mono">
                            {JSON.stringify(logData, null, 2)}
                        </div>
                        <div className="flex justify-between">
                            <div>
                                <span className="text-sm text-gray-500">ID события: {event.id}</span>
                            </div>
                            <div>
                                <button
                                    className="text-blue-600 text-sm mr-3 hover:underline"
                                    onClick={() => navigate(`/events/details/${event.id}`)}
                                >
                                    <i className="fas fa-eye mr-1"></i> Детали
                                </button>
                                <button className="text-gray-600 text-sm hover:underline">
                                    <i className="fas fa-bell-slash mr-1"></i> Игнорировать
                                </button>
                            </div>
                        </div>
                    </div>
                );
            })}

            {/* Пагинация (заглушка) */}
            <div className="flex justify-between items-center mt-6">
                <div className="text-sm text-gray-500">
                    Показано {filteredEvents.length} из {events.length} событий
                </div>
                <div className="flex space-x-2">
                    <button className="px-3 py-1 border border-gray-300 rounded text-gray-600 bg-white">
                        <i className="fas fa-chevron-left"></i>
                    </button>
                    <button className="px-3 py-1 border border-blue-500 rounded text-white bg-blue-500">1</button>
                    <button className="px-3 py-1 border border-gray-300 rounded text-gray-600 bg-white">2</button>
                    <button className="px-3 py-1 border border-gray-300 rounded text-gray-600 bg-white">3</button>
                    <button className="px-3 py-1 border border-gray-300 rounded text-gray-600 bg-white">
                        <i className="fas fa-chevron-right"></i>
                    </button>
                </div>
            </div>
        </div>
    );
}

export default EventsByType;