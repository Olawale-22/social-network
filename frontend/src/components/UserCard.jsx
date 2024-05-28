import React, { useContext, useState } from 'react'
import { useSelector, useDispatch } from 'react-redux';
import WebSocketContext from '../contexts/WebSocketContext';
import { Link, useNavigate } from 'react-router-dom';
import { updateRequestedUser } from '../redux/dataSlice';

const UserCard = ({ user, userId }) => {
    const dispatch = useDispatch();
    const users = useSelector(state => state.data.users);
    const follows = useSelector(state => state.user.follows);
    const followers = useSelector(state => state.user.followers);
    const { socket } = useContext(WebSocketContext);
    const [error, setError] = useState("");
    const navigate = useNavigate();


    const checkFollowing = otherId => {
        return follows.some(item => item === otherId)
    }

    const checkPublic = id => {
        return users.some(user => user.ID === id && user.IsPublic)
    }

    const handleFollowToggle = (event, userIdToFollow) => {
        event.preventDefault();
        const data = {
            event: "",
            ID: userId,
            FollowId: 0
        };

        if (checkFollowing(userIdToFollow)) {
            data.event = "unfollow_user";
            data.FollowId = String(userIdToFollow);
        } else {
            if (checkPublic(userIdToFollow)) {
                data.event = "follow_user";
                data.FollowId = String(userIdToFollow);
            } else {
                data.event = "user_follow_request_user";
                data.FollowId = JSON.stringify([userIdToFollow])
                dispatch(updateRequestedUser(userIdToFollow))
            }
        }
        socket.send(JSON.stringify(data));
    };
    

    const isUserAllowedToChat = id => {
        return checkFollowing(id) || followers.some(item => item === id);
    }

    const handleChat = id => {
        if (isUserAllowedToChat(id)) {
            navigate(`/chat/${id}`);
        } else {
            setError("none of you are following you");
            setTimeout(() => {
                setError("")
            }, 1500)
        }
    }


    return (
        <div className="user-card" key={user.ID}>
            <div className="user-card-image">
                <img
                    src={user.Avatar}
                    alt={`${user.ID}'s Profile`}
                    className="profile-image"
                />

                <div className="user-card-username">
                    <Link to={`/profile/${user.ID}`}>
                        <h3 className="username">{user.Nickname}</h3>
                    </Link>
                </div>
            </div>

            <div className="buttons">
                <button
                    onClick={event => handleFollowToggle(event, user.ID)}
                >
                    {user.IsRequested ? 'Waiting...' : checkFollowing(user.ID) ? 'Unfollow' : checkPublic(user.ID) ? 'Follow' : 'Request'}
                </button>
                <button onClick={() => handleChat(user.ID)} className="chat">Chat</button>
                {error && <p className="error">{error}</p>}
            </div>
        </div >
    )
}

export default UserCard