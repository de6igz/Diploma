import React from 'react';

function StatusCard({ title, value, change, icon, iconBg, iconColor }) {
    return (
        <div className="bg-white rounded-xl shadow-sm p-6">
            <div className="flex justify-between items-start">
                <div>
                    <p className="text-gray-500 text-sm">{title}</p>
                    <h3 className="text-3xl font-bold mt-1">{value}</h3>
                    <p className="text-green-600 flex items-center mt-2">
                        <i className="fas fa-arrow-up mr-1 text-xs"></i>
                        <span>{change}</span>
                    </p>
                </div>
                <div className={`${iconBg} p-3 rounded-lg`}>
                    <i className={`${icon} ${iconColor} text-xl`}></i>
                </div>
            </div>
        </div>
    );
}

export default StatusCard;