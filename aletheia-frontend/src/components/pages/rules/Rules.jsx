import React, { useState, useEffect } from 'react';
import { useNavigate } from "react-router-dom";
import aletheiaClient from "../api/aletheiaClient.js"; // Ensure this is correctly imported

function RuleCard({ rule, onDelete }) {
    const navigate = useNavigate();
    return (
        <div className="bg-white rounded-xl shadow-sm p-4">
            <div className="flex justify-between items-center">
                <div>
                    <h3 className="text-lg font-semibold">
                        {rule.name}
                        <span className="text-sm text-gray-500 ml-2">RuleType:({rule.ruleType})</span>
                    </h3>


                    <p className="text-gray-600">{rule.description || 'No description'}</p>
                </div>
                <div className="flex space-x-2">
                    <button
                        className="flex items-center space-x-1 hover:bg-[#2044B8] hover:text-[#DBE9FE] bg-[#DBE9FE] text-[#2044B8] px-3 py-1 rounded "
                        onClick={() => navigate(`/rules/${rule.ruleType}/${rule.id}/edit`)}
                    >
                        <i className="fas fa-check"></i>
                        <span>Edit</span>
                    </button>
                    <button
                        className="flex items-center space-x-1 px-3 py-1 rounded hover:bg-[#DC2625] hover:text-[#FEE2E1] bg-[#FEE2E1] text-[#DC2625]"
                        onClick={() => onDelete(rule)}
                    >
                        <i className="fas fa-trash " ></i>
                        <span>delete</span>
                    </button>
                </div>
            </div>
        </div>
    );
}
//TODO: щас только errorType: rules приходят, надо будет проверить после обновления ручки
function RulesPage() {
    const navigate = useNavigate();
    const [rules, setRules] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState(null);
    // Fetch rules from the API
    const fetchRules = async () => {
        try {
            const response = await aletheiaClient.get('/rules');
            setRules(response.data.items.rules); // Update state with API data
            console.log(response.data.items)
            setIsLoading(false);
        } catch (err) {
            console.error('Error fetching rules:', err);
            setError('Failed to load rules');
            setIsLoading(false);
        }
    };

    // Load rules when the component mounts
    useEffect(() => {
        fetchRules();
    }, []);

    // Handle rule deletion
    const handleDeleteRule = async (rule) => {
        try {
            await aletheiaClient.delete('/rules', {
                data: {
                    req: {
                        ruleId: rule.id,
                        ruleType: rule.ruleType,
                    },
                },
            });
            setRules(rules.filter(r => r.id !== rule.id));
        } catch (err) {
            console.error('Error deleting rule:', err);
            alert('Failed to delete rule');
        }
    };


    // Show loading state
    if (isLoading) {
        return <div>Loading...</div>;
    }

    // Show error state
    if (error) {
        return <div className="text-red-500">{error}</div>;
    }

    return (
        <section id="rules" className="mb-12">
            <header className="mb-6 flex justify-between items-center">
                <span>
                    <h1 className="text-2xl font-bold">Rules</h1>
                    <p className="text-gray-600">Create and edit project rules</p>
                </span>
                <div className="flex space-x-4">
                    <button className="hover:bg-[#2044B8] hover:text-[#DBE9FE] bg-[#DBE9FE] text-[#2044B8] px-4 py-2 rounded ">
                        Фильтр по проектам
                    </button>
                    <button
                        className="hover:bg-[#2044B8] hover:text-[#DBE9FE] bg-[#DBE9FE] text-[#2044B8] px-4 py-2 rounded "
                        onClick={() => navigate('/rules/new')}
                    >
                        + Создать правило
                    </button>
                </div>
            </header>
            <div className="space-y-4">
                {rules ? (rules.map((rule) => (
                    <RuleCard key={rule.id} rule={rule} onDelete={handleDeleteRule} />
                ))) : <p>У вас пока нет правил</p>}
            </div>
        </section>
    );
}

export default RulesPage;