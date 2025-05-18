import React, { useState, useEffect } from 'react';
import { Link, NavLink } from 'react-router-dom';
import aletheiaClient from "../pages/api/aletheiaClient.js";

function Sidebar({ isOpen }) {
    const [projects, setProjects] = useState([]);

    useEffect(() => {
        async function fetchProjects() {
            try {
                const response = await aletheiaClient.get('/projects');
                const projectsData = response.data.items.projects || [];
                setProjects(projectsData);
            } catch (error) {
                console.error('Error fetching projects:', error);
                setProjects([]);
            }
        }
        fetchProjects();
    }, []);

    return (
        <aside
            className={`bg-white shadow-md h-screen fixed md:sticky top-0 left-0 overflow-y-auto z-40 transition-transform transform ${
                isOpen ? 'translate-x-0' : '-translate-x-full'
            } md:translate-x-0 w-64`}
        >
            <div className="py-6 px-4">
                <nav className="space-y-1">
                    <NavLink
                        to="/dashboard"
                        className={({ isActive }) =>
                            `flex items-center space-x-3 px-3 py-2 rounded-lg ${
                                isActive ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                            }`
                        }
                        end
                    >
                        <i className="fas fa-tachometer-alt"></i>
                        <span>Dashboard</span>
                    </NavLink>

                    <NavLink
                        to="/events"
                        className={({ isActive }) =>
                            `flex items-center space-x-3 px-3 py-2 rounded-lg ${
                                isActive ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                            }`
                        }
                    >
                        <i className="fas fa-calendar-check"></i>
                        <span>Events</span>
                    </NavLink>
                    <NavLink
                        to="/projects"
                        className={({ isActive }) =>
                            `flex items-center space-x-3 px-3 py-2 rounded-lg ${
                                isActive ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                            }`
                        }
                    >
                        <i className="fa-solid fa-list-check"></i>
                        <span>Projects</span>
                    </NavLink>
                    <NavLink
                        to="/rules"
                        className={({ isActive }) =>
                            `flex items-center space-x-3 px-3 py-2 rounded-lg ${
                                isActive ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                            }`
                        }
                    >
                        <i className="fa-solid fa-scale-balanced"></i>
                        <span>Rules</span>
                    </NavLink>

                    <NavLink
                        to="/settings"
                        className={({ isActive }) =>
                            `flex items-center space-x-3 px-3 py-2 rounded-lg ${
                                isActive ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                            }`
                        }
                    >
                        <i className="fas fa-cog"></i>
                        <span>Settings</span>
                    </NavLink>
                </nav>

                {/* Раздел "Recent projects" */}
                <div className="mt-10 pt-6 border-t border-gray-200">
                    <h3 className="text-xs uppercase text-gray-500 font-semibold px-3 mb-2">
                        Recent projects
                    </h3>
                    <nav className="space-y-1">
                        {projects.length === 0 ? (
                            <p className="px-3 py-2 text-gray-500">No recent projects</p>
                        ) : (
                            projects.slice(0, 3).reverse().map((project) => (
                                <Link
                                    key={project.id}
                                    to={`/projects/${project.id}`}
                                    className="flex items-center justify-between px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-100"
                                >
                                    <div className="flex items-center space-x-3 truncate">
                                        <div className="h-2 w-2 rounded-full bg-green-500"></div>
                                        <span className="truncate">{project.projectName}</span>
                                    </div>
                                    <span className="bg-green-100 text-green-800 text-xs px-2 py-1 rounded">
                                        Healthy
                                    </span>
                                </Link>
                            ))
                        )}
                    </nav>
                </div>
            </div>
        </aside>
    );
}

export default Sidebar;