import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import Markdown from 'marked-react';
import './BlogpostPage.css';


function BlogpostPage() {
    const [blogpost, setBlogpost] = useState({});
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    let { slug } = useParams();
  
    useEffect(() => {
        fetch(`https://forecast-project-backend.vercel.app/forecaster/api/blogposts/?slug=${slug}`, {
          headers : {
            'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
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
            setBlogpost(data[0]); 
          } else {
            throw new Error('No blog posts found');
          }
          setLoading(false);
        })
        .catch(error => {
          setError(error);
          setLoading(false);
        });
    }, [slug]);

    
    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error loading the forecast: {error.message}</div>;
    return (
      <div>
              <article>
                  <header>
                      <h1 className="blogpost-title">ASSASAS {blogpost.title}</h1>
                  </header>
                  <section className="blogpost-content">
                      <Markdown>{blogpost.post}</Markdown>
                  </section>
              </article>
      </div>
  );
  
}

export default BlogpostPage;