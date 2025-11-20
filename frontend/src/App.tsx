import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import ProtectedRoute from './components/ProtectedRoute'
import MainLayout from './layouts/MainLayout'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import UserList from './pages/users/UserList'
import UserDetail from './pages/users/UserDetail'
import CardList from './pages/cards/CardList'
import CardDetail from './pages/cards/CardDetail'
import CoachList from './pages/coaches/CoachList'
import CoachDetail from './pages/coaches/CoachDetail'
import CourseList from './pages/courses/CourseList'
import CourseBooking from './pages/courses/CourseBooking'

const App: React.FC = () => {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/" element={
        <ProtectedRoute>
          <MainLayout />
        </ProtectedRoute>
      }>
        <Route index element={<Navigate to="/dashboard" replace />} />
        <Route path="dashboard" element={<Dashboard />} />
        <Route path="users" element={<UserList />} />
        <Route path="users/:id" element={<UserDetail />} />
        <Route path="cards" element={<CardList />} />
        <Route path="cards/:id" element={<CardDetail />} />
        <Route path="coaches" element={<CoachList />} />
        <Route path="coaches/:id" element={<CoachDetail />} />
        <Route path="courses" element={<CourseList />} />
        <Route path="courses/booking" element={<CourseBooking />} />
      </Route>
    </Routes>
  )
}

export default App