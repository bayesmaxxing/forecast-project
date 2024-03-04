import React, { useState, useEffect } from 'react';
import { marked } from 'marked';

function BlogPosts() {
    const [posts, setPosts] = useState([]);

    useEffect(() => {
        // List of your markdown files
        const postFiles = [
            'posts/grasped_sword.md',
            'posts/policy_dollars.md',
            'posts/social_status_framework.md',
            'posts/fermi_grabby_aliens.md'
            // Add more posts as needed
        ];

        const fetchPosts = async () => {
            const postsData = await Promise.all(
                postFiles.map(async (file) => {
                    const response = await fetch(file);
                    const text = await response.text();
                    return marked(text);
                })
            );
            setPosts(postsData);
        };

        fetchPosts();
    }, []);

    return (
        <div>
            {posts.map((postHtml, index) => (
                <div key={index} dangerouslySetInnerHTML={{ __html: postHtml }} />
            ))}
        </div>
    );
}

export default BlogPosts;
