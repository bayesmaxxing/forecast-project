import React, { useState, useEffect } from 'react';
import { 
  FormControl, 
  InputLabel, 
  Select, 
  MenuItem, 
  Box, 
  CircularProgress 
} from '@mui/material';
import { useUserData } from '../services/hooks/useUserData';

const UserSelector = ({ onUserChange, selectedUserId = 'all' }) => {
  const { users = [], userLoading, error } = useUserData() || {};

  const handleChange = (event) => {
    const userId = event.target.value;
    onUserChange(userId);
  };

  return (
    <Box sx={{ minWidth: 220 }}>
      <FormControl fullWidth>
        <InputLabel id="user-selector-label">Select User</InputLabel>
        <Select
          labelId="user-selector-label"
          id="user-selector"
          value={selectedUserId}
          label="Select User"
          onChange={handleChange}
          disabled={userLoading}
        >
          <MenuItem value='all'>All Users</MenuItem>
          {userLoading ? (
            <MenuItem disabled>
              <CircularProgress size={20} sx={{ mr: 1 }} />
              Loading users...
            </MenuItem>
          ) : (
            users.map((user) => (
              <MenuItem key={user.id} value={user.id}>
                {user.username}
              </MenuItem>
            ))
          )}
        </Select>
      </FormControl>
      {error && <Box sx={{ color: 'error.main', mt: 1, fontSize: 'small' }}>{error}</Box>}
    </Box>
  );
};

export default UserSelector;
