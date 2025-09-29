import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Card,
  CardContent,
  CardActionArea,
  Grid,
  CircularProgress,
  Alert
} from '@mui/material';
import { useNavigate } from 'react-router-dom';

const BlogList = () => {
  const [blogPosts, setBlogPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchBlogList = async () => {
      try {
        setLoading(true);
        // In a real app, this would come from an API or be statically generated
        // For now, we'll use a hardcoded list of blog posts
        const posts = [
          {
            slug: 'ai-agent-forecasting',
            title: 'AI Agent Forecasting',
            excerpt: 'An exploration of using AI agents for forecasting.',
            date: '2025-09-28'
          },
        ];

        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 500));
        setBlogPosts(posts);
      } catch (err) {
        setError('Failed to load blog posts');
      } finally {
        setLoading(false);
      }
    };

    fetchBlogList();
  }, []);

  const handlePostClick = (slug) => {
    navigate(`/blog/${slug}`);
  };

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
    <Box sx={{ py: 4 }}>
      <Typography variant="h3" component="h1" gutterBottom>
        
      </Typography>
      <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
        
      </Typography>

      <Grid container spacing={3}>
        {blogPosts.map((post) => (
          <Grid item xs={12} md={6} key={post.slug}>
            <Card
              elevation={2}
              sx={{
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                '&:hover': {
                  elevation: 4,
                  transform: 'translateY(-2px)',
                  transition: 'all 0.2s ease-in-out'
                }
              }}
            >
              <CardActionArea
                onClick={() => handlePostClick(post.slug)}
                sx={{ height: '100%', display: 'flex', flexDirection: 'column', alignItems: 'stretch' }}
              >
                <CardContent sx={{ flexGrow: 1 }}>
                  <Typography variant="h5" component="h2" gutterBottom>
                    {post.title}
                  </Typography>
                  <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                    {new Date(post.date).toLocaleDateString()}
                  </Typography>
                  <Typography variant="body1" color="text.secondary">
                    {post.excerpt}
                  </Typography>
                </CardContent>
              </CardActionArea>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default BlogList;