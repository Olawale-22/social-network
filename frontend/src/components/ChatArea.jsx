import { useContext, useEffect, useRef, useState } from "react";
import { useSelector } from 'react-redux';
import { useParams } from "react-router-dom";
import WebSocketContext from '../contexts/WebSocketContext';
import EventCard from "./EventCard";
import '../styles/chatarea.css';
import { formatDate } from "../utils/utils";

const ChatArea = () => {
  const allUsers = useSelector(state => state.data.users);
  const events = useSelector(state => state.data.events);
  const groupChats = useSelector(state => state.data.gchats);
  const { socket } = useContext(WebSocketContext);
  const { id } = useParams();
  const [message, setMessage] = useState("");
  const [showForm, setShowForm] = useState(false);
  const scrollableDiv = useRef();
  const userId = localStorage.getItem("user_id");
  const filteredGroupChats = groupChats && groupChats.filter(gc => gc.GroupID === parseInt(id));


  useEffect(() => {
    if (scrollableDiv.current) {
      scrollableDiv.current.scrollTop = scrollableDiv.current.scrollHeight;
    }
  }, [groupChats]);


  const handleEvent = event => {
    event.preventDefault();

    console.log(event);
    const formData = event.target.elements;

    const data = {
      event: "event",
      userID: userId,
      groupID: id,
      title: formData.title.value,
      description: formData.description.value,
      time: formData.time.value,
    };
    socket.send(JSON.stringify(data));
    setMessage("");
    setShowForm(curr => !curr);
  }

  const handleMessage = event => {
    event.preventDefault();

    const formData = event.target.elements;

    if (formData.message.value !== "") {
      const data = {
        event: "group_chat",
        SenderID: userId,
        GroupID: id,
        content: formData.message.value,
      };
      socket.send(JSON.stringify(data));
      setMessage("");
    }
  }

  const renderMessages = messages => {
    return messages.map((message) => {
      if (message && message.SenderID !== undefined) {
        const user = allUsers.find((item) => item.ID === message.SenderID);
        if (user) {
          return (
            <div key={message.ID} className={message.SenderID !== parseInt(userId) ? "message member-message" : "message own-message"}>
                <span className="nickname">{user.Nickname}</span>
                <span className="date">{formatDate(message.Time)}</span>
              <p className="content">{message.Content}</p>
            </div>
          )
        }
      }
      return null;
    });
  }


  return (
    <div className="container-events">
      <div className="events-part">
        <div className="btn-create-event">
          <button onClick={() => setShowForm(curr => !curr)}>CREATE EVENT</button>
        </div>

        {showForm && (
          <div className="container-form-events">
            <form className="form-events" onSubmit={handleEvent}>
              <div className="form-title">
                <label htmlFor="title">Title </label>
                <input type="text" name="title" id="title" placeholder="Enter a title..." />
              </div>

              <div className="form-description">
                <label htmlFor="description">Description </label>
                <textarea name="description" id="description" cols="40" rows="7"></textarea>
              </div>

              <div className="form-time">
                <label htmlFor="time">Time </label>
                <input type="datetime-local" name="time" id="time" />
              </div>


              <div className="btn-create-event">
                <button type="submit">Create</button>
              </div>
            </form>
          </div>
        )}

        {events && events.length > 0 ? (
          <div className="events">
            {events.map(event => (
              <EventCard key={event.ID} event={event} userID={userId} />
            ))}
          </div>
        ) : (
          <div className="header-event">
            <h3>There is no incoming event</h3>
          </div>
        )}
      </div>


      <div className="chat-part">
        <div className="container-chat">
          {groupChats && (
            <div className="container-messages" ref={scrollableDiv}>
              {renderMessages(filteredGroupChats)}
            </div>
          )}

          <div className="container-form-chat">
            <form className="form-message" onSubmit={handleMessage}>
              <input
                type="text"
                placeholder="Send a message..."
                name="message"
                id="message"
                value={message}
                onChange={event => setMessage(event.target.value)}
              />
              <button type="submit">Send</button>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ChatArea;

