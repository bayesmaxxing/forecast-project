import React from 'react';
import { useParams } from 'react-router-dom';
import { Container, Box, Button } from '@mui/material';
import { ArrowBack } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import BlogPost from '../components/BlogPost';
import BlogList from '../components/BlogList';

const BlogPage = () => {
  const { slug } = useParams();
  const navigate = useNavigate();

  const handleBackToBlog = () => {
    navigate('/blog');
  };

  return (
    <Container maxWidth="lg">
      <Box sx={{ py: 4 }}>
        {slug ? (
          <>
            <Button
              startIcon={<ArrowBack />}
              onClick={handleBackToBlog}
              sx={{ mb: 3 }}
            >
              Back to Blog
            </Button>
            <BlogPost slug={slug} />
          </>
        ) : (
          <BlogList />
        )}
      </Box>
    </Container>
  );
};

export default BlogPage;