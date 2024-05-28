import React, { useContext } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { getUserById, updateRequestedGroup } from '../redux/dataSlice';
import { Link } from 'react-router-dom';
import WebSocketContext from '../contexts/WebSocketContext';

const GroupCard = ({ userId, group, isMember }) => {
    const dispatch = useDispatch();
    const { socket } = useContext(WebSocketContext);
    const user = useSelector(state => getUserById(state.data.users, group.Admin_id));

    const handleGroupRequest = (groupID, adminID) => {
        console.log(groupID);
        socket.send(JSON.stringify({
            event: "user_group_request_admin",
            ID: userId,
            Recipient: JSON.stringify([adminID]),
            GroupID: String(groupID)
        }))
        dispatch(updateRequestedGroup(groupID));
    }

    const commonContent = (
        <div className="group">
            <div className="group-header">
                <h3 className="name">{group.Name}</h3>
            </div>

            <div className="group-info">
                <span>Admin: {user.Nickname}</span>
                <p className="group-description">{group.Descriptions}</p>
                {group.Mentions !== null && <span className="group-count-members">Members: {1 + group.Mentions.length}</span>}
                {isMember ? null : <button
                    onClick={() => handleGroupRequest(group.ID, group.Admin_id)}>
                    {group.IsRequested ? "Waiting..." : "REQUEST"}
                </button>}
            </div>
        </div>
    );

    return isMember ? (
        <Link to={`/messenger/${group.ID}`}>
            {commonContent}
        </Link>
    ) : commonContent
}

export default GroupCard