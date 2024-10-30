import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import Markdown from 'marked-react';
import './BlogPage.css';

function BlogPosts() {
    const [blogposts, setBlogposts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    
    useEffect(() => {
      fetch('https://forecasting-389105.ey.r.appspot.com/blogposts', {
        headers : {
          "Accept": "application/json"
        }
      })
      .then(response => {
        if (!response.ok) {
          throw new Error('Error fetching data');
        }
        return response.json();
      })
      .then(data => {
        if (data.length > 0) {
          setBlogposts(data);
        } else {
          throw new Error('No blogposts found');
        }
        setLoading(false);
      })
      .catch(error => {
        setError(error);
        setLoading(false);
      });
    }, []);

    const sortedBlogposts = [...blogposts].sort((a, b)=>{
        return b.post_id - a.post_id;
      });
  
    const truncatePostContent = (content, wordLimit) => {
      const wordsArray = content.split(' ');
      if (wordsArray.length > wordLimit) {
        return wordsArray.slice(0, wordLimit).join(' ') + '...';
      }
      return content;
    };

    return (
        <div>
          <ul className="blogpost-list">
            {sortedBlogposts.map(blogposts => (
              <li key={blogposts.slug} className="blog-item">
                <div className="blog-container">
                <Link to={`/blog/${blogposts.slug}`} className="blog-link">
                    <Markdown>{blogposts.title}</Markdown>
                  </Link>
                  <Markdown>{truncatePostContent(blogposts.post, 80)}</Markdown>
                </div>
              </li>
            ))}
          </ul>
        </div>
      );
    };
    
export default BlogPosts;
