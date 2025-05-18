import React from 'react';

function SystemHealth() {
    return (
        <div className="bg-white rounded-xl shadow-sm">
            <div className="p-6 border-b border-gray-100">
                <h2 className="text-lg font-semibold">System Health</h2>
            </div>
            <div className="p-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div>
                        <h3 className="text-sm font-medium text-gray-500 mb-2">CPU Usage</h3>
                        <div className="h-4 bg-gray-200 rounded-full overflow-hidden">
                            <div className="h-full bg-blue-600 rounded-full" style={{width: '45%'}}></div>
                        </div>
                        <div className="flex justify-between mt-1 text-xs text-gray-500">
                            <span>45%</span>
                            <span>8 cores / 3.5 GHz</span>
                        </div>
                    </div>
                    <div>
                        <h3 className="text-sm font-medium text-gray-500 mb-2">Memory Usage</h3>
                        <div className="h-4 bg-gray-200 rounded-full overflow-hidden">
                            <div className="h-full bg-green-600 rounded-full" style={{width: '28%'}}></div>
                        </div>
                        <div className="flex justify-between mt-1 text-xs text-gray-500">
                            <span>28%</span>
                            <span>8 GB / 16 GB</span>
                        </div>
                    </div>
                    <div>
                        <h3 className="text-sm font-medium text-gray-500 mb-2">Disk Usage</h3>
                        <div className="h-4 bg-gray-200 rounded-full overflow-hidden">
                            <div className="h-full bg-yellow-500 rounded-full" style={{width: '72%'}}></div>
                        </div>
                        <div className="flex justify-between mt-1 text-xs text-gray-500">
                            <span>72%</span>
                            <span>320 GB / 450 GB</span>
                        </div>
                    </div>
                    <div>
                        <h3 className="text-sm font-medium text-gray-500 mb-2">Network Traffic</h3>
                        <div className="h-4 bg-gray-200 rounded-full overflow-hidden">
                            <div className="h-full bg-purple-600 rounded-full" style={{width: '65%'}}></div>
                        </div>
                        <div className="flex justify-between mt-1 text-xs text-gray-500">
                            <span>65%</span>
                            <span>280 Mbps / 1 Gbps</span>
                        </div>
                    </div>
                </div>

            </div>
        </div>
    );
}

export default SystemHealth;