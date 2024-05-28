// userSlice.js
import { createSlice } from "@reduxjs/toolkit";

const userSlice = createSlice({
	name: 'user',
	initialState: {
		userInfo: null,
		follows: [],
		followers: [],
		posts: {
			public: [],
			private: [],
			mentions: []
		},
		notifications: []
	},
	reducers: {
		setUserInfo: (state, action) => {
			state.userInfo = action.payload;
		},
		setFollows: (state, action) => {
			if (action.payload !== null) {
				state.follows = action.payload;
			}
		},
		setFollowers: (state, action) => {
			if (action.payload !== null) {
				state.followers = action.payload;
			}
		},
		setPosts: (state, action) => {
			state.posts.public = action.payload.Public !== null ? action.payload.Public : state.posts.public;
			state.posts.private = action.payload.Private !== null ? action.payload.Private : state.posts.private;
			state.posts.mentions = action.payload.Mentions !== null ? action.payload.Mentions : state.posts.mentions;
		},
		addPostToUser: (state, action) => {
			const { privacy, post } = action.payload;
			switch (privacy) {
				case "private":
					if (state.posts.private.length > 0) {
						state.posts.private.unshift(post);
					} else {
						state.posts.private.push(post);
					}
					break;
				case "mentions":
					if (state.posts.mentions.length > 0) {
						state.posts.mentions.unshift(post);
					} else {
						state.posts.mentions.push(post);
					}
					break;
				default:
					if (state.posts.public.length > 0) {
						state.posts.public.unshift(post);
					} else {
						state.posts.public.push(post);
					}
					break;
			}
		},
		addFollowToUser: (state, action) => {
			state.follows.push(action.payload);
		},
		deleteFollow: (state, action) => {
			const searchedFollowIndex = state.follows.findIndex(item => item === action.payload);
			if (searchedFollowIndex !== -1) {
				state.follows.splice(searchedFollowIndex, 1);
			}
		},
		addFollowerToUser: (state, action) => {
			state.followers.push(action.payload);
		},
		deleteFollower: (state, action) => {
			const searchedFollowIndex = state.followers.findIndex(item => item === action.payload);
			if (searchedFollowIndex !== -1) {
				state.followers.splice(searchedFollowIndex, 1);
			}
		},
		updateOwnPrivacy: (state) => {
			state.userInfo.IsPublic = !state.userInfo.IsPublic;
		},
		setNotifications: (state, action) => {
			state.notifications = action.payload !== null ? action.payload : state.notifications;
		},
		addNotificationToUser: (state, action) => {
			console.log(action.payload);
			if (state.notifications.length > 0) {
				state.notifications.unshift(action.payload);
			} else {
				state.notifications.push(action.payload);
			}
		},
		deleteNotification: (state, action) => {
			state.notifications = state.notifications.filter(notification => notification.ID !== action.payload);
		},
	},
});


export const { setUserInfo, setFollows, setFollowers, setPosts, addPostToUser, addFollowToUser, addFollowerToUser, deleteFollow, deleteFollower, updateOwnPrivacy, setNotifications, addNotificationToUser, deleteNotification } = userSlice.actions;
export default userSlice.reducer;