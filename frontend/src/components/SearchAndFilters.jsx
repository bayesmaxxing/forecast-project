import React from 'react';
import { Grid2 } from '@mui/material';
import SearchBar from './SearchBar';
import UserSelector from './UserSelector';

function SearchAndFilters({ onSearch, selectedUserId, onUserChange }) {
  return (
    <Grid2 container spacing={2} alignItems="center">
      <Grid2 xs={12} md={8}>
        <SearchBar 
          onSearch={onSearch}
          placeholder="Search questions..."
        />
      </Grid2>
      <Grid2 xs={12} md={4}>
        <UserSelector
          selectedUserId={selectedUserId}
          onUserChange={onUserChange}
        />
      </Grid2>
    </Grid2>
  );
}

export default SearchAndFilters;