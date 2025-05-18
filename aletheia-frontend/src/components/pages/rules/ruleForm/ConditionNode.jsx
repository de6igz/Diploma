import React from 'react';
import Condition from './Condition';

function ConditionNode({
                           node,
                           path,
                           onOperatorChange,
                           onConditionChange,
                           onAddCondition,
                           onAddChildNode,
                           onDeleteCondition,
                           onDeleteChildNode,
                       }) {
    return (
        <div className="border border-gray-200 bg-gray-50 p-4 rounded-md mb-4">
            <div className="flex items-center space-x-2 mb-4">
                <label className="font-medium text-sm text-gray-700">Логический оператор:</label>
                <select
                    value={node.operator}
                    onChange={(e) => onOperatorChange(path, e.target.value)}
                    className="p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500"
                >
                    <option value="AND">AND</option>
                    <option value="OR">OR</option>
                </select>
            </div>
            {node.conditions.map((condition, index) => (
                <div key={index} className="mb-4">
                    <Condition
                        condition={condition}
                        onChange={(field, value) => onConditionChange(path, index, field, value)}
                    />
                    <div className="flex justify-end mt-2">
                        <button
                            className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 text-sm"
                            onClick={() => onDeleteCondition(path, index)}
                        >
                            Удалить
                        </button>
                    </div>
                </div>
            ))}
            {node.children.map((child, index) => (
                <div key={index} className="mb-4">
                    <ConditionNode
                        node={child}
                        path={[...path, index]}
                        onOperatorChange={onOperatorChange}
                        onConditionChange={onConditionChange}
                        onAddCondition={onAddCondition}
                        onAddChildNode={onAddChildNode}
                        onDeleteCondition={onDeleteCondition}
                        onDeleteChildNode={onDeleteChildNode}
                    />
                </div>
            ))}
            <div className="flex space-x-2 mt-4">
                <button
                    className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm"
                    onClick={() => onAddCondition(path)}
                >
                    Добавить Условие
                </button>
                <button
                    className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm"
                    onClick={() => onAddChildNode(path)}
                >
                    Добавить Дочерний Узел
                </button>
                {/* Кнопка "Удалить узел" только для дочерних узлов */}
                {path.length > 0 && (
                    <button
                        className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 text-sm"
                        onClick={() => onDeleteChildNode(path.slice(0, -1), path[path.length - 1])}
                    >
                        Удалить узел
                    </button>
                )}
            </div>
        </div>
    );
}

export default ConditionNode;