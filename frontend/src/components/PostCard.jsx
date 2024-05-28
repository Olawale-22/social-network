import '../styles/post.css';

import React, { useContext, useState } from 'react';
import { formatDate } from '../utils/utils';
import { useSelector } from 'react-redux';
import WebSocketContext from '../contexts/WebSocketContext';

const PostCard = ({ post, user }) => {
    const comments = useSelector(state => state.data.comments);
    const { socket } = useContext(WebSocketContext);
    const postComments = comments.filter(comment => comment.Post_Id === post.PiD)
    const users = useSelector(state => state.data.users);
    const [commentsForm, setCommentsForm] = useState(false);


    const renderComments = (comments, allUsers) => {
        return comments.map((comment) => {
            if (comment && comment.User_Id !== undefined) {
                const user = allUsers.find((item) => item.ID === post.User_Id);
                if (user) {
                    return (
                        <div className="comment" key={comment.ID}>
                            <div className="comment-header">
                                <img
                                    src={user.Avatar}
                                    alt={`${post.User_Id}'s Profile`}
                                    className="profile-image"
                                />
                                <h3 className="username">{user.Nickname}</h3>
                                <span className="timestamp">{formatDate(post.Time)}</span>
                            </div>
                            <div className="comment-info">
                                <p className="comment-content">{comment.Content}</p>
                            </div>
                        </div>
                    )
                }
            }
            return null;
        });
    };

    const handleComments = () => {
        setCommentsForm(curr => !curr);
    }

    const handleSubmit = event => {
        event.preventDefault();
        const formData = event.target.elements;
        const data = {
            event: "comments",
            ID: localStorage.getItem("user_id"),
            postId: String(post.PiD),
            comment: formData.comment.value,
        }

        socket.send(JSON.stringify(data));
    }


    return (
        <div className="post-card" key={post.PiD}>
            <div className="post-header">
                <img
                    src={user.Avatar}
                    alt={`${post.User_Id}'s Profile`}
                    className="profile-image"
                />
                <h3 className="username">{user.Nickname}</h3>
                <span className="timestamp">{formatDate(post.Time)}</span>
            </div>
            <div className="post-info">
                <p className="post-content">{post.Content}</p>
                {post.Image && (
                    <img
                        src={post.Image}
                        alt={`${post.User_Id}' post`}
                        className="post-image"
                    />
                )}
            </div>

            <div className="toggle-comments">
                <span onClick={handleComments}>{postComments.length} Comments</span>
            </div>

            {commentsForm && (
                <div className="post_comments">
                    <form onSubmit={handleSubmit}>
                        <textarea name="comment" id="comment"></textarea>
                        <button type="submit">Comment</button>
                    </form>

                    {postComments.length > 0 && users && (
                        <div className="comments">
                            {renderComments(postComments, users)}
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};

export default PostCard;
