import React, {useEffect, useState} from 'react';
import {useNavigate} from "react-router-dom";
import aletheiaClient from "../api/aletheiaClient.js";

function ProjectCard({project, onDelete}) {
    const navigate = useNavigate();
    return (
        <div className="bg-white rounded-xl shadow-sm p-4 hover:scale-101 transition-transform" >
            <div className="flex justify-between items-center">
                <div className="cursor-pointer" onClick={() => navigate(`/projects/${project.id}`)}>
                    <h3 className="text-lg font-semibold">
                        {project.projectName}, <span className="text-gray-500">ID: {project.id}</span>
                    </h3>                    <p className="text-gray-600">{project.description}</p>
                </div>
                <div className="flex space-x-2">
                    <button
                        className="flex items-center space-x-1 hover:bg-[#2044B8] hover:text-[#DBE9FE] bg-[#DBE9FE] text-[#2044B8] px-3 py-1 rounded "
                        onClick={() => navigate(`/projects/${project.id}/edit`)}
                    >
                        <i className="fas fa-check"></i>
                        <span>Edit</span>
                    </button>
                    <button
                        className="flex items-center space-x-1 hover:bg-[#DC2625] hover:text-[#FEE2E1] bg-[#FEE2E1] text-[#DC2625] px-3 py-1 rounded"
                        onClick={() => onDelete(project.id)}
                    >
                        <i className="fas fa-trash"></i>
                        <span>delete</span>
                    </button>
                </div>
            </div>
        </div>
    );
}

// страница со всеми проектами
function ProjectsPage() {

    const navigate = useNavigate();

    const [projects, setProjects] = useState([

    ]);

    const handleDeleteProjects = (ruleId) => {
        //fetch удалить с бебебека
        try{
            aletheiaClient.delete(`/project/${ruleId}`);
            setProjects(projects.filter(rule => rule.id !== ruleId));
        }
        catch(error) {
            console.error('Error deleting projects:', error.message);
        }
    };

    const fetchProjects = async () => {
        const response = await aletheiaClient.get("/projects")
        setProjects(response.data.items.projects);
    }
    useEffect(() => {
        fetchProjects()
    }, []);
    return (
        <section id="rules" className="mb-12">
            <header className="mb-6 flex justify-between items-center">
                <span><h1 className="text-2xl font-bold">Projects</h1>
                <p className="text-gray-600">
                    Create and edit projects
                </p>
                    </span>
                <div className="flex space-x-4">
                    <button className="hover:bg-[#2044B8] hover:text-[#DBE9FE] bg-[#DBE9FE] text-[#2044B8] px-4 py-2 rounded ">
                        Фильтр по проектам
                    </button>
                    <button className=" px-4 py-2 rounded hover:bg-[#2044B8] hover:text-[#DBE9FE] bg-[#DBE9FE] text-[#2044B8]" onClick={() => navigate('/projects/new')}>
                        + Создать проект
                    </button>
                </div>
            </header>
            <div className="space-y-4">
                {projects ? (projects.map((project, index) => (
                    <ProjectCard key={index} project={project} onDelete={handleDeleteProjects}/>
                ))) : <p>У вас пока нет проектов</p>}
            </div>
        </section>
    );
}

export default ProjectsPage;