import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import aletheiaClient from "../api/aletheiaClient.js";

function RecentEventsTable() {
    const [recentEvents, setRecentEvents] = useState([]);

    useEffect(() => {
        async function fetchRecentEvents() {
            try {
                const response = await aletheiaClient.get('/events');
                // API возвращает данные в формате { "items": { "events": [...] } }
                setRecentEvents(response.data.items.events || []);
            } catch (error) {
                console.error('Error fetching recent events:', error);
                setRecentEvents([]); // В случае ошибки устанавливаем пустой массив
            }
        }
        fetchRecentEvents();
    }, []);

    return (
        <div className="bg-white rounded-xl shadow-sm mb-8">
            <div className="p-6 border-b border-gray-100">
                <div className="flex justify-between items-center">
                    <h2 className="text-lg font-semibold">Recent Events</h2>
                    <Link to="/events" className="text-blue-600 hover:text-blue-800 text-sm font-medium">
                        View All
                    </Link>
                </div>
            </div>
            <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                    <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Type
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Source
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Status
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Language
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                            Actions
                        </th>
                    </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                    {recentEvents.length === 0 ? (
                        <tr>
                            <td colSpan="5" className="px-6 py-4 text-center text-gray-500">
                                No recent events
                            </td>
                        </tr>
                    ) : (
                        recentEvents.slice(0, 3).map((event) => (
                            <tr key={event.id}>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="text-sm font-medium text-gray-900">{event.eventType}</div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                    {event.serviceName}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                    {event.eventsCount}
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                <span className="bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded">
                    {event.language}
                </span>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                    <Link
                                        to={`/events/${event.eventType}`}
                                        className="text-blue-600 hover:text-blue-800"
                                    >
                                        View
                                    </Link>
                                </td>
                            </tr>
                        ))
                    )}
                    </tbody>
                </table>
            </div>
        </div>
    );
}

export default RecentEventsTable;