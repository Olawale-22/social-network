import React, { useEffect, useContext } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import WebSocketContext from '../contexts/WebSocketContext';
import { useLocation } from 'react-router-dom';
import PostCard from '../components/PostCard';
import { addCommentToPost } from '../redux/dataSlice';

const Post = () => {
    const comments = useSelector(state => state.data.comments);
    const { socket, websocketData } = useContext(WebSocketContext);
    const dispatch = useDispatch();
    const location = useLocation();
    const nickname = location.state.nickname;
    const post = location.state.post;
    const postComments = comments.filter(comment => comment.Post_Id === post.PiD)

    useEffect(() => {
        if (websocketData) {
            if (websocketData.Event === 'CommentSuccessfull') {
                dispatch(addCommentToPost(websocketData.Data));
            }
        }
    }, [websocketData]);

    const handleSubmit = event => {
        event.preventDefault();
        const formData = event.target.elements;
        const data = {
            event: "comments",
            ID: localStorage.getItem("user_id"),
            postId: String(post.PiD),
            comment: formData.comment.value,
            image: ""
        }

        socket.send(JSON.stringify(data));
    }

    return (
        <>
            <div>
                <PostCard post={post} nickname={nickname} />
            </div>

            <div>
                <form onSubmit={handleSubmit}>
                    <textarea name="comment" id="comment" cols="15" rows="10"></textarea>
                    <button type="submit">Comment</button>
                </form>
            </div>


            {postComments && (
                postComments.map(comment => (
                    <p>{comment.Content}</p>
                ))
            )}
        </>
    )
}

export default Post