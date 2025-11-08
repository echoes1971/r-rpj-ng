import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import AppNavbar from "./Navbar";
import './App.css';

import Login from "./Login";
import Users from "./Users";

function App() {
  const token = localStorage.getItem("token");
  return (
    <Router>
      <AppNavbar />
      <Routes>
        {/* Rotta login sempre accessibile */}
        <Route path="/login" element={<Login />} />

        {/* Rotta utenti protetta */}
        <Route
          path="/users"
          element={token ? <Users /> : <Navigate to="/login" />}
        />

        {/* Rotta di default → redirect a /login */}
        <Route path="*" element={<Navigate to="/login" />} />
      </Routes>
    </Router>
  );
}

export default App;
