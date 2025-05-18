import React from 'react';

const actionIcons = {
    TELEGRAM: 'fab fa-telegram text-blue-500',
    DISCORD: 'fab fa-discord text-indigo-500',
    EMAIL: 'fas fa-envelope text-orange-500',
};

function Action({ action, onChange, onDelete, errors, index }) {
    return (
        <div className="bg-gray-100 p-4 rounded-md mb-4">
            <div className="flex justify-between items-center mb-2">
                <label className="text-sm font-medium text-gray-700">Тип действия</label>
                <button
                    className="text-gray-500 hover:text-red-600"
                    onClick={onDelete}
                >
                    <i className="fas fa-times"></i>
                </button>
            </div>
            <select
                value={action.type}
                onChange={(e) => onChange('type', e.target.value)}
                className={`w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 ${errors[`actionType${index}`] ? 'border-red-500' : ''}`}
            >
                <option value="">Выберите тип действия</option>
                <option value="TELEGRAM">Telegram</option>
                <option value="DISCORD">Discord</option>
                <option value="EMAIL">Email</option>
            </select>
            {errors[`actionType${index}`] && <p className="text-red-500 text-xs mt-1">{errors[`actionType${index}`]}</p>}
            {action.type && (
                <div className="mt-2">
                    <i className={actionIcons[action.type]}></i>
                </div>
            )}
            <div className="mt-2">
                <label className="block text-sm font-medium text-gray-700">Ключ параметра</label>
                <input
                    type="text"
                    value={action.params.key}
                    onChange={(e) => onChange('key', e.target.value)}
                    placeholder="Введите ключ (например, chat_id)"
                    className={`w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 ${errors[`actionParams${index}`] ? 'border-red-500' : ''}`}
                />
            </div>
            <div className="mt-2">
                <label className="block text-sm font-medium text-gray-700">Значение параметра</label>
                <input
                    type="text"
                    value={action.params.value}
                    onChange={(e) => onChange('value', e.target.value)}
                    placeholder="Введите значение (например, 123456)"
                    className={`w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 ${errors[`actionParams${index}`] ? 'border-red-500' : ''}`}
                />
            </div>
            {errors[`actionParams${index}`] && <p className="text-red-500 text-xs mt-1">{errors[`actionParams${index}`]}</p>}
        </div>
    );
}

export default Action;