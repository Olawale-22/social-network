import React, { useState } from 'react';
import { useSelector } from 'react-redux';
import { Link, } from 'react-router-dom';

import '../styles/home.css';
import PostForm from '../components/PostForm';
import PostCard from '../components/PostCard';
import Users from './Users';

const Home = () => {
  const userId = localStorage.getItem("user_id");
  const users = useSelector(state => state.data.users);
  const posts = useSelector(state => state.user.posts);
  const [privacy, setPrivacy] = useState("public");
  const [postForm, setPostForm] = useState(false);


  const renderPosts = (posts, allUsers) => {
    return posts.map((post) => {
      if (post && post.User_Id !== undefined) {
        const user = allUsers.find((item) => item.ID === post.User_Id);
        if (user) {
          return (
            <PostCard post={post} user={user} key={post.PiD} />
          )
        }
      }
      return null;
    });
  };


  return (
    <>
      <div className="users">
        <Users />
      </div>
      
      <div className="container-home">

        <div className="home-posts">
          <div className="btn-create-post">
            <button onClick={() => setPostForm(curr => !curr)}>POST</button>
          </div>

          {postForm && (
            <div className="home-post-creation">
              <PostForm userId={userId} privacy={privacy} setPrivacy={setPrivacy} />
            </div>
          )}

          {posts ? (
            <div className="posts">
              <h2>PUBLIC</h2>
              <hr />
              {posts && posts.public && renderPosts(posts.public, users)}
              <h2>PRIVATE</h2>
              <hr />
              {posts && posts.private && renderPosts(posts.private, users)}
              <h2>ALMOST PRIVATE</h2>
              <hr />
              {posts && posts.mentions && renderPosts(posts.mentions, users)}
            </div>
          ) : <p></p>}
        </div>
      </div>
    </>
  );
};

export default Home;