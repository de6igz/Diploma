import React from 'react';

function EventConfigurationForm() {
    return (
        <section id="events" className="mb-12">
            <header className="mb-6">
                <h1 className="text-2xl font-bold">Event Type Configuration</h1>
                <p className="text-gray-600">Create and manage event types to be monitored</p>
            </header>

            <div className="bg-white rounded-xl shadow-sm">
                <div className="p-6 border-b border-gray-100">
                    <h2 className="text-lg font-semibold">New Event Type</h2>
                </div>
                <div className="p-6">
                    <form>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                            <div>
                                <label
                                    htmlFor="event-name"
                                    className="block text-sm font-medium text-gray-700 mb-1"
                                >
                                    Name
                                </label>
                                <input
                                    type="text"
                                    id="event-name"
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    placeholder="API Latency Alert"
                                />
                            </div>
                            <div>
                                <label
                                    htmlFor="event-source"
                                    className="block text-sm font-medium text-gray-700 mb-1"
                                >
                                    Source System
                                </label>
                                <select
                                    id="event-source"
                                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                >
                                    <option>Production API</option>
                                    <option>Database Cluster</option>
                                    <option>Authentication Service</option>
                                    <option>Payment Processing</option>
                                    <option>Add new source...</option>
                                </select>
                            </div>
                        </div>

                        <div className="mt-8">
                            <h3 className="text-md font-semibold mb-4">Notification Settings</h3>
                            <div className="space-y-3">
                                <div className="flex items-center">
                                    <input
                                        id="notify-email"
                                        type="checkbox"
                                        className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                                    />
                                    <label
                                        htmlFor="notify-email"
                                        className="ml-2 block text-sm text-gray-700"
                                    >
                                        Email notifications
                                    </label>
                                </div>
                                <div className="flex items-center">
                                    <input
                                        id="notify-sms"
                                        type="checkbox"
                                        className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                                    />
                                    <label htmlFor="notify-sms" className="ml-2 block text-sm text-gray-700">
                                        SMS notifications
                                    </label>
                                </div>
                                <div className="flex items-center">
                                    <input
                                        id="notify-webhook"
                                        type="checkbox"
                                        className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                                    />
                                    <label
                                        htmlFor="notify-webhook"
                                        className="ml-2 block text-sm text-gray-700"
                                    >
                                        Webhook notifications
                                    </label>
                                </div>
                            </div>
                        </div>

                        <div className="mt-8 flex justify-end">
                            <button
                                type="button"
                                className="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 mr-3"
                            >
                                Cancel
                            </button>
                            <button
                                type="submit"
                                className="px-4 py-2 bg-yellow-400 hover:bg-yellow-500 text-yellow-900 rounded-md text-sm font-medium"
                            >
                                Create Event Type
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </section>
    );
}

export default EventConfigurationForm;