import React, {useState, useEffect} from 'react';
import {useNavigate, useParams} from 'react-router-dom';
import ConditionNode from './ConditionNode';
import Action from './Action';
import RuleSettings from './RuleSettings';
import aletheiaClient from "../../api/aletheiaClient.js";

function RuleForm() {
    const navigate = useNavigate();
    const [rule, setRule] = useState({
        name: '',
        description: '',
        ruleType: 'errors',
        rootNode: {operator: 'AND', conditions: [], children: []},
        actions: [],
        settings: {priority: 'medium', status: 'active'},
    });
    const [errors, setErrors] = useState({});
    const {id, ruleType} = useParams();

    // Загрузка данных правила при редактировании
    const fetchRuleData = async () => {
        try {
            console.log(id, ruleType)
            const response = await aletheiaClient.get(`rule/byId?ruleId=${id}&ruleType=${ruleType}`, {
                request: {
                    ruleId: id,
                    ruleType: ruleType,
                },
            });
            const ruleData = response.data.rule;
            console.log(ruleData)
            setRule({
                ...ruleData,
                rootNode: ruleData.root_node, // Приведение к внутренней структуре
                ruleType: ruleType,
                actions: ruleData.actions.map((action) => ({
                    type: action.type,
                    params: {
                        key: action.params.key || '',
                        value: action.params.value || '',
                    },
                })),
            });
        } catch (error) {
            console.error('Ошибка при загрузке правила:', error);
        }
    };
    console.log(rule.ruleType)

    useEffect(() => {
        if (id) {
            fetchRuleData();
        }
    }, [id]);

    // Валидация формы
    const validateForm = () => {
        const newErrors = {};
        if (!rule.name.trim()) {
            newErrors.name = 'Название правила не может быть пустым';
        }
        // Валидация условий
        const validateConditions = (node) => {
            node.conditions.forEach((condition) => {
                if (!condition.field || !condition.operator || !condition.value) {
                    newErrors.conditions = 'Все поля условия должны быть заполнены';
                }
            });
            node.children.forEach((child) => validateConditions(child));
        };
        validateConditions(rule.rootNode);
        // Валидация действий
        rule.actions.forEach((action, index) => {
            if (!action.type) {
                newErrors[`actionType${index}`] = 'Тип действия не может быть пустым';
            }
            if (!action.params.key || !action.params.value) {
                newErrors[`actionParams${index}`] = 'Параметры действия не могут быть пустыми';
            }
        });
        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    // Обработчик сохранения правила
    const handleSaveRule = async () => {
        if (validateForm()) {
            const ruleData = {
                request: {
                    actions: rule.actions.map(action => ({
                        type: action.type,
                        params: {
                            key: action.params.key,
                            value: action.params.value
                        },
                    })),
                    description: rule.description,
                    name: rule.name,
                    root_node: rule.rootNode,
                    ruleType: rule.ruleType,
                    ...(id && {ruleId: id}), // Добавляем ruleId для обновления
                },
            };
            console.log(ruleData)
            try {
                if (id) {
                    await aletheiaClient.put('rules/update', ruleData);
                    alert('Правило успешно обновлено!');
                    navigate("/rules");
                } else {
                    await aletheiaClient.post('rules/create', ruleData);
                    alert('Правило успешно создано!');
                    navigate("/rules");
                }
            } catch (error) {
                console.error('Ошибка при сохранении правила:', error);
                alert('Ошибка при сохранении правила');
            }
        } else {
            alert('Пожалуйста, исправьте ошибки в форме');
        }
    };

    // Обработчики для условий и действий
    const handleAddCondition = (path) => {
        setRule((prevRule) => ({
            ...prevRule,
            rootNode: addConditionToNode(prevRule.rootNode, path),
        }));
    };

    const handleAddChildNode = (path) => {
        setRule((prevRule) => ({
            ...prevRule,
            rootNode: addChildNodeToNode(prevRule.rootNode, path),
        }));
    };

    const handleDeleteCondition = (path, conditionIndex) => {
        setRule((prevRule) => ({
            ...prevRule,
            rootNode: deleteConditionFromNode(prevRule.rootNode, path, conditionIndex),
        }));
    };

    const handleDeleteChildNode = (path, childIndex) => {
        setRule((prevRule) => ({
            ...prevRule,
            rootNode: deleteChildNodeFromNode(prevRule.rootNode, path, childIndex),
        }));
    };

    const handleConditionChange = (path, conditionIndex, field, value) => {
        setRule((prevRule) => ({
            ...prevRule,
            rootNode: changeConditionInNode(prevRule.rootNode, path, conditionIndex, field, value),
        }));
    };

    const handleOperatorChange = (path, value) => {
        setRule((prevRule) => ({
            ...prevRule,
            rootNode: changeOperatorInNode(prevRule.rootNode, path, value),
        }));
    };

    const handleAddAction = () => {
        setRule({
            ...rule,
            actions: [...rule.actions, {type: '', params: {key: '', value: ''}}],
        });
    };

    const handleDeleteAction = (index) => {
        setRule({
            ...rule,
            actions: rule.actions.filter((_, i) => i !== index),
        });
    };

    const handleActionChange = (index, field, value) => {
        const updatedActions = rule.actions.map((action, i) => {
            if (i === index) {
                if (field === 'type') return {...action, type: value};
                return {...action, params: {...action.params, [field]: value}};
            }
            return action;
        });
        setRule({...rule, actions: updatedActions});
    };

    const handleSettingsChange = (field, value) => {
        setRule({...rule, settings: {...rule.settings, [field]: value}});
    };

    const generateJSON = () => {
        console.log(JSON.stringify(rule, null, 2));
        document.getElementById('output').textContent = JSON.stringify(rule, null, 4);
    };

    // Вспомогательные функции для работы с узлами
    const addConditionToNode = (node, path) => {
        if (path.length === 0) {
            return {
                ...node,
                conditions: [...node.conditions, {field: '', operator: 'eq', value: ''}],
            };
        }
        const [index, ...restPath] = path;
        return {
            ...node,
            children: node.children.map((child, i) =>
                i === index ? addConditionToNode(child, restPath) : child
            ),
        };
    };

    const addChildNodeToNode = (node, path) => {
        if (path.length === 0) {
            return {
                ...node,
                children: [...node.children, {operator: 'AND', conditions: [], children: []}],
            };
        }
        const [index, ...restPath] = path;
        return {
            ...node,
            children: node.children.map((child, i) =>
                i === index ? addChildNodeToNode(child, restPath) : child
            ),
        };
    };

    const deleteConditionFromNode = (node, path, conditionIndex) => {
        if (path.length === 0) {
            return {
                ...node,
                conditions: node.conditions.filter((_, i) => i !== conditionIndex),
            };
        }
        const [index, ...restPath] = path;
        return {
            ...node,
            children: node.children.map((child, i) =>
                i === index ? deleteConditionFromNode(child, restPath, conditionIndex) : child
            ),
        };
    };

    const deleteChildNodeFromNode = (node, path, childIndex) => {
        if (path.length === 0) {
            return {
                ...node,
                children: node.children.filter((_, i) => i !== childIndex),
            };
        }
        const [index, ...restPath] = path;
        return {
            ...node,
            children: node.children.map((child, i) =>
                i === index ? deleteChildNodeFromNode(child, restPath, childIndex) : child
            ),
        };
    };

    const changeConditionInNode = (node, path, conditionIndex, field, value) => {
        if (path.length === 0) {
            return {
                ...node,
                conditions: node.conditions.map((condition, i) =>
                    i === conditionIndex ? {...condition, [field]: value} : condition
                ),
            };
        }
        const [index, ...restPath] = path;
        return {
            ...node,
            children: node.children.map((child, i) =>
                i === index ? changeConditionInNode(child, restPath, conditionIndex, field, value) : child
            ),
        };
    };

    const changeOperatorInNode = (node, path, value) => {
        if (path.length === 0) {
            return {...node, operator: value};
        }
        const [index, ...restPath] = path;
        return {
            ...node,
            children: node.children.map((child, i) =>
                i === index ? changeOperatorInNode(child, restPath, value) : child
            ),
        };
    };

    return (
        <div className="max-w-4xl mx-auto my-10 bg-white rounded-lg shadow-md p-8">
            <h1 className="text-2xl font-bold mb-6 text-gray-800">
                {id ? 'Редактировать правило' : 'Создать новое правило'}
            </h1>
            <div className="mb-8 space-y-4">
                <div>
                    <label className="block font-medium text-sm text-gray-700 mb-1">Название правила</label>
                    <input
                        type="text"
                        value={rule.name}
                        onChange={(e) => setRule({...rule, name: e.target.value})}
                        placeholder="Введите название правила"
                        className={`placeholder:text-gray-400 w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-blue-500 ${errors.name ? 'border-red-500' : ''}`}
                    />
                    {errors.name && <p className="text-red-500 text-xs mt-1">{errors.name}</p>}
                </div>
                <div>
                    <label className="block font-medium text-sm text-gray-700 mb-1">Описание правила</label>
                    <textarea
                        value={rule.description}
                        onChange={(e) => setRule({...rule, description: e.target.value})}
                        placeholder="Краткое описание правила"
                        className="placeholder:text-gray-400 w-full p-2 border border-gray-300 rounded-md text-sm min-h-[100px] focus:ring-2 focus:ring-blue-500"
                    />
                </div>
                <div>
                    <label className="block font-medium text-sm text-gray-700 mb-1">Тип правила</label>
                    <div className="flex space-x-4">
                        <label className="flex items-center">
                            <input
                                type="radio"
                                name="ruleType"
                                value="errors"
                                checked={rule.ruleType === 'errors'}
                                onChange={() => setRule({...rule, ruleType: 'errors'})}
                                className="mr-2 text-blue-600 focus:ring-blue-500"
                            />
                            <span className="text-sm text-gray-700">Ошибки</span>
                        </label>
                        <label className="flex items-center">
                            <input
                                type="radio"
                                name="ruleType"
                                value="resources"
                                checked={rule.ruleType === 'resources'}
                                onChange={() => setRule({...rule, ruleType: 'resources'})}
                                className="mr-2 text-blue-600 focus:ring-blue-500"
                            />
                            <span className="text-sm text-gray-700">Ресурсы</span>
                        </label>
                    </div>
                </div>
            </div>
            <div className="mb-8">
                <h2 className="text-lg font-semibold mb-2 pb-1 border-b border-gray-200 text-gray-800 mb-3">Условия</h2>
                <p className="mb-3">Настройте условия с помощью оператора <strong>AND</strong> или <strong>OR</strong>.
                    Вы можете добавлять дочерние узлы для более сложной логики.</p>
                <ConditionNode
                    node={rule.rootNode}
                    path={[]}
                    onOperatorChange={handleOperatorChange}
                    onConditionChange={handleConditionChange}
                    onAddCondition={handleAddCondition}
                    onAddChildNode={handleAddChildNode}
                    onDeleteCondition={handleDeleteCondition}
                    onDeleteChildNode={handleDeleteChildNode}
                />
                {errors.conditions && <p className="text-red-500 text-xs mt-1">{errors.conditions}</p>}
            </div>
            <div className="mb-8">
                <h2 className="text-lg font-semibold mb-2 pb-1 border-b border-gray-200 text-gray-800 mb-3">Действия</h2>
                <p>Укажите, какие действия выполнять при выполнении условий.</p>
                {rule.actions.map((action, index) => (
                    <Action
                        key={index}
                        action={action}
                        onChange={(field, value) => handleActionChange(index, field, value)}
                        onDelete={() => handleDeleteAction(index)}
                        errors={errors}
                        index={index}
                    />
                ))}
                <button
                    className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm"
                    onClick={handleAddAction}
                >
                    Добавить Действие
                </button>
            </div>

            <div className="flex space-x-3 mb-8">
                <button
                    className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 text-sm"
                    onClick={generateJSON}
                >
                    Сгенерировать JSON
                </button>
                <button
                    className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm"
                    onClick={handleSaveRule}
                >
                    Сохранить правило
                </button>
            </div>
            <div className="bg-gray-50 p-4 rounded-md">
                <pre id="output" className="text-xs text-gray-800 whitespace-pre-wrap"></pre>
            </div>
        </div>
    );
}

export default RuleForm;