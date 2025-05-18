import React, { useState, useEffect } from 'react';
import StatusCard from './StatusCard';
import RecentEventsTable from './RecentEventsTable';
import SystemHealth from './SystemHealth';
import aletheiaClient from "../api/aletheiaClient.js";

function Dashboard() {
    const [events, setEvents] = useState([]);

    useEffect(() => {
        async function fetchEvents() {
            try {
                const response = await aletheiaClient.get('events');
                setEvents(response.data.items.events || []);
            } catch (error) {
                console.error('Error fetching events:', error);
                setEvents([]); // В случае ошибки устанавливаем пустой массив
            }
        }
        fetchEvents();
    }, []);

    const totalEvents = events.reduce((sum, event) => sum + (event.eventsCount || 0), 0);
    const activeAlerts = events.length;

    return (
        <section id="dashboard" className="mb-12">
            <header className="mb-6">
                <h1 className="text-2xl font-bold">Dashboard</h1>
                <p className="text-gray-600">Overview of your monitoring systems and events</p>
            </header>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                <StatusCard
                    title="Total Events"
                    value={totalEvents.toString()}
                    change="12.5% from last week" // Заглушка, так как нет исторических данных
                    icon="fas fa-calendar-check"
                    iconBg="bg-blue-100"
                    iconColor="text-blue-600"
                />
                <StatusCard
                    title="Active Alerts"
                    value={activeAlerts.toString()}
                    change="3 new since yesterday" // Заглушка, так как нет исторических данных
                    icon="fas fa-exclamation-circle"
                    iconBg="bg-red-100"
                    iconColor="text-red-600"
                />
                <StatusCard
                    title="System Uptime"
                    value="99.8%"
                    change="All systems operational"
                    icon="fas fa-server"
                    iconBg="bg-green-100"
                    iconColor="text-green-600"
                />
            </div>

            <RecentEventsTable />
            <SystemHealth />
        </section>
    );
}

export default Dashboard;