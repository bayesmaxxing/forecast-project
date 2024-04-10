import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Marked } from 'marked';
import './BlogpostPage.css';


function BlogpostPage() {
    const [blogpost, setBlogpost] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    let { slug } = useParams();
  
    useEffect(() => {
        fetch(`https://forecast-project-backend.vercel.app/forecaster/api/blogposts/?slug=${slug}`, {
          headers : {
            'Authorization': `Token ${process.env.REACT_APP_API_TOKEN}`
          }
        })
      .then( blogpostData => {
        if (!blogpostData.ok) {
          throw new Error('Error fetching data');
        }
        const blogpostJson =blogpostData.json()
        return blogpostJson;
      })
      .then(blogpostJson => {
        setBlogpost(blogpostJson);
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
          {blogpost ? ( // Check if the blogpost object is not null
              <article>
                  <header>
                      <h1 className="blogpost-title">{blogpost.title}</h1>
                      <div className="blogpost-meta">
                      </div>
                  </header>
                  <section className="blogpost-content">
                      <p>{blogpost.post}</p> {/* Assuming 'content' is a field in your blogpost object */}
                  </section>
              </article>
          ) : (
              <div>Blog post not found.</div>
          )}
      </div>
  );
  
}

export default BlogpostPage;