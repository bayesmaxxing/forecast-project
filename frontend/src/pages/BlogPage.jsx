import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Container,
  Box,
  Button,
  Typography,
  Card,
  CardContent,
  CardActionArea,
  Grid
} from '@mui/material';
import { ArrowBack } from '@mui/icons-material';
import BlogPost from '../components/BlogPost';

// Blog posts data - add new posts here
const BLOG_POSTS = [
  {
    slug: 'update-agent-forecasting',
    title: 'Update on AI Agent Forecasting',
    excerpt: 'An update on the performance of the AI agents running on this website.',
    date: '2026-01-19'
  },
  {
    slug: 'ai-agent-forecasting',
    title: 'AI Agent Forecasting',
    excerpt: 'An exploration of using AI agents for forecasting.',
    date: '2025-09-28'
  },
];

const BlogPage = () => {
  const { slug } = useParams();
  const navigate = useNavigate();

  // If viewing a specific post
  if (slug) {
    return (
      <Container maxWidth="lg" sx={{ mt: { xs: 6, sm: 8 } }}>
        <Box sx={{ py: 4 }}>
          <Button
            startIcon={<ArrowBack />}
            onClick={() => navigate('/blog')}
            sx={{ mb: 3 }}
          >
            Back to Blog
          </Button>
          <BlogPost slug={slug} />
        </Box>
      </Container>
    );
  }

  // Blog list view
  return (
    <Container maxWidth="lg" sx={{ mt: { xs: 6, sm: 8 } }}>
      <Box sx={{ py: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom sx={{ mb: 4 }}>
          Blog
        </Typography>

        <Grid container spacing={3}>
          {BLOG_POSTS.map((post) => (
            <Grid item xs={12} md={6} key={post.slug}>
              <Card
                elevation={2}
                sx={{
                  height: '100%',
                  display: 'flex',
                  flexDirection: 'column',
                  '&:hover': {
                    transform: 'translateY(-2px)',
                    transition: 'all 0.2s ease-in-out'
                  }
                }}
              >
                <CardActionArea
                  onClick={() => navigate(`/blog/${post.slug}`)}
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
    </Container>
  );
};

export default BlogPage;
