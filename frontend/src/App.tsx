import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { Layout } from 'antd'
import MainLayout from './layouts/MainLayout'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import UserList from './pages/users/UserList'
import CardList from './pages/cards/CardList'
import CoachList from './pages/coaches/CoachList'
import CourseList from './pages/courses/CourseList'

const App: React.FC = () => {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/" element={<MainLayout />}>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="users" element={<UserList />} />
        <Route path="cards" element={<CardList />} />
        <Route path="coaches" element={<CoachList />} />
        <Route path="courses" element={<CourseList />} />
      </Route>
    </Routes>
  )
}

export default App
