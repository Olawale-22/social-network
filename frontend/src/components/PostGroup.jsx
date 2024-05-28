
import React from 'react';
import { Link } from "react-router-dom"


const PostGroup = ({ group }) => {
    return (

        <div className="post-card" key={group.ID}>
        <div className="post-header">
          <Link to={`/profile/${group.ID}`}>
            <div className="post-info">
              <h3 className="username">{group.Name}</h3>
              <p className="timestamp">{group.Mentions.length} <i> members</i></p>
              {/* <p className="timestamp">About Me: {item.AboutMe}</p> */}
            </div>
          </Link>
        </div>
      </div>
    );
};

export default PostGroup;
