import React, { useState } from 'react';
import LoginForm from '../components/LoginForm';
import AddForecast from '../components/AddForecast'

const AdminPage = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  const handleLoginSuccess = () => {
    setIsLoggedIn(true);
  };

  if (!isLoggedIn) {
    return <LoginForm onLoginSuccess={handleLoginSuccess} />;
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Add forecast</h1>
      <AddForecast></AddForecast>
    </div>
  );
};

export default AdminPage;