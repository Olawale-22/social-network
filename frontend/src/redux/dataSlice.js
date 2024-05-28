// userSlice.js
import { createSlice } from "@reduxjs/toolkit";

const dataSlice = createSlice({
    name: 'data',
    initialState: {
        users: [],
        comments: [],
        conversations: [],
        groups: [],
        gchats: [],
        events: [],
        votes: []
    },
    reducers: {
        setUsers: (state, action) => {
            state.users = action.payload !== null ? action.payload : state.users;
        },
        updatePrivacy: (state, action) => {
            const { id } = action.payload;
            const userIndex = state.users.findIndex(user => user.ID === id);
            if (userIndex !== -1) {
                const updatedUsers = [...state.users];
                updatedUsers[userIndex].IsPublic = !updatedUsers[userIndex].IsPublic;
                state.users = updatedUsers;
            }
        },
        updateRequestedUser: (state, action) => {
            const userIndex = state.users.findIndex(user => user.ID === parseInt(action.payload));
            if (userIndex !== -1) {
                const updatedUsers = [...state.users];
                updatedUsers[userIndex].IsRequested = !updatedUsers[userIndex].IsRequested;
                state.users = updatedUsers;
            }
        },
        setComments: (state, action) => {
            state.comments = action.payload !== null ? action.payload : state.comments;
        },
        addCommentToPost: (state, action) => {
            if (state.comments.length > 0) {
                state.comments.unshift(action.payload);
            } else {
                state.comments.push(action.payload);
            }
        },
        setConversations: (state, action) => {
            state.conversations = action.payload !== null ? action.payload : state.conversations;
        },
        addMessageToConversation: (state, action) => {
            state.conversations.push(action.payload);
        },
        setGroups: (state, action) => {
            state.groups = action.payload !== null ? action.payload : state.groups;
        },
        addGroupToUser: (state, action) => {
            if (action.payload !== null) {
                state.groups.unshift(action.payload)
            }
        },
        addMemberToGroup: (state, action) => {
            const { groupID, member } = action.payload;
            const groupIndex = state.groups.findIndex(group => group.ID === groupID);

            if (groupIndex !== -1) {
                const updatedGroups = [...state.groups];
                const updatedGroup = { ...updatedGroups[groupIndex] };
                updatedGroup.Mentions = [...updatedGroup.Mentions, member];
                updatedGroups[groupIndex] = updatedGroup;

                state.groups = updatedGroups;
            }
        },
        updateRequestedGroup: (state, action) => {
            const groupIndex = state.groups.findIndex(group => group.ID === parseInt(action.payload));
            if (groupIndex !== -1) {
                const updatedGroups = [...state.groups];
                updatedGroups[groupIndex].IsRequested = !updatedGroups[groupIndex].IsRequested;
                state.groups = updatedGroups;
            }
        },
        setGroupChats: (state, action) => {
            if (action.payload !== null) {
                state.gchats = action.payload;
            }
        },
        addGroupChatToUser: (state, action) => {
            state.gchats.push(action.payload);
        },
        setEvents: (state, action) => {
            if (action.payload !== null) {
                state.events = action.payload;
            }
        },
        addEventToGroup: (state, action) => {
            if (state.events.length > 0) {
                state.events.unshift(action.payload);
            } else {
                state.events.push(action.payload);
            }
        },
        setVotes: (state, action) => {
            state.votes = action.payload !== null ? action.payload : state.votes;
        },
        addVoteToEvent: (state, action) => {
            const { HasAlreadyVoted, Vote } = action.payload;
            if (HasAlreadyVoted) {
                const voteIndex = state.votes.findIndex(vote => vote.EventID === Vote.EventID && vote.UserID === Vote.UserID);
                if (voteIndex !== -1) {
                    const updatedVotes = [...state.votes];
                    updatedVotes[voteIndex].VoteOption = Vote.VoteOption;
                    state.votes = updatedVotes;
                }
            } else {
                state.votes.push(Vote);
            }
        },
    },
});

export const getUserById = (users, id) => {
    return users.find(user => user.ID === id);
}


export const { setUsers, updatePrivacy, updateRequestedUser, setComments, addCommentToPost, setConversations, addMessageToConversation, setGroups, addGroupToUser, addMemberToGroup, updateRequestedGroup, setGroupChats, addGroupChatToUser, setEvents, addEventToGroup, setVotes, addVoteToEvent } = dataSlice.actions;
export default dataSlice.reducer;