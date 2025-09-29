import React, { useState, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import { Box, CircularProgress, Alert, Paper } from '@mui/material';

const BlogPost = ({ slug }) => {
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchMarkdownFile = async () => {
      try {
        setLoading(true);
        const response = await fetch(`/blog/${slug}.md`);
        if (!response.ok) {
          throw new Error('Blog post not found');
        }
        const text = await response.text();
        setContent(text);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    if (slug) {
      fetchMarkdownFile();
    }
  }, [slug]);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" my={4}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box my={4}>
        <Alert severity="error">{error}</Alert>
      </Box>
    );
  }

  return (
    <Paper
      elevation={1}
      sx={{
        p: 4,
        my: 2,
        '& h1, & h2, & h3, & h4, & h5, & h6': {
          color: 'primary.main',
          mb: 2
        },
        '& p': {
          mb: 2,
          lineHeight: 1.6
        },
        '& pre': {
          backgroundColor: 'grey.800',
          color: 'grey.100',
          p: 2,
          borderRadius: 1,
          overflow: 'auto'
        },
        '& blockquote': {
          borderLeft: 4,
          borderColor: 'primary.main',
          ml: 0,
          pl: 2,
          fontStyle: 'italic'
        }
      }}
    >
      <ReactMarkdown>{content}</ReactMarkdown>
    </Paper>
  );
};

export default BlogPost;