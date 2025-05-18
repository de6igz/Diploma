import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Header from './components/ui/Header';
import Sidebar from './components/ui/Sidebar';
import Dashboard from './components/pages/dashboard/Dashboard';
import RulesPage from "./components/pages/rules/Rules.jsx";
import RuleForm from "./components/pages/rules/ruleForm/RuleForm.jsx";
import ProjectsPage from "./components/pages/projects/Projects.jsx";
import ProjectForm from "./components/pages/projects/ProjectForm.jsx";
import ProjectPage from "./components/pages/projects/ProjectPage.jsx";
import EventsByType from "./components/pages/event/EventsByType.jsx";
import EventsPage from "./components/pages/event/EventsPage.jsx";
import EventDetailsPage from "./components/pages/event/EventDetails.jsx";
import AuthPage from "./components/pages/auth/Auth.jsx";
import { SnackbarProvider } from "notistack";
import Settings from "./components/pages/settings/Settings.jsx";

// Компоненты для защиты маршрутов
const ProtectedRoute = ({ children, isAuthenticated }) => {
    return isAuthenticated ? children : <Navigate to="/" />;
};

const AuthRoute = ({ children, isAuthenticated }) => {
    return isAuthenticated ? <Navigate to="/dashboard" /> : children;
};

function App() {
    const [isSidebarOpen, setIsSidebarOpen] = useState(false);
    const [isAuthenticated, setIsAuthenticated] = useState(!!localStorage.getItem('accessToken'));

    return (
        <SnackbarProvider maxSnack={3}>
            <Router>
                <div className="h-screen flex flex-col overflow-hidden">
                    {isAuthenticated && <Header onToggleSidebar={() => setIsSidebarOpen(!isSidebarOpen)} />}
                    <div className="flex flex-1 overflow-hidden">
                        {isAuthenticated && <Sidebar isOpen={isSidebarOpen} />}
                        <main className={isAuthenticated ? "flex-1 overflow-y-auto p-6 md:p-8 bg-gray-50" : "flex-1 overflow-y-auto bg-gray-50"}>
                            <Routes>
                                <Route
                                    path="/"
                                    element={
                                        <AuthRoute isAuthenticated={isAuthenticated}>
                                            <AuthPage setIsAuthenticated={setIsAuthenticated} />
                                        </AuthRoute>
                                    }
                                />
                                <Route
                                    path="/dashboard"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <Dashboard />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route
                                    path="/events"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <EventsPage />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route
                                    path="/events/details/:id"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <EventDetailsPage />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route
                                    path="/events/:eventType"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <EventsByType />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route
                                    path="/projects"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <ProjectsPage />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route
                                    path="/projects/:id"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <ProjectPage />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route
                                    path="/projects/new"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <ProjectForm />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route
                                    path="/projects/:id/edit"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <ProjectForm />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route
                                    path="/rules"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <RulesPage />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route
                                    path="/rules/new"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <RuleForm />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route
                                    path="/rules/:ruleType/:id/edit"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <RuleForm />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route
                                    path="/settings"
                                    element={
                                        <ProtectedRoute isAuthenticated={isAuthenticated}>
                                            <Settings />
                                        </ProtectedRoute>
                                    }
                                />
                                <Route path="*" element={<h1 className="text-2xl">404 - Not Found</h1>} />
                            </Routes>
                        </main>
                    </div>
                </div>
            </Router>
        </SnackbarProvider>
    );
}

export default App;