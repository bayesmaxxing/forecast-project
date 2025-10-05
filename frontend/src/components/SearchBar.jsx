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
          borderRadius: 2,
          transition: 'all 0.2s ease-in-out',
          '& fieldset': {
            borderColor: 'divider',
            borderWidth: '1px',
          },
          '&:hover': {
            backgroundColor: 'background.paper',
            '& fieldset': {
              borderColor: 'primary.main',
            },
          },
          '&.Mui-focused': {
            backgroundColor: 'background.paper',
            boxShadow: '0 0 0 3px rgba(59, 130, 246, 0.1)',
            '& fieldset': {
              borderColor: 'primary.main',
              borderWidth: '2px',
            },
          },
        },
        '& input': {
          color: 'text.primary',
          fontSize: '0.95rem',
          '&::placeholder': {
            color: 'text.secondary',
            opacity: 0.6,
          },
        },
      }}
      slots={{
        startAdornment: () => (
          <InputAdornment position="start">
            <SearchIcon
              sx={{
                color: 'text.secondary',
                transition: 'color 0.2s ease-in-out',
                '.MuiOutlinedInput-root.Mui-focused &': {
                  color: 'primary.main',
                },
              }}
            />
          </InputAdornment>
        ),
      }}
    />
  );
}

export default SearchBar;