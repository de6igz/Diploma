import React from 'react';

function RuleSettings({ settings, onChange }) {
    return (
        <div className="space-y-4">
            <div>
                <label className="block font-medium text-sm text-gray-700 mb-1">Приоритет:</label>
                <select
                    value={settings.priority}
                    onChange={(e) => onChange('priority', e.target.value)}
                    className="w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
                >
                    <option value="low">Низкий</option>
                    <option value="medium">Средний</option>
                    <option value="high">Высокий</option>
                </select>
            </div>
            <div>
                <label className="block font-medium text-sm text-gray-700 mb-1">Статус:</label>
                <select
                    value={settings.status}
                    onChange={(e) => onChange('status', e.target.value)}
                    className="w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
                >
                    <option value="active">Активно</option>
                    <option value="inactive">Неактивно</option>
                </select>
            </div>
        </div>
    );
}

export default RuleSettings;