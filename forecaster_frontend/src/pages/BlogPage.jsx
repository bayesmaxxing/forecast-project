import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import Markdown from 'marked-react';
import './BlogPage.css';

function BlogPosts() {
    const [blogposts, setBlogposts] = useState([]);
    
    useEffect(() => {
        const CACHE_DURATION = 5 * 60 * 1000; // Cache duration in milliseconds, e.g., 5 minutes
        const now = new Date().getTime(); // Current time
    
        const blogpostsCache = localStorage.getItem('blogposts');
    
        const blogpostsDataValid = blogpostsCache && now - JSON.parse(blogpostsCache).timestamp < CACHE_DURATION;
        // Try to load data from cache
        if (blogpostsDataValid) {
            setBlogposts(JSON.parse(blogpostsCache).data);
        } else {
        // Fetch the list of resolutions from the API if cache is empty
        fetch('https://forecast-project-backend.vercel.app/forecaster/api/blogposts/', {
            headers : {
                'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
            }
            })
            .then(response => response.json())
            .then(data => {
            // Update state with fetched data
            setBlogposts(data);
            // Update cache with new data
            localStorage.setItem('blogposts', JSON.stringify({data: data, timestamp: now}));
            })
            .catch(error => console.error('Error fetching data: ', error));
        }
    }, []);

    const sortedBlogposts = [...blogposts].sort((a, b)=>{
        return b.post_id - a.post_id;
      });
  
      const formatDate = (dateString) => dateString.split('T')[0];

    return (
        <div>
          <ul className="blogpost-list">
            {sortedBlogposts.map(blogposts => (
              <li key={blogposts.slug} className="blog-item">
                <div className="blog-container">
                <Link to={`/blog/${blogposts.slug}`} className="blog-link">
                    <Markdown>{blogposts.title}</Markdown>
                  </Link>
                </div>
                <div>
                <p>{blogposts.summary}</p>
                </div>
              </li>
            ))}
          </ul>
        </div>
      );
    };
    
export default BlogPosts;
