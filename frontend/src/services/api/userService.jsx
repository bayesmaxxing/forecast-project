import { API_BASE_URL } from './index';
export const fetchUsers = async () => {
  const response = await fetch(`${API_BASE_URL}/users`, {
    headers: { "Accept": "application/json" }
  });
  
  if (!response.ok) {
    throw new Error(`Error fetching users: ${response.status}`);
  }

  return response.json();
};

export const createUser = async (username, password) => {
  const response = await fetch(`${API_BASE_URL}/users`, {
    method: 'POST',
    headers: { "Accept": "application/json", "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  });

  if (!response.ok) {
    throw new Error(`Error creating user: ${response.status}`);
  }

  return response.json();
};

export const loginUser = async (username, password) => {
  const response = await fetch(`${API_BASE_URL}/users/login`, {
    method: 'POST',
    headers: { "Accept": "application/json", "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  });

  if (!response.ok) {
    throw new Error(`Error logging in: ${response.status}`);
  }

  return response.json();
};