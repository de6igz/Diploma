import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import aletheiaClient from '../api/aletheiaClient'; // Предполагается, что aletheiaClient настроен для запросов

function EventsPage() {
    const navigate = useNavigate();
    const [events, setEvents] = useState([]);
    const [filteredEvents, setFilteredEvents] = useState([]);
    const [filters, setFilters] = useState({
        eventType: 'All',
        serviceName: 'All',
        language: 'All',
    });

    // Получение данных из API
    useEffect(() => {
        const fetchEvents = async () => {
            try {
                const response = await aletheiaClient.get('/events');
                const eventsData = response.data.items.events;
                setEvents(eventsData);
                setFilteredEvents(eventsData);
            } catch (error) {
                console.error('Ошибка при получении событий:', error);
            }
        };
        fetchEvents();
    }, []);

    // Собираем уникальные значения для фильтров
    const uniqueEventTypes = [...new Set(events.map(event => event.eventType))];
    const uniqueServiceNames = [...new Set(events.map(event => event.serviceName))];
    const uniqueLanguages = [...new Set(events.map(event => event.language))];

    // Применение фильтров при их изменении
    useEffect(() => {
        let filtered = events;
        if (filters.eventType !== 'All') {
            filtered = filtered.filter(event => event.eventType === filters.eventType);
        }
        if (filters.serviceName !== 'All') {
            filtered = filtered.filter(event => event.serviceName === filters.serviceName);
        }
        if (filters.language !== 'All') {
            filtered = filtered.filter(event => event.language === filters.language);
        }
        setFilteredEvents(filtered);
    }, [filters, events]);

    return (
        <div className="">
            <div className="flex justify-between items-center mb-6">
                <header>
                    <h1 className="text-2xl font-bold">All Events</h1>
                    <p className="text-gray-600">Overview of your events</p>
                </header>
            </div>

            {/* Фильтры */}
            <div className="bg-white rounded-lg shadow-md border border-gray-200 mb-6">
                <div className="bg-gray-50 p-4 border-b border-gray-200">
                    <h2 className="text-lg font-semibold">Event Filters</h2>
                </div>
                <div className="p-5">
                    <div className="grid grid-cols-3 gap-4">
                        <div>
                            <label className="block text-sm text-gray-600 mb-1">Event Type</label>
                            <select
                                className="w-full border rounded-md px-3 py-2 text-sm"
                                value={filters.eventType}
                                onChange={(e) => setFilters({ ...filters, eventType: e.target.value })}
                            >
                                <option>All</option>
                                {uniqueEventTypes.map((type, index) => (
                                    <option key={index} value={type}>{type}</option>
                                ))}
                            </select>
                        </div>
                        <div>
                            <label className="block text-sm text-gray-600 mb-1">Service</label>
                            <select
                                className="w-full border rounded-md px-3 py-2 text-sm"
                                value={filters.serviceName}
                                onChange={(e) => setFilters({ ...filters, serviceName: e.target.value })}
                            >
                                <option>All</option>
                                {uniqueServiceNames.map((service, index) => (
                                    <option key={index} value={service}>{service}</option>
                                ))}
                            </select>
                        </div>
                        <div>
                            <label className="block text-sm text-gray-600 mb-1">Language</label>
                            <select
                                className="w-full border rounded-md px-3 py-2 text-sm"
                                value={filters.language}
                                onChange={(e) => setFilters({ ...filters, language: e.target.value })}
                            >
                                <option>All</option>
                                {uniqueLanguages.map((language, index) => (
                                    <option key={index} value={language}>{language}</option>
                                ))}
                            </select>
                        </div>
                    </div>
                </div>
            </div>

            {/* Таблица событий */}
            <div className="bg-white rounded-lg shadow-md border border-gray-200">
                <div className="bg-gray-50 p-4 flex justify-between items-center border-b border-gray-200">
                    <h2 className="text-lg font-semibold">Event List</h2>
                    <div>
                        <button className="text-gray-600 border border-gray-300 rounded-md px-3 py-1 text-sm mr-2">
                            <i className="fas fa-download mr-1"></i> Export
                        </button>
                        <button className="text-gray-600 border border-gray-300 rounded-md px-3 py-1 text-sm">
                            <i className="fas fa-sync-alt mr-1"></i> Refresh
                        </button>
                    </div>
                </div>
                <div className="p-0">
                    <table className="w-full">
                        <thead>
                        <tr className="bg-gray-50 text-left text-gray-600 text-sm">
                            <th className="px-6 py-3 font-medium">Language</th>
                            <th className="px-6 py-3 font-medium">Event Type</th>
                            <th className="px-6 py-3 font-medium">Service</th>
                            <th className="px-6 py-3 font-medium">Events Count</th>
                            <th className="px-6 py-3 font-medium">Actions</th>
                        </tr>
                        </thead>
                        <tbody>
                        {filteredEvents.map((event, index) => (
                            <tr
                                key={index}
                                className="border-t border-gray-200 hover:bg-gray-50 cursor-pointer"
                                onClick={() => navigate(`/events/${event.eventType}`)}
                            >
                                <td className="px-6 py-4">
                                        <span className="px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-700">
                                            {event.language}
                                        </span>
                                </td>
                                <td className="px-6 py-4">{event.eventType}</td>
                                <td className="px-6 py-4">{event.serviceName}</td>
                                <td className="px-6 py-4">{event.eventsCount}</td>
                                <td className="px-6 py-4">
                                    <button className="text-blue-500 hover:text-blue-700">
                                        <i className="fas fa-eye"></i>
                                    </button>
                                </td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                    <div className="p-4 flex justify-between items-center border-t border-gray-200">
                        <div className="text-sm text-gray-600">Showing 1 to {filteredEvents.length} of {events.length} events</div>
                        <ul className="pagination flex space-x-2">
                            <li>
                                <button className="px-2 py-1 border border-gray-300 rounded text-gray-600">«</button>
                            </li>
                            <li>
                                <button className="px-2 py-1 border border-blue-500 bg-blue-500 text-white rounded">1</button>
                            </li>
                            <li>
                                <button className="px-2 py-1 border border-gray-300 rounded text-gray-600">2</button>
                            </li>
                            <li>
                                <button className="px-2 py-1 border border-gray-300 rounded text-gray-600">3</button>
                            </li>
                            <li>
                                <button className="px-2 py-1 border border-gray-300 rounded text-gray-600">4</button>
                            </li>
                            <li>
                                <button className="px-2 py-1 border border-gray-300 rounded text-gray-600">5</button>
                            </li>
                            <li>
                                <button className="px-2 py-1 border border-gray-300 rounded text-gray-600">»</button>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default EventsPage;