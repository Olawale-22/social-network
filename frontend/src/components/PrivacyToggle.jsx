import React, { useContext, useEffect, useState } from 'react'
import { useSelector, useDispatch } from 'react-redux';
import WebSocketContext from '../contexts/WebSocketContext'
import { updateOwnPrivacy } from '../redux/userSlice';
import { updatePrivacy } from '../redux/dataSlice';

const PrivacyToggle = ({ setMessage }) => {
    const { socket, websocketData } = useContext(WebSocketContext);
    const userInfos = useSelector(state => state.user.userInfo);
    const dispatch = useDispatch();
    const userId = localStorage.getItem("user_id");


    useEffect(() => {
        if (websocketData) {
            if (websocketData.Event === 'user_update_profil') {
                setMessage("Profile's privacy updated !")
                dispatch(updateOwnPrivacy());
                dispatch(updatePrivacy(parseInt(userId)));
                setTimeout(() => {
                    setMessage("")
                }, 2000);
            }
        }
    }, [websocketData]);

    const handleToggle = () => {
        socket.send(JSON.stringify({
            event: "profile",
            id: userId
        }))
    }
    return (
        <>
            {userInfos.IsPublic !== null &&
                <>
                    <label>
                        Public
                        <input
                            type="checkbox"
                            checked={userInfos.IsPublic}
                            onChange={handleToggle}
                        />
                    </label>
                    <label>
                        Private
                        <input
                            type="checkbox"
                            checked={!userInfos.IsPublic}
                            onChange={handleToggle}
                        />
                    </label>
                </>
            }
        </>
    )
}
export default PrivacyToggle