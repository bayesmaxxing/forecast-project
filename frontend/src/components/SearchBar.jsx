import React, { useState } from 'react';
import { TextField, InputAdornment } from '@mui/material';
import { Search as SearchIcon } from '@mui/icons-material';

function SearchBar({ onSearch, placeholder = "Search..." }) {
  const [searchQuery, setSearchQuery] = useState('');

  const handleSearchChange = (e) => {
    const query = e.target.value.toLowerCase();
    setSearchQuery(query);
    onSearch(query);
  };

  return (
    <TextField
      fullWidth
      variant="outlined"
      placeholder={placeholder}
      value={searchQuery}
      onChange={handleSearchChange}
      sx={{
        mb: 2,
        '& .MuiOutlinedInput-root': {
          backgroundColor: 'background.paper',
          '& fieldset': {
            borderColor: 'primary.main',
          },
          '&:hover fieldset': {
            borderColor: 'primary.light',
          },
        },
        '& input': {
          color: 'primary.light',
        }
      }}
      slots={{
        startAdornment: () => (
          <InputAdornment position="start">
            <SearchIcon sx={{ color: 'primary.main' }} />
          </InputAdornment>
        ),
      }}
    />
  );
}

export default SearchBar;