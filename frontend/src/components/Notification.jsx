import React from 'react';
import "../styles/notification.css";
import { useSelector } from 'react-redux';
import { Link } from 'react-router-dom';


const Notification = ({ onAccept, onDecline }) => {
  const users = useSelector(state => state.data.users);
  const notifications = useSelector(state => state.user.notifications);


  const getText = (nickname, type) => {
    console.log(nickname, type);
    if (type === "follow") {
      return `${nickname} wants to follow you`
    } else if (type === "admin_group_invitation_user") {
      return `${nickname} wants you join his group`
    } else if (type === "user_group_request_admin") {
      return `${nickname} want to join your group`
    } else if (type === "user_follow_request_user") {
      return `${nickname} wants to follow you`
    } else if (type === "member_event_creation_member") {
      console.log("NOTIFICATION EVENT");
      return `${nickname} created a new event`;
    }
  }

  const renderNotifications = currNotifications => {
    return currNotifications.map(notif => {
      const user = users.find(item => item.ID === notif.Sender);
      console.log(user);
      if (user) {
        return (
          <div key={notif.ID} className="notification">
            {notif.Type !== "member_event_creation_member" ? (
              <>
                <img src={user.Avatar} alt="avatar" />
                <p className="notification-text">{getText(user.Nickname, notif.Type)}</p>
                <div className="buttons">
                  <button style={{ backgroundColor: "green" }} onClick={() => onAccept(notif.Type)}>Accept</button>
                  <button style={{ backgroundColor: "red" }} onClick={() => onDecline(notif.Type)}>Decline</button>
                </div>
              </>
            ) : (
              <Link to={`/messenger/${notif.ItemID}`} onClick={() => onAccept(notif.Type)}>
                <img src={user.Avatar} alt="avatar" />
                <span>{getText(user.Nickname, notif.Type)}</span>
              </Link>
            )}
          </div>
        )
      }
      return null;
    });
  };

  console.log(notifications);

  return (
    notifications.length > 0 && (
      <div className="container-notifications">
        {renderNotifications(notifications)}
      </div>
    )
  );
}

export default Notification;