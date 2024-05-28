import React, { useEffect, useState, useContext } from 'react';
import { Routes, Route, useNavigate } from 'react-router-dom';
import WebSocketContext from './contexts/WebSocketContext';
import Login from './pages/Login';
import Register from './pages/Register';
import Home from './pages/Home';
import Profile from './pages/Profile';
import Messenger from './pages/Messenger';
import { useDispatch, useSelector } from 'react-redux';
import { setUserInfo, setFollows, setPosts, addFollowToUser, addNotificationToUser, setNotifications, deleteNotification, addPostToUser, setFollowers, addFollowerToUser, deleteFollower, deleteFollow } from './redux/userSlice';
import { setComments, setConversations, setUsers, setGroupChats, addCommentToPost, updateRequestedUser, addMessageToConversation, setGroups, addGroupToUser, updateRequestedGroup, addMemberToGroup, addGroupChatToUser, addEventToGroup, setEvents, setVotes, addVoteToEvent, updatePrivacy } from './redux/dataSlice';
import { checkCookie, setCookie } from './utils/utils';
import axios from 'axios';
import Chat from './components/PrivateChat';
import NavBar from './components/NavBar';
import ChatArea from './components/ChatArea';

function App() {
  const dispatch = useDispatch();
  const notifications = useSelector(state => state.user.notifications);
  const { socket, websocketData } = useContext(WebSocketContext);
  const navigate = useNavigate();
  const [id, setId] = useState(parseInt(localStorage.getItem("user_id")) || null);
  const [sender, setSender] = useState(0);
  const [groupID, setGroupID] = useState(0);
  const [countNotifications, setCountNotifications] = useState(0);
  const userId = localStorage.getItem("user_id");


  useEffect(() => {
    console.log("PASSE");
    const fetchData = async () => {
      try {
        const response = await axios.get(`http://localhost:8080/api/data?studentID=${id}`);
        const { userInfos, follows, followers, posts, users, comments, privateChats, groups, groupChats, events, votes, notifications } = response.data;

        dispatch(setUserInfo(userInfos));
        dispatch(setFollows(follows));
        dispatch(setFollowers(followers));
        dispatch(setPosts(posts));
        dispatch(setUsers(users));
        dispatch(setComments(comments));
        dispatch(setConversations(privateChats));
        dispatch(setGroups(groups));
        dispatch(setGroupChats(groupChats));
        dispatch(setEvents(events));
        dispatch(setVotes(votes))
        dispatch(setNotifications(notifications));

        setCountNotifications(() => notifications !== null ? notifications.length : 0);

        console.log("PROVIDE REDUX STORE");
      } catch (error) {
        if (axios.isCancel(error)) {
          console.log("REQUEST CANCELED");
        } else {
          console.log("ERROR AXIOS: ", error);
        }
      }
    }

    if (checkCookie()) {
      console.log("CONNECTED");
      fetchData();
    }
  }, [id]);


  useEffect(() => {
    if (!checkCookie()) {
      navigate("/login")
      return
    }

    if (websocketData) {
      console.log(websocketData.Event);
      let notificationToDelete;

      switch (websocketData.Event) {
        case 'LogoutSuccess':
          setCookie("", -1);
          localStorage.removeItem("user_id");
          navigate("/login");
          break;

        case "postSuccessful":
          console.log("SERVER RESPONSE: UPLOAD SUCCESS");
          dispatch(addPostToUser({ privacy: websocketData.Data.Privacy, post: websocketData.Data }));
          break;

        case 'CommentSuccessfull':
          dispatch(addCommentToPost(websocketData.Data));
          break;

        case 'follow_user_success_callback':
          console.log(websocketData.Data);
          dispatch(addFollowerToUser(websocketData.Data));
          break;

        case 'follow_user_success':
          console.log("YOU FOLLOWED: ", websocketData.Data);
          dispatch(addFollowToUser(websocketData.Data));
          break;

        case 'follow_user_failed':
          console.log("FAILED FOLLOW: ", websocketData.Data);
          break;

        case 'unfollow_user_success_callback':
          console.log(websocketData.Data);
          dispatch(deleteFollower(websocketData.Data));
          break;

        case 'unfollow_user_success':
          dispatch(deleteFollow(websocketData.Data));
          break;

        case 'unfollowed_user_failed':
          console.log("FAILED UNFOLLOW: ", websocketData.Data);
          break;

        case 'user_follow_request_user':
          console.log("Follow Request received from user successfully from user-ID", websocketData.Data);
          dispatch(addNotificationToUser(websocketData.Data));
          dispatch(updateRequestedUser(websocketData.Data.ItemID))
          console.log(websocketData.Data);
          setSender(websocketData.Data.Sender);
          setCountNotifications(current => current + 1);
          break;

        case 'user_follow_request_user_accepted_callback':
          dispatch(updateRequestedUser(websocketData.Data));
          dispatch(addFollowToUser(websocketData.Data));
          console.log("UPDATED FOLLOW AND REQUESTED");
          console.log(websocketData.Data);
          break;

        case 'user_follow_request_user_accepted':
          console.log(websocketData.Data);
          dispatch(addFollowerToUser(websocketData.Data.FollowerID));
          notificationToDelete = notifications.find(notification => notification.ItemID === websocketData.Data.ItemID);
          console.log(notificationToDelete);
          socket.send(JSON.stringify({
            event: "delete_notification",
            ID: String(notificationToDelete.ID),
            Recipient: userId
          }));
          break;

        case 'user_follow_request_user_declined_callback':
          console.log(websocketData.Data);
          dispatch(updateRequestedUser(websocketData.Data));
          break;

        case 'user_follow_request_user_declined':
          console.log(websocketData.Data);
          notificationToDelete = notifications.find(notification => notification.ItemID === websocketData.Data);
          console.log(notificationToDelete);
          socket.send(JSON.stringify({
            event: "delete_notification",
            ID: String(notificationToDelete.ID),
            Recipient: userId
          }));
          break;

        case 'MessageSuccessfull':
          dispatch(addMessageToConversation(websocketData.Data));
          break;

        case 'group_created':
          dispatch(addGroupToUser(websocketData.Data));
          console.log("GROUP CREATED");
          break;

        case 'admin_group_invitation_user':
          dispatch(addNotificationToUser(websocketData.Data));
          console.log(websocketData.Data);
          setSender(websocketData.Data.Sender);
          setGroupID(websocketData.Data.ItemID);
          setCountNotifications(current => current + 1);
          break;

        case 'admin_group_invitation_user_accepted_callback':
          console.log(websocketData.Data);
          dispatch(addMemberToGroup({ groupID: websocketData.Data.GroupID, member: websocketData.Data.Member }));
          break;

        case 'admin_group_invitation_user_accepted':
          dispatch(addMemberToGroup({ groupID: websocketData.Data.GroupID, member: websocketData.Data.Member }));
          notificationToDelete = notifications.find(notification => notification.ItemID === websocketData.Data.GroupID);
          console.log(notificationToDelete);
          socket.send(JSON.stringify({
            event: "delete_notification",
            ID: String(notificationToDelete.ID),
            Recipient: userId
          }));
          break;

        case 'admin_group_invitation_user_declined_callback':
          console.log(websocketData.Data);
          break;

        case 'admin_group_invitation_user_declined':
          notificationToDelete = notifications.find(notification => notification.ItemID === websocketData.Data);
          console.log(notificationToDelete);
          socket.send(JSON.stringify({
            event: "delete_notification",
            ID: String(notificationToDelete.ID),
            Recipient: userId
          }));
          break;

        case 'user_group_request_admin':
          dispatch(addNotificationToUser(websocketData.Data));
          console.log(websocketData.Data);
          setSender(websocketData.Data.Sender);
          setGroupID(websocketData.Data.ItemID);
          setCountNotifications(current => current + 1);
          break;

        case 'user_group_request_admin_accepted_callback':
          console.log(websocketData.Data);
          dispatch(addMemberToGroup({ groupID: websocketData.Data.GroupID, member: websocketData.Data.Member }));
          break;

        case 'user_group_request_admin_accepted':
          console.log(websocketData.Data);
          dispatch(updateRequestedGroup(websocketData.Data.GroupID));
          dispatch(addMemberToGroup({ groupID: websocketData.Data.GroupID, member: websocketData.Data.Member }));
          notificationToDelete = notifications.find(notification => notification.ItemID === websocketData.Data.GroupID);
          console.log(notificationToDelete);
          socket.send(JSON.stringify({
            event: "delete_notification",
            ID: String(notificationToDelete.ID),
            Recipient: userId
          }));
          break;

        case 'user_group_request_admin_declined':
          console.log(websocketData.Data);
          notificationToDelete = notifications.find(notification => notification.ItemID === websocketData.Data);
          console.log(notificationToDelete);
          socket.send(JSON.stringify({
            event: "delete_notification",
            ID: String(notificationToDelete.ID),
            Recipient: userId
          }));
          break;

        case 'user_group_request_admin_declined_callback':
          console.log(websocketData.Data);
          dispatch(updateRequestedGroup(websocketData.Data));
          break;

        case 'group_chat':
          console.log(websocketData.Data);
          dispatch(addGroupChatToUser(websocketData.Data));
          break;

        case 'event':
          dispatch(addEventToGroup(websocketData.Data));
          break;

        case 'member_event_creation_member':
          console.log(websocketData.Data);
          dispatch(addNotificationToUser(websocketData.Data));
          setSender(websocketData.Data.Sender);
          setGroupID(websocketData.Data.ItemID);
          setCountNotifications(current => current + 1);
          break;

        case 'member_event_creation_member_checked':
          notificationToDelete = notifications.find(notification => notification.ItemID === websocketData.Data);
          console.log(notificationToDelete);
          socket.send(JSON.stringify({
            event: "delete_notification",
            ID: String(notificationToDelete.ID),
            Recipient: userId
          }));
          break;

        case 'vote_successfull':
          console.log(websocketData.Data);
          dispatch(addVoteToEvent(websocketData.Data));
          break;

        case 'delete_notification_success':
          console.log(websocketData);
          dispatch(deleteNotification(websocketData.Data));
          setCountNotifications(current => current - 1);
          break;

        case 'user_update_profil_callback':
          console.log(websocketData.Data);
          dispatch(updatePrivacy({ id: websocketData.Data }));
          break;

        default:
          console.log("ERROR UNKNOW EVENT");
      }
    }
  }, [websocketData])


  const handleAccept = type => {
    console.log("sender: ", sender);
    const data = {
      event: "",
      ID: String(userId),
      Sender: String(sender),
    };

    if (type === "user_follow_request_user") {
      data.event = "user_follow_request_user_accepted"
    } else if (type === "admin_group_invitation_user") {
      console.log(groupID);
      data.event = "admin_group_invitation_user_accepted"
      data.GroupID = String(groupID)
    } else if (type === "user_group_request_admin") {
      console.log("GroupID: ", groupID);
      data.event = "user_group_request_admin_accepted"
      data.GroupID = String(groupID)
    } else if (type === "member_event_creation_member") {
      console.log("GroupID: ", groupID);
      data.event = "member_event_creation_member_checked"
      data.GroupID = String(groupID)
    }
    socket.send(JSON.stringify(data));
    setSender(0);
    setGroupID(0);
  };

  const handleDecline = type => {
    console.log("sender: ", sender);
    const data = {
      event: "",
      ID: String(userId),
      Sender: String(sender),
    };

    if (type === "user_follow_request_user") {
      data.event = "user_follow_request_user_declined"
    } else if (type === "admin_group_invitation_user") {
      data.event = "admin_group_invitation_user_declined"
      data.GroupID = String(groupID)
    } else if (type === "user_group_request_admin") {
      data.event = "user_group_request_admin_declined"
    }

    socket.send(JSON.stringify(data));
    setSender(0);
    setGroupID(0);
  };

  console.log("ID:" , id);


  return (
    <>
      {window.location.pathname !== "/login" && window.location.pathname !== "/register" && <NavBar handleAccept={handleAccept} handleDecline={handleDecline} countNotifications={countNotifications} />}

      <Routes>
        <Route path="/login" element={<Login setId={setId} />} />
        <Route path="/register" element={<Register />} />
        <Route path="/" element={<Home />} />
        <Route path="/profile/:id" element={<Profile />} />
        <Route path="/chat/:id" element={<Chat />} />
        <Route path="/messenger" element={<Messenger />} />
        <Route path="/messenger/:id" element={<ChatArea />} />
      </Routes>
    </>
  )
}

export default App;