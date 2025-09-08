import React, { useState } from 'react';
import {
  Paper,
  Box,
  Button,
  Typography,
  Collapse,
  IconButton,
  CircularProgress,
  Alert,
  Chip
} from '@mui/material';
import {
  NewspaperOutlined as NewsIcon,
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon,
  Refresh as RefreshIcon
} from '@mui/icons-material';
import MarkdownRenderer from 'marked-react';
import { newsService } from '../services/api/index';

function ForecastNews({ forecastQuestion, forecastId }) {
  const [loading, setLoading] = useState(false);
  const [newsContent, setNewsContent] = useState(null);
  const [error, setError] = useState(null);
  const [showNews, setShowNews] = useState(false);
  const [fetchTimestamp, setFetchTimestamp] = useState(null);

  const handleFetchNews = async () => {
    if (!forecastQuestion) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const response = await newsService.fetchForecastNews(forecastQuestion);
      setNewsContent(response.news || response);
      setFetchTimestamp(new Date());
      setShowNews(true);
    } catch (err) {
      setError(err.message || 'Failed to fetch news');
      console.error('Error fetching forecast news:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleRetry = () => {
    setError(null);
    handleFetchNews();
  };

  const toggleNewsDisplay = () => {
    setShowNews(!showNews);
  };

  const formatTimestamp = (date) => {
    return date.toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <Paper elevation={3} sx={{ p: 3, mb: 3 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 2 }}>
        <Typography variant="h6" sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <NewsIcon />
          Latest News & Context
        </Typography>
        
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          {fetchTimestamp && (
            <Chip 
              label={`Updated: ${formatTimestamp(fetchTimestamp)}`}
              size="small"
              variant="outlined"
            />
          )}
          
          {newsContent && (
            <IconButton onClick={toggleNewsDisplay} size="small">
              {showNews ? <ExpandLessIcon /> : <ExpandMoreIcon />}
            </IconButton>
          )}
        </Box>
      </Box>

      {/* Fetch Button */}
      {!newsContent && !loading && (
        <Box sx={{ textAlign: 'center', py: 2 }}>
          <Button
            variant="contained"
            onClick={handleFetchNews}
            startIcon={<NewsIcon />}
            disabled={loading || !forecastQuestion}
            size="large"
          >
            Get Latest News
          </Button>
          {!forecastQuestion && (
            <Typography variant="caption" display="block" sx={{ mt: 1, color: 'text.secondary' }}>
              Forecast question required to fetch news
            </Typography>
          )}
        </Box>
      )}

      {/* Loading State */}
      {loading && (
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'center', py: 4 }}>
          <CircularProgress size={30} sx={{ mr: 2 }} />
          <Typography variant="body1" color="text.secondary">
            Fetching latest news and analysis...
          </Typography>
        </Box>
      )}

      {/* Error State */}
      {error && (
        <Alert 
          severity="error" 
          action={
            <Button 
              color="inherit" 
              size="small" 
              onClick={handleRetry}
              startIcon={<RefreshIcon />}
            >
              Retry
            </Button>
          }
          sx={{ mb: 2 }}
        >
          {error}
        </Alert>
      )}

      {/* News Content */}
      {newsContent && (
        <>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
            <Button
              variant="outlined"
              onClick={handleFetchNews}
              startIcon={<RefreshIcon />}
              disabled={loading}
              size="small"
            >
              Refresh News
            </Button>
          </Box>

          <Collapse in={showNews}>
            <Box
              sx={{
                bgcolor: 'background.default',
                borderRadius: 1,
                p: 3,
                border: 1,
                borderColor: 'divider',
                maxHeight: '70vh',
                overflow: 'auto',
                '& p': {
                  mb: 2,
                  lineHeight: 1.7,
                },
                '& p:last-child': {
                  mb: 0,
                }
              }}
            >
              <Box
                component="div"
                sx={{ 
                  wordBreak: 'break-word',
                  '& h1, & h2, & h3, & h4, & h5, & h6': {
                    mt: 2,
                    mb: 1,
                    fontWeight: 600,
                  },
                  '& h1': { fontSize: '1.5rem' },
                  '& h2': { fontSize: '1.3rem' },
                  '& h3': { fontSize: '1.1rem' },
                  '& p': {
                    mb: 1,
                    lineHeight: 1.7,
                  },
                  '& ul, & ol': {
                    pl: 2,
                    mb: 1,
                  },
                  '& li': {
                    mb: 0.5,
                  },
                  '& blockquote': {
                    borderLeft: '4px solid',
                    borderColor: 'primary.main',
                    pl: 2,
                    ml: 0,
                    fontStyle: 'italic',
                    opacity: 0.9,
                  },
                  '& code': {
                    bgcolor: 'action.hover',
                    px: 0.5,
                    py: 0.25,
                    borderRadius: 0.5,
                    fontSize: '0.875rem',
                  },
                  '& pre': {
                    bgcolor: 'action.hover',
                    p: 1.5,
                    borderRadius: 1,
                    overflow: 'auto',
                  }
                }}
              >
                {typeof newsContent === 'string' 
                  ? <MarkdownRenderer>{newsContent}</MarkdownRenderer>
                  : <Typography component="pre">{JSON.stringify(newsContent, null, 2)}</Typography>
                }
              </Box>
            </Box>
          </Collapse>
        </>
      )}
    </Paper>
  );
}

export default ForecastNews;