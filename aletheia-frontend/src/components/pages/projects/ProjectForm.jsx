import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import aletheiaClient from "../api/aletheiaClient.js";

function ProjectForm() {
    const [projectName, setProjectName] = useState('');
    const [description, setDescription] = useState('');
    const [services, setServices] = useState([]);
    const [availableRules, setAvailableRules] = useState([]);
    const [allRules, setAllRules] = useState([]); // Новое состояние
    const [errors, setErrors] = useState({
        projectName: '',
        services: [],
    });
    const { id } = useParams();

    const fetchAvailableRules = async () => {
        try {
            const response = await aletheiaClient.get('/rules/available');
            setAllRules(response.data.items.rules);
        } catch (error) {
            console.error('Ошибка при загрузке доступных правил:', error);
        }
    };

    const fetchProjectInfo = async () => {
        try {
            const response = await aletheiaClient.get(`/project/${id}`);
            const project = response.data.project.project;
            setProjectName(project.projectName || '');

            const projectServices = Array.isArray(project.services) ? project.services : [];
            const formattedServices = projectServices.map((service) => ({
                serviceName: service.serviceName || '',
                errorRules: Array.isArray(service.errorRules)
                    ? service.errorRules.map((rule) => ({ id: rule.ruleId, name: rule.ruleName }))
                    : [],
                resourceRules: Array.isArray(service.resourceRules)
                    ? service.resourceRules.map((rule) => ({ id: rule.ruleId, name: rule.ruleName }))
                    : [],
            }));
            setServices(formattedServices);
        } catch (error) {
            console.error('Ошибка при загрузке проекта:', error);
        }
    };

    useEffect(() => {
        fetchAvailableRules();
        if (id) {
            fetchProjectInfo();
        } else {
            setServices([{ serviceName: '', errorRules: [], resourceRules: [] }]);
        }
    }, [id]);

    useEffect(() => {
        if (allRules.length > 0) {
            const selectedRules = new Set();
            services.forEach((service) => {
                service.errorRules.forEach((rule) => {
                    if (rule.id) selectedRules.add(rule.id);
                });
                service.resourceRules.forEach((rule) => {
                    if (rule.id) selectedRules.add(rule.id);
                });
            });
            setAvailableRules(allRules.filter((rule) => !selectedRules.has(rule.id)));
        }
    }, [services, allRules]);

    // Получение всех выбранных правил
    const getAllSelectedRules = () => {
        const selectedRules = new Set();
        services.forEach((service) => {
            service.errorRules.forEach((rule) => selectedRules.add(rule.id));
            service.resourceRules.forEach((rule) => selectedRules.add(rule.id));
        });
        return selectedRules;
    };

    // Добавление нового сервиса
    const handleAddService = () => {
        setServices([...services, { serviceName: '', errorRules: [], resourceRules: [] }]);
    };

    // Удаление сервиса
    const handleRemoveService = (index) => {
        const removedService = services[index];
        const updatedServices = services.filter((_, i) => i !== index);
        setServices(updatedServices);

        // Возвращаем правила в availableRules, если они больше не используются
        const selectedRules = getAllSelectedRules();
        [...removedService.errorRules, ...removedService.resourceRules].forEach((rule) => {
            if (!selectedRules.has(rule.id)) {
                setAvailableRules((prev) => [...prev, rule]);
            }
        });
    };

    // Изменение названия сервиса
    const handleServiceNameChange = (index, value) => {
        const updatedServices = [...services];
        updatedServices[index].serviceName = value;
        setServices(updatedServices);
    };

    // Добавление нового правила в сервис
    const handleAddRule = (serviceIndex, ruleType) => {
        const updatedServices = [...services];
        updatedServices[serviceIndex][ruleType].push({ id: '', name: '' });
        setServices(updatedServices);
    };

    // Удаление правила из сервиса
    const handleRemoveRule = (serviceIndex, ruleType, ruleIndex) => {
        const updatedServices = [...services];
        const removedRule = updatedServices[serviceIndex][ruleType][ruleIndex];
        updatedServices[serviceIndex][ruleType].splice(ruleIndex, 1);
        setServices(updatedServices);

        // Возвращаем правило в availableRules, если оно больше не используется
        if (removedRule.id) {
            const selectedRules = getAllSelectedRules();
            if (!selectedRules.has(removedRule.id)) {
                setAvailableRules((prev) => [...prev, removedRule]);
            }
        }
    };

    // Изменение значения правила
    const handleRuleChange = (serviceIndex, ruleType, ruleIndex, selectedRuleId) => {
        const updatedServices = [...services];
        const previousRule = updatedServices[serviceIndex][ruleType][ruleIndex];

        if (selectedRuleId) {
            const selectedRule = availableRules.find((rule) => rule.id === selectedRuleId);
            if (selectedRule) {
                updatedServices[serviceIndex][ruleType][ruleIndex] = {
                    id: selectedRule.id,
                    name: selectedRule.name,
                };
                setAvailableRules((prev) => prev.filter((rule) => rule.id !== selectedRuleId));
            }
        } else {
            if (previousRule.id) {
                const selectedRules = getAllSelectedRules();
                selectedRules.delete(previousRule.id);
                if (!selectedRules.has(previousRule.id)) {
                    setAvailableRules((prev) => [...prev, previousRule]);
                }
            }
            updatedServices[serviceIndex][ruleType][ruleIndex] = { id: '', name: '' };
        }
        setServices(updatedServices);
    };

    // Валидация формы
    const validateForm = () => {
        let isValid = true;
        const newErrors = {
            projectName: '',
            services: services.map(() => ({
                serviceName: '',
                errorRules: [],
                resourceRules: [],
            })),
        };

        if (!projectName.trim()) {
            newErrors.projectName = 'Название проекта не может быть пустым';
            isValid = false;
        }

        services.forEach((service, index) => {
            if (!service.serviceName.trim()) {
                newErrors.services[index].serviceName = 'Название сервиса не может быть пустым';
                isValid = false;
            }
            service.errorRules.forEach((rule, ruleIndex) => {
                if (!rule.id) {
                    newErrors.services[index].errorRules[ruleIndex] = 'Выберите правило';
                    isValid = false;
                }
            });
            service.resourceRules.forEach((rule, ruleIndex) => {
                if (!rule.id) {
                    newErrors.services[index].resourceRules[ruleIndex] = 'Выберите правило';
                    isValid = false;
                }
            });
        });

        setErrors(newErrors);
        return isValid;
    };

    // Сохранение проекта
    const handleSaveProject = async () => {
        if (validateForm()) {
            const projectData = {
                project: {
                    projectName: projectName,
                    description: description,
                    services: services.map((service) => ({
                        serviceName: service.serviceName,
                        errorRules: service.errorRules.map((rule) => parseInt(rule.id)),
                        resourceRules: service.resourceRules.map((rule) => parseInt(rule.id)),
                    })),
                },
            };

            try {
                console.log(JSON.stringify(projectData));
                if (id) {
                    // TODO: Добавить ручку для обновления
                    console.log(id)
                    await aletheiaClient.put(`/project/${id}`, projectData);

                } else {
                    await aletheiaClient.post('/project/create', projectData);
                }
                alert('Проект успешно сохранен!');
            } catch (error) {
                console.error('Ошибка при сохранении проекта:', error);
                alert('Ошибка при сохранении проекта');
            }
        } else {
            alert('Пожалуйста, исправьте ошибки в форме');
        }
    };

    return (
        <div className="max-w-4xl mx-auto my-10 bg-white rounded-lg shadow-md p-8">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold text-gray-800">
                    {id ? 'Редактирование проекта' : 'Создание проекта'}
                </h1>
                <div className="flex space-x-3">
                    <button
                        className="px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 text-sm"
                        onClick={() => window.history.back()}
                    >
                        Отмена
                    </button>
                    <button
                        className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 text-sm"
                        onClick={handleSaveProject}
                    >
                        Сохранить проект
                    </button>
                </div>
            </div>

            {/* Название проекта */}
            <div className="mb-8 space-y-4">
                <div>
                    <label className="block font-medium text-sm text-gray-700 mb-1">
                        Название проекта
                    </label>
                    <input
                        type="text"
                        value={projectName}
                        onChange={(e) => setProjectName(e.target.value)}
                        placeholder="Введите название проекта"
                        className={`placeholder:text-gray-400 w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-yellow-500 ${errors.projectName ? 'border-red-500' : ''}`}
                    />
                    {errors.projectName && (
                        <p className="text-red-500 text-xs mt-1">{errors.projectName}</p>
                    )}
                </div>
            </div>

            {/* Описание проекта */}
            <div className="mb-8 space-y-4">
                <div>
                    <label className="block font-medium text-sm text-gray-700 mb-1">
                        Описание проекта
                    </label>
                    <input
                        type="text"
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        placeholder="Введите описание проекта"
                        className="placeholder:text-gray-400 w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-yellow-500"
                    />
                </div>
            </div>

            {/* Сервисы */}
            <div className="mb-8">
                <div className="flex justify-between items-center mb-3">
                    <h2 className="text-lg font-semibold text-gray-800">Сервисы</h2>
                    <button
                        className="px-4 py-2 hover:bg-[#2044B8] hover:text-[#DBE9FE] bg-[#DBE9FE] text-[#2044B8] rounded-md text-sm"
                        onClick={handleAddService}
                    >
                        Добавить сервис
                    </button>
                </div>

                {services.length === 0 && (
                    <p className="text-sm text-gray-500">Добавьте сервисы для проекта.</p>
                )}

                {services.map((service, serviceIndex) => (
                    <div
                        key={serviceIndex}
                        className="border border-gray-200 p-4 rounded-md mb-4"
                    >
                        <div className="flex justify-between items-center mb-3">
                            <div className="font-medium text-gray-800">
                                Сервис {serviceIndex + 1}
                            </div>
                            <button
                                className="px-3 py-1 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 text-sm"
                                onClick={() => handleRemoveService(serviceIndex)}
                            >
                                Удалить
                            </button>
                        </div>

                        <div className="mb-3">
                            <label className="block font-medium text-sm text-gray-700 mb-1">
                                Название сервиса
                            </label>
                            <input
                                type="text"
                                value={service.serviceName}
                                onChange={(e) => handleServiceNameChange(serviceIndex, e.target.value)}
                                placeholder="Введите название сервиса"
                                className={`placeholder:text-gray-400 w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-yellow-500 ${errors.services[serviceIndex]?.serviceName ? 'border-red-500' : ''}`}
                            />
                            {errors.services[serviceIndex]?.serviceName && (
                                <p className="text-red-500 text-xs mt-1">{errors.services[serviceIndex].serviceName}</p>
                            )}
                        </div>

                        {/* Error Rules */}
                        <div className="mb-4">
                            <div className="flex justify-between items-center mb-2">
                                <label className="text-sm font-medium text-gray-700">
                                    Error Rules
                                </label>
                                <button
                                    className="text-sm text-blue-600 hover:text-blue-800 flex items-center"
                                    onClick={() => handleAddRule(serviceIndex, 'errorRules')}
                                >
                                    <i className="fas fa-plus mr-1"></i> Добавить правило
                                </button>
                            </div>
                            {service.errorRules.map((rule, ruleIndex) => (
                                <div
                                    key={ruleIndex}
                                    className="bg-gray-50 p-3 rounded-md mb-3"
                                >
                                    <div className="flex justify-between items-center mb-2">
                                        <div className="font-medium text-gray-800">
                                            Правило {ruleIndex + 1}
                                        </div>
                                        <button
                                            className="text-gray-500 hover:text-red-600"
                                            onClick={() => handleRemoveRule(serviceIndex, 'errorRules', ruleIndex)}
                                        >
                                            <i className="fas fa-times"></i>
                                        </button>
                                    </div>
                                    <select
                                        value={rule.id}
                                        onChange={(e) => handleRuleChange(serviceIndex, 'errorRules', ruleIndex, e.target.value)}
                                        className={`w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-yellow-500 ${errors.services[serviceIndex]?.errorRules[ruleIndex] ? 'border-red-500' : ''}`}
                                    >
                                        <option value="">Выберите правило</option>
                                        {availableRules
                                            .filter((r) => r.ruleType === 'errors')
                                            .map((r) => (
                                                <option key={r.id} value={r.id}>
                                                    {r.name}
                                                </option>
                                            ))}
                                        {rule.id && !availableRules.some((r) => r.id === rule.id) && (
                                            <option value={rule.id}>{rule.name}</option>
                                        )}
                                    </select>
                                    {errors.services[serviceIndex]?.errorRules[ruleIndex] && (
                                        <p className="text-red-500 text-xs mt-1">{errors.services[serviceIndex].errorRules[ruleIndex]}</p>
                                    )}
                                </div>
                            ))}
                        </div>

                        {/* Resource Rules */}
                        <div className="mb-4">
                            <div className="flex justify-between items-center mb-2">
                                <label className="text-sm font-medium text-gray-700">
                                    Resource Rules
                                </label>
                                <button
                                    className="text-sm text-blue-600 hover:text-blue-800 flex items-center"
                                    onClick={() => handleAddRule(serviceIndex, 'resourceRules')}
                                >
                                    <i className="fas fa-plus mr-1"></i> Добавить правило
                                </button>
                            </div>
                            {service.resourceRules.map((rule, ruleIndex) => (
                                <div
                                    key={ruleIndex}
                                    className="bg-gray-50 p-3 rounded-md mb-3"
                                >
                                    <div className="flex justify-between items-center mb-2">
                                        <div className="font-medium text-gray-800">
                                            Правило {ruleIndex + 1}
                                        </div>
                                        <button
                                            className="text-gray-500 hover:text-red-600"
                                            onClick={() => handleRemoveRule(serviceIndex, 'resourceRules', ruleIndex)}
                                        >
                                            <i className="fas fa-times"></i>
                                        </button>
                                    </div>
                                    <select
                                        value={rule.id}
                                        onChange={(e) => handleRuleChange(serviceIndex, 'resourceRules', ruleIndex, e.target.value)}
                                        className={`w-full p-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-yellow-500 ${errors.services[serviceIndex]?.resourceRules[ruleIndex] ? 'border-red-500' : ''}`}
                                    >
                                        <option value="">Выберите правило</option>
                                        {availableRules
                                            .filter((r) => r.ruleType === 'resources')
                                            .map((r) => (
                                                <option key={r.id} value={r.id}>
                                                    {r.name}
                                                </option>
                                            ))}
                                        {rule.id && !availableRules.some((r) => r.id === rule.id) && (
                                            <option value={rule.id}>{rule.name}</option>
                                        )}
                                    </select>
                                    {errors.services[serviceIndex]?.resourceRules[ruleIndex] && (
                                        <p className="text-red-500 text-xs mt-1">{errors.services[serviceIndex].resourceRules[ruleIndex]}</p>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}

export default ProjectForm;