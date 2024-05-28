import React, { useContext } from 'react'
import { useSelector } from 'react-redux';
import WebSocketContext from '../contexts/WebSocketContext';
import { formatDate } from '../utils/utils';
import { useParams } from 'react-router-dom';

const EventCard = ({ event, userID }) => {
  const allUsers = useSelector(state => state.data.users);
  const votes = useSelector(state => state.data.votes);
  const { socket } = useContext(WebSocketContext);
  const { id } = useParams();

  const totalEventVotes = votes.filter(vote => vote.EventID === event.ID).length;
  const totalEventGoingVotes = votes.filter(vote => vote.EventID === event.ID && vote.VoteOption === "Going").length;
  const totalEventNotGoingVotes = votes.filter(vote => vote.EventID === event.ID && vote.VoteOption === "Not Going").length;
  const goingPercentage = totalEventVotes > 0 ? Math.round((totalEventGoingVotes / totalEventVotes) * 100) : 0;
  const notGoingPercentage = totalEventVotes > 0 ? Math.round((totalEventNotGoingVotes / totalEventVotes) * 100) : 0;
  const currentEventVotesOfUser = votes.find(vote => vote.EventID === event.ID && vote.UserID === parseInt(userID));


  const getNickname = id => {
    const user = allUsers.find(user => user.ID === id);
    return user.Nickname;
  }


  const handleVote = option => {
    const data = {
      event: "vote",
      groupID: id,
      eventID: String(event.ID),
      userID: userID,
      option
    };

    console.log("DATA VOTE SENT");

    socket.send(JSON.stringify(data));
  }


  return (
    <div className="event">
      <div className="event-header">
        <h3>{event.Title}</h3>
      </div>

      <div className="event-info">
        <span>Creator: {getNickname(event.UserID)}</span>
        <span>Time: {formatDate(event.Time)}</span>
        <span>{event.Description}</span>
      </div>

      <span>{totalEventVotes} people voted</span>
      <div className="event-choices">
        <div className="going-vote">
          <div className="vote-stats">
            <div className="first-part">
              <input
                type="radio"
                name={`vote${event.ID}`}
                id="going" 
                checked={currentEventVotesOfUser ? currentEventVotesOfUser.VoteOption === "Going" : false}
                onChange={() => handleVote('Going')}
              />
              <label htmlFor="going">Going</label>
            </div>

            <div className="second-part">
              <span className="percentage">{goingPercentage} %</span>
            </div>
          </div>

          <div className="progress-bar">
            <div className="bar" style={{ width: `${goingPercentage}%` }}></div>
          </div>
        </div>

        <div className="not-going-vote">
          <div className="vote-stats">
            <div className="first-part">
              <input type="radio"
                name={`vote${event.ID}`}
                id="not-going"
                checked={currentEventVotesOfUser ? currentEventVotesOfUser.VoteOption === "Not Going" : false}
                onChange={() => handleVote('Not Going')}
              />
              <label htmlFor="going">Not Going</label>
            </div>

            <div className="second-part">
              <span className="percentage">{notGoingPercentage} %</span>
            </div>
          </div>

          <div className="progress-bar">
            <div className="bar" style={{ width: `${notGoingPercentage}%` }}></div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default EventCard