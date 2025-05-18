import React from 'react';


function Condition({ condition, onChange }) {
    // Вспомогательная функция для обработки изменения одного из полей условия.
    const handleFieldChange = (fieldName, value) => {
        onChange(fieldName, value);
    };

    return (
        <div className="border border-dashed border-gray-300 bg-white p-4 rounded-md">
            <div className="space-y-4">
                <div>
                    <label className="block font-medium text-sm text-gray-700 mb-1">
                        Поле:
                    </label>
                    <input
                        type="text"
                        value={condition.field}
                        onChange={(e) => handleFieldChange('field', e.target.value)}
                        placeholder="Например: node_name"
                        className="placeholder:text-gray-400 w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
                    />
                </div>
                <div>
                    <label className="block font-medium text-sm text-gray-700 mb-1">
                        Оператор:
                    </label>
                    <select
                        value={condition.operator}
                        onChange={(e) => handleFieldChange('operator', e.target.value)}
                        className="w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
                    >
                        <option value="eq">=</option>
                        <option value="neq">!=</option>
                        <option value="gt">&gt;</option>
                        <option value="lt">&lt;</option>
                        <option value="repeat_over">repeat_over</option>
                    </select>
                </div>
                {condition.operator === 'repeat_over' ? (
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block font-medium text-sm text-gray-700 mb-1">
                                Порог:
                            </label>
                            <input
                                type="number"
                                value={(condition.value && condition.value.threshold) || ''}
                                onChange={(e) => {
                                    const threshold = parseInt(e.target.value, 10) || 0;
                                    const current =
                                        condition.value && typeof condition.value === 'object'
                                            ? condition.value
                                            : { minutes: 0, threshold: 0 };
                                    handleFieldChange('value', { ...current, threshold });
                                }}
                                className="placeholder:text-gray-400 w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                        <div>
                            <label className="block font-medium text-sm text-gray-700 mb-1">
                                Минуты:
                            </label>
                            <input
                                type="number"
                                value={(condition.value && condition.value.minutes) || ''}
                                onChange={(e) => {
                                    const minutes = parseInt(e.target.value, 10) || 0;
                                    const current =
                                        condition.value && typeof condition.value === 'object'
                                            ? condition.value
                                            : { minutes: 0, threshold: 0 };
                                    handleFieldChange('value', { ...current, minutes });
                                }}
                                className="placeholder:text-gray-400 w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                    </div>
                ) : (
                    <div>
                        <label className="block font-medium text-sm text-gray-700 mb-1">
                            Значение:
                        </label>
                        <input
                            type="text"
                            value={condition.value || ''}
                            onChange={(e) => handleFieldChange('value', e.target.value)}
                            placeholder="Например: example"
                            className="placeholder:text-gray-400 w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
                        />
                    </div>
                )}
            </div>
        </div>
    );
}

export default Condition;


