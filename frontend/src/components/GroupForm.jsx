import React, { useContext } from 'react'
import WebSocketContext from '../contexts/WebSocketContext';
import UserSearch from './UserSearch';

const GroupForm = ({ checkedUsers, setCheckedUsers }) => {
    const { socket } = useContext(WebSocketContext);
    const userId = localStorage.getItem("user_id");


    const handleSubmit = event => {
        event.preventDefault();
        const formData = new FormData(event.target);
        const group_name = formData.get("group_name");
        const descriptions = formData.get("description");

        const checkedIDs = checkedUsers.map(user => user.ID);

        const groupData = {
            event: "new_group",
            name: group_name,
            user_id: userId,
            descriptions: descriptions,
            mentions: JSON.stringify(checkedIDs),
        };

        console.log("Username:", group_name);
        console.log("Email:", descriptions)
        console.log("Tagged:", checkedUsers)
        socket.send(JSON.stringify(groupData));

        setCheckedUsers([]);
        event.target.reset();
    };


    return (
        <>
            <form onSubmit={handleSubmit}>
                <label>Title</label>
                <input type="text" name="group_name" />

                <label>Description</label>
                <input type="text" name="description" />

                <label>Invite users</label>
                <UserSearch userId={userId} checkedUsers={checkedUsers} setCheckedUsers={setCheckedUsers} />

                <div className="create-group-wrapper">
                    <button type="submit">CREATE</button>
                </div>
            </form>
        </>
    )
}

export default GroupForm