import React, { useContext, useEffect, useRef, useState } from 'react'
import '../styles/privateChat.css';
import { useSelector, } from 'react-redux';
import { useParams } from 'react-router-dom';
import WebSocketContext from '../contexts/WebSocketContext';
import EmojiPicker from 'emoji-picker-react';
import { formatDate } from '../utils/utils';

const Chat = () => {
    const conversations = useSelector(state => state.data.conversations);
    const allUsers = useSelector(state => state.data.users);
    const { socket } = useContext(WebSocketContext);
    const { id } = useParams();
    const [message, setMessage] = useState("");
    const scrollableDivRef = useRef(null);
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);
    const filteredMessages = conversations.filter(message => {
        if ((message.Sender === parseInt(localStorage.getItem("user_id")) && parseInt(message.Recipient) === parseInt(id)) ||
            (message.Sender === parseInt(id) && parseInt(message.Recipient) === parseInt(localStorage.getItem("user_id")))) {
            return true
        }
        return false
    })
    const userId = localStorage.getItem("user_id");
    

    useEffect(() => {
        if (scrollableDivRef.current) {
            scrollableDivRef.current.scrollTop = scrollableDivRef.current.scrollHeight;
        }
    }, [conversations]);


    const handleSubmit = event => {
        event.preventDefault();

        const formData = event.target.elements;

        if (formData.message.value !== "") {
            const data = {
                event: "chat",
                ID: localStorage.getItem("user_id"),
                recipient: id,
                content: formData.message.value
            }
            socket.send(JSON.stringify(data));
            setMessage("");
            setShowEmojiPicker(false);
        }
    }

    const renderMessages = messages => {
        return messages.map((message) => {
            if (message && message.Sender !== undefined) {
                const user = allUsers.find((item) => item.ID === message.Sender);
                if (user) {
                    return (
                        <div key={message.ID} className={message.Sender !== parseInt(userId) ? "private-message" : "private-message space"}>
                            <span>{user.Nickname}</span>
                            <span className="private-message-date">{formatDate(message.Time)}</span>
                            <p>{message.Content}</p>
                        </div>
                    )
                }
            }
            return null;
        });
    }

    console.log(filteredMessages);

    return (
        <div className="private-chat-container">

            <div className="chat-container">
                <div className="private-messages" ref={scrollableDivRef}>
                    {conversations && renderMessages(filteredMessages)}
                </div>

                <form className="form-message" onSubmit={handleSubmit}>
                    <div className="form-elements">
                        <div className="message-wrapper">
                            <img src="/emoji.png" alt="emoji icon" onClick={() => setShowEmojiPicker(!showEmojiPicker)} />
                            <input
                                type="text"
                                placeholder="Send a message..."
                                name="message"
                                id="message"
                                value={message}
                                onChange={event => setMessage(event.target.value)}
                            />
                        </div>
                        {showEmojiPicker && (
                            <div className="emojis-wrapper">
                                <EmojiPicker
                                    width="100%"
                                    onEmojiClick={emoji => setMessage(current => current + emoji.emoji)}
                                />
                            </div>
                        )}
                        <button type="submit">Send</button>
                    </div>
                </form>
            </div>
        </div>
    )
}

export default Chat