package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"social-network/backend/pkg/db/sqlite"
	"social-network/backend/pkg/db/sqlite/queries"
	"social-network/backend/pkg/helper"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ID  int
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	refreshed bool
}

// readPump pumps messages from the websocket connection to the hub.

// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.deconnection <- c
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.handleEvent(message)
	}
}

func (c *Client) handleEvent(message []byte) {
	var data map[string]string
	var dataToSend sqlite.ServToClientMessage
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println("Erreur lors de l'analyse du message JSON:", err)
		return
	}

	event := data["event"]

	fmt.Println(event)

	switch event {
	case "register":
		if check, err := queries.GetByEmailOrNickname(data["email"], data["nickname"]); check && err == nil {
			if err := queries.CreateUser(data["email"], data["password"], data["firstname"], data["lastname"], data["birthdate"], data["avatar"], data["nickname"], data["aboutme"]); err != nil {
				log.Println("Error during creation of user", err.Error())
				return
			}

			dataToSend.Event = "RegisterSuccess"
			dataToSend.Data = "GOT THE CONNECTION"
		} else {
			dataToSend.Event = "RegisterFailure"
			dataToSend.Data = "SIGN-UP CONNECTION IS FUCKED UP"
		}
		sendWsMessage(dataToSend, c)
	case "login":
		if check, err := queries.VerifyLogin(data["username"], data["password"]); err == nil {
			if check {
				info, err := queries.GetInfo(data["username"])
				if err != nil {
					fmt.Println("error getting user info: ", err.Error())
					return
				}

				Uuid := SessionUuid(info.Nickname)
				info.Uuid = Uuid
				dataToSend.Event = "LoginSuccess"
				dataToSend.Data = info

				c.hub.mutex.Lock()
				c.hub.connections[c] = sqlite.Connection{Id: info.ID, IsConnected: true}
				c.hub.mutex.Unlock()
			} else {
				dataToSend.Event = "LoginFailure"
				dataToSend.Data = "Login information is incorrect"
			}
		} else {
			fmt.Println("Error checking credentials: ", err.Error())
			return
		}

		sendWsMessage(dataToSend, c)

	case "posts":
		fmt.Println("\nNEW POST IS GOING TO BE CREATED: ", data["post"])
		if data["privacy"] != "mentions" {
			err := queries.CreatePost(data["user_id"], data["post"], data["privacy"], data["image"])
			if err != nil {
				fmt.Printf("Error Calling Function INSERT POSTS: %v", err.Error())
				return
			}
		} else {
			err := queries.InsertMentionedPost(data["user_id"], data["post"], data["image"], data["privacy"], data["mentions"])
			if err != nil {
				fmt.Printf("Error Calling Function INSERT POSTS: %v", err.Error())
				return
			}
		}

		lastPost, err := queries.GetLastPost()
		if err != nil {
			fmt.Printf("Error Calling Function GET LAST POST: %v", err.Error())
			return
		}

		fmt.Println("POST RETRIEVED: ", lastPost)

		dataToSend.Event = "postSuccessful"
		dataToSend.Data = lastPost

		id, err := helper.MakeInt(data["user_id"])
		if err != nil {
			fmt.Println("Error Calling Function MAKE INT: ", err.Error())
			return
		}

		if data["privacy"] == "public" {
			for client := range c.hub.connections {
				sendWsMessage(dataToSend, client)
			}
		} else if data["privacy"] == "private" {
			followers, err := queries.GetFollowers(id)
			if err != nil {
				fmt.Println("Error Calling Function GET FOLLOWERS: ", err.Error())
				return
			}

			fmt.Println("id: ", id)
			for client, connec := range c.hub.connections {
				if connec.Id == id {
					sendWsMessage(dataToSend, client)
				}
				for _, follower := range followers {
					if follower == connec.Id {
						sendWsMessage(dataToSend, client)
					}
				}
			}
		} else {
			mentionedUsers, ok, err := queries.GetMentionedId(lastPost.PiD, "post_id", "MENTIONS")
			if !ok {
				if err != nil {
					log.Println("Error Calling Function GET MENTIONED ID: ", err.Error())
					return
				}
			}

			for client, connec := range c.hub.connections {
				if connec.Id == id {
					sendWsMessage(dataToSend, client)
				}
				for _, user := range mentionedUsers {
					if user == connec.Id {
						sendWsMessage(dataToSend, client)
					}
				}
			}
		}

	case "comments":
		err := queries.InsertComment(data["ID"], data["postId"], data["comment"], data["image"])
		if err != nil {
			fmt.Printf("Error detected: %v", err.Error())
		}

		lastComment, err := queries.GetLastComment(data["postId"])
		if err != nil {
			fmt.Printf("Error Calling Function GET LAST COMMENT: %v", err.Error())
			return
		}

		fmt.Println("POST RETRIEVED: ", lastComment)

		dataToSend.Event = "CommentSuccessfull"
		dataToSend.Data = lastComment

		for client := range c.hub.connections {
			sendWsMessage(dataToSend, client)
		}

		fmt.Println("COMMENT CREATED ON POST " + data["postId"])

	case "logout":
		id, err := strconv.Atoi(data["id"])
		if err != nil {
			log.Println("error converting id to int logout")
			return
		}

		if err = queries.DeleteSession(id); err != nil {
			log.Println("error deleting session: ", err.Error())
			return
		}

		dataToSend.Event = "LogoutSuccess"
		dataToSend.Data = "Successfully disconnected !"

		c.hub.mutex.Lock()
		c.hub.connections[c] = sqlite.Connection{IsConnected: false}
		c.hub.mutex.Unlock()

		sendWsMessage(dataToSend, c)

	case "follow_user":
		fmt.Printf("FOLLOW: user_id: %s follower_id: %s ", data["ID"], data["FollowId"])
		ok, err := queries.CreateFollower(data["ID"], data["FollowId"])
		if err != nil {
			fmt.Printf("Error detected: %v", err.Error())
		}

		followId, _ := helper.MakeInt(data["FollowId"])
		if ok {
			dataToSend.Event = "follow_user_success"
			dataToSend.Data = followId
		} else {
			dataToSend.Event = "follow_user_exists"
		}

		id, _ := helper.MakeInt(data["ID"])

		for client, connec := range c.hub.connections {
			if connec.Id == id {
				sendWsMessage(dataToSend, client)
			} else if connec.Id == followId {
				sendWsMessage(sqlite.ServToClientMessage{Event: "follow_user_success_callback", Data: id}, client)
			}
		}

		fmt.Println("SENT RESULT TO: ", data["FollowId"])

	case "user_follow_request_user":
		uId, _ := strconv.Atoi(data["ID"])

		itemID := strconv.Itoa(helper.SliceAtoi(data["FollowId"])[0])

		err := queries.CreateNotification(data["ID"], data["FollowId"], "user_follow_request_user", itemID)
		if err != nil {
			fmt.Printf("Error Calling Function CREATE NOTIFICATION: %v", err.Error())
			return
		}

		fmt.Println("NOTIF INSERTED")

		lastNotification, err := queries.GetLastNotification()
		if err != nil {
			fmt.Printf("Error Calling Function GET LAST NOTIFICATION: %v", err.Error())
			return
		}

		err = queries.UpdateRequestedUser(itemID, true)
		if err != nil {
			fmt.Printf("Error Calling Function UPDATE REQUESTED: %v", err.Error())
			return
		}

		dataToSend.Event = "user_follow_request_user"
		dataToSend.Data = lastNotification

		receivers := helper.SliceAtoi(data["FollowId"])

		for client, connec := range c.hub.connections {
			for _, receiver := range receivers {
				if connec.Id == receiver {
					fmt.Println("FOUND TARGET")
					sendWsMessage(dataToSend, client)
				}
			}
		}

		fmt.Printf("REQUEST FROM %d SENT TO %d\n", uId, receivers[0])

	case "unfollow_user":
		ok, err := queries.CheckFollowing(data["ID"], data["FollowId"])
		if err != nil {
			fmt.Printf("Error detected: %v", err.Error())
		}

		id, _ := helper.MakeInt(data["ID"])
		followId, _ := helper.MakeInt(data["FollowId"])
		if ok {
			if err := queries.Unfollow(id, followId); err != nil {
				fmt.Printf("Error detected: %v", err.Error())
			}
			dataToSend.Event = "unfollow_user_success"
			dataToSend.Data = followId
		} else {
			dataToSend.Event = "unfollow_user_failed"
		}

		for client, connec := range c.hub.connections {
			if connec.Id == id {
				sendWsMessage(dataToSend, client)
			} else if connec.Id == followId {
				sendWsMessage(sqlite.ServToClientMessage{Event: "unfollow_user_success_callback", Data: id}, client)
			}
		}

	case "user_follow_request_user_accepted":
		ok, err := queries.CreateFollower(data["Sender"], data["ID"])
		if err != nil {
			fmt.Printf("Error detected: %v", err.Error())
		}

		id, _ := helper.MakeInt(data["ID"])
		if ok {
			dataToSend.Event = "user_follow_request_user_accepted_callback"
			dataToSend.Data = id
		} else {
			dataToSend.Event = "follower_user_exists"
			fmt.Println("FOLLOW EXISTS: ", ok)
		}

		followerId, _ := strconv.Atoi(data["Sender"])

		err = queries.UpdateRequestedUser(data["ID"], false)
		if err != nil {
			fmt.Printf("Error Calling Function UPDATE REQUESTED: %v", err.Error())
			return
		}

		result := struct {
			ItemID     int
			FollowerID int
		}{
			ItemID:     id,
			FollowerID: followerId,
		}

		for client, connec := range c.hub.connections {
			if connec.Id == followerId {
				sendWsMessage(dataToSend, client)
			} else if connec.Id == id {
				sendWsMessage(sqlite.ServToClientMessage{Event: "user_follow_request_user_accepted", Data: result}, client)
			}
		}

	case "user_follow_request_user_declined":
		err = queries.UpdateRequestedUser(data["ID"], false)
		if err != nil {
			fmt.Printf("Error Calling Function UPDATE REQUESTED: %v", err.Error())
			return
		}

		followerId, _ := strconv.Atoi(data["Sender"])
		id, _ := helper.MakeInt(data["ID"])
		dataToSend.Event = "user_follow_request_user_declined_callback"
		dataToSend.Data = id

		for client, connec := range c.hub.connections {
			if connec.Id == followerId {
				sendWsMessage(dataToSend, client)
			} else if connec.Id == id {
				sendWsMessage(sqlite.ServToClientMessage{Event: "user_follow_request_user_declined", Data: id}, client)
			}
		}

	case "new_group":
		fmt.Printf("NEW_GROUP: group_name: %s user_id: %s group_descriptions: %s tagged_users: %s", data["name"], data["user_id"], data["descriptions"], data["mentions"])
		err := queries.CreateGroup(data["name"], data["user_id"], data["descriptions"])
		if err != nil {
			fmt.Printf("Error detected: %v", err.Error())
		}

		dataToSend.Event = "group_created"
		gp, err := queries.GetLastGroup()
		fmt.Println("GP: ", gp)
		if err != nil {
			fmt.Printf("Error GetLastGroup detected: %v", err.Error())
		}
		dataToSend.Data = gp
		for client, c := range c.hub.connections {
			fmt.Println("SENT GROUP TO USERS: " + strconv.Itoa(c.Id))
			sendWsMessage(dataToSend, client)
		}

		mentions := helper.SliceAtoi(data["mentions"])
		if len(mentions) > 0 {
			err := queries.CreateNotification(data["user_id"], data["mentions"], "admin_group_invitation_user", strconv.Itoa(gp.ID))
			if err != nil {
				fmt.Printf("Error Calling Function CREATE NOTIFICATION: %v", err.Error())
				return
			}

			fmt.Println("NOTIF INSERTED")

			lastNotification, err := queries.GetLastNotification()
			if err != nil {
				fmt.Printf("Error Calling Function GET LAST NOTIFICATION: %v", err.Error())
				return
			}

			dataToSend.Event = "admin_group_invitation_user"
			dataToSend.Data = lastNotification

			for client, c := range c.hub.connections {
				for _, mention := range mentions {
					if c.Id == mention {
						sendWsMessage(dataToSend, client)
					}
				}
			}
		}

		fmt.Println("GROUP REQUEST SENT")

	case "admin_group_invitation_user_accepted":
		err := queries.UpdateMentionedUsers(data["GroupID"], data["ID"])
		if err != nil {
			fmt.Printf("Error detected: %v", err.Error())
		}

		gp, err := queries.GetGroupByID(data["GroupID"])
		fmt.Println("GP: ", gp)
		if err != nil {
			fmt.Printf("Error GetLastGroup detected: %v", err.Error())
		}

		groupID, _ := helper.MakeInt(data["GroupID"])
		id, _ := helper.MakeInt(data["ID"])
		sender, _ := helper.MakeInt(data["Sender"])
		dataToSend.Event = "admin_group_invitation_user_accepted_callback"

		result := struct {
			GroupID int
			Member  int
		}{
			GroupID: groupID,
			Member:  id,
		}

		dataToSend.Data = result

		for client, connec := range c.hub.connections {
			if connec.Id == sender {
				sendWsMessage(dataToSend, client)
			} else if connec.Id == id {
				sendWsMessage(sqlite.ServToClientMessage{Event: "admin_group_invitation_user_accepted", Data: result}, client)
			}
		}

	case "admin_group_invitation_user_declined":
		id, _ := helper.MakeInt(data["ID"])
		sender, _ := strconv.Atoi(data["Sender"])

		gp, err := queries.GetGroupByID(data["GroupID"])
		fmt.Println("GP: ", gp)
		if err != nil {
			fmt.Printf("Error GetGroupByID detected: %v", err.Error())
		}

		dataToSend.Event = "admin_group_invitation_user_declined_callback"
		dataToSend.Data = data["ID"] + "refused your group invitation for group n°" + data["GroupID"]

		for client, connec := range c.hub.connections {
			if connec.Id == sender {
				sendWsMessage(dataToSend, client)
			} else if connec.Id == id {
				sendWsMessage(sqlite.ServToClientMessage{Event: "admin_group_invitation_user_declined", Data: gp.ID}, client)
			}
		}

	case "user_group_request_admin":
		err := queries.CreateNotification(data["ID"], data["Recipient"], "user_group_request_admin", data["GroupID"])
		if err != nil {
			fmt.Printf("Error Calling Function CREATE NOTIFICATION: %v", err.Error())
			return
		}

		fmt.Println("NOTIF INSERTED")

		gp, err := queries.GetGroupByID(data["GroupID"])
		fmt.Println("GP: ", gp)
		if err != nil {
			fmt.Printf("Error GetLastGroup detected: %v", err.Error())
		}

		lastNotification, err := queries.GetLastNotification()
		if err != nil {
			fmt.Printf("Error Calling Function GET LAST NOTIFICATION: %v", err.Error())
			return
		}

		err = queries.UpdateRequested(data["GroupID"], true)
		if err != nil {
			fmt.Printf("Error Calling Function UPDATE REQUESTED: %v", err.Error())
			return
		}

		fmt.Println("GROUP ID: ", gp.ID)

		dataToSend.Event = "user_group_request_admin"
		dataToSend.Data = lastNotification
		receivers := helper.SliceAtoi(data["Recipient"])

		for client, connec := range c.hub.connections {
			for _, receiver := range receivers {
				if connec.Id == receiver {
					fmt.Println("FOUND ADMIN")
					sendWsMessage(dataToSend, client)
				}
			}
		}

		fmt.Printf("REQUEST SENT TO ADMIN")

	case "user_group_request_admin_accepted":
		err := queries.UpdateMentionedUsers(data["GroupID"], data["Sender"])
		if err != nil {
			fmt.Printf("Error detected: %v", err.Error())
		}

		gp, err := queries.GetGroupByID(data["GroupID"])
		fmt.Println("GP: ", gp)
		if err != nil {
			fmt.Printf("Error GetGroupByID detected: %v", err.Error())
		}

		err = queries.UpdateRequested(data["GroupID"], false)
		if err != nil {
			fmt.Printf("Error Calling Function UPDATE REQUESTED: %v", err.Error())
			return
		}

		groupID, _ := helper.MakeInt(data["GroupID"])
		id, _ := helper.MakeInt(data["ID"])
		sender, _ := strconv.Atoi(data["Sender"])

		result := struct {
			GroupID int
			Member  int
		}{
			GroupID: groupID,
			Member:  sender,
		}

		dataToSend.Event = "user_group_request_admin_accepted_callback"
		dataToSend.Data = result

		for client, connec := range c.hub.connections {
			if connec.Id == sender {
				sendWsMessage(dataToSend, client)
			} else if connec.Id == id {
				sendWsMessage(sqlite.ServToClientMessage{Event: "user_group_request_admin_accepted", Data: result}, client)
			}
		}

	case "user_group_request_admin_declined":
		groupID, _ := helper.MakeInt(data["GroupID"])
		id, _ := helper.MakeInt(data["ID"])
		sender, _ := strconv.Atoi(data["Sender"])
		dataToSend.Event = "user_group_request_admin_declined_callback"
		dataToSend.Data = groupID

		err = queries.UpdateRequested(data["GroupID"], false)
		if err != nil {
			fmt.Printf("Error Calling Function UPDATE REQUESTED: %v", err.Error())
			return
		}

		for client, connec := range c.hub.connections {
			if connec.Id == sender {
				sendWsMessage(dataToSend, client)
			} else if connec.Id == id {
				sendWsMessage(sqlite.ServToClientMessage{Event: "user_group_request_admin_declined", Data: groupID}, client)
			}
		}

	case "delete_notification":
		id, _ := helper.MakeInt(data["ID"])
		notifID, err := queries.DeleteCheckedLastNofification(id)
		if err != nil {
			fmt.Printf("Error detected: %v", err.Error())
		}

		dataToSend.Event = "delete_notification_success"
		dataToSend.Data = notifID

		recipient, _ := helper.MakeInt(data["Recipient"])

		for client, connec := range c.hub.connections {
			if connec.Id == recipient {
				sendWsMessage(dataToSend, client)
			}
		}

	case "profile":
		id, err := strconv.Atoi(data["id"])
		if err != nil {
			log.Println("error converting id to int profile: ", err.Error())
			return
		}

		if currentState, err := queries.GetPrivacyByID(id); err == nil {
			newState := !currentState

			if err = queries.UpdatePrivacy(id, newState); err != nil {
				log.Println("error updating pricacy: ", err.Error())
				return
			}
			dataToSend.Event = "user_update_profil"
		} else {
			log.Println("error getting user info: ", err.Error())
			return
		}

		for client, connec := range c.hub.connections {
			if connec.Id == id {
				sendWsMessage(dataToSend, client)
			} else {
				sendWsMessage(sqlite.ServToClientMessage{Event: "user_update_profil_callback", Data: id}, client)
			}
		}

	case "group_chat":
		err := queries.InsertGroupChat(data["SenderID"], data["GroupID"], data["content"])
		if err != nil {
			fmt.Printf("Error Calling Function INSERT GROUP CHAT: %v", err.Error())
			return
		}

		lastChat, err := queries.GetLastGroupChat()
		if err != nil {
			fmt.Printf("Error Calling Function GET LAST CHAT: %v", err.Error())
			return
		}

		dataToSend.Event = "group_chat"
		dataToSend.Data = lastChat

		admin, err := queries.GetAdmin(data["GroupID"])
		if err != nil {
			fmt.Printf("Error Calling Function Get admin : %v", err.Error())
			return
		}

		members, err := queries.GetGroupChatMembers(lastChat.GroupID)
		if err != nil {
			fmt.Printf("Error Calling Function GET : %v", err.Error())
			return
		}

		membersConv := helper.SliceAtoi(members)
		membersConv = append(membersConv, admin)

		for client, c := range c.hub.connections {
			for _, member := range membersConv {
				if member == c.Id {
					fmt.Println("MATCHED MEMBER ID")
					sendWsMessage(dataToSend, client)
				}
			}
		}

	case "event":
		err := queries.CreateEvent(data["userID"], data["groupID"], data["title"], data["description"], data["time"])
		if err != nil {
			fmt.Printf("Error Calling Function CREATE EVENT: %v", err.Error())
			return
		}

		lastEvent, err := queries.GetLastEvent(data["groupID"])
		if err != nil {
			fmt.Printf("Error Calling Function GET LAST EVENT: %v", err.Error())
			return
		}

		userID, _ := helper.MakeInt(data["userID"])
		groupID, _ := helper.MakeInt(data["groupID"])
		admin, err := queries.GetAdmin(data["groupID"])
		if err != nil {
			fmt.Printf("Error Calling Function GET ADMIN: %v", err.Error())
			return
		}

		members, err := queries.GetGroupChatMembers(groupID)
		if err != nil {
			fmt.Printf("Error Calling Function GET GROUP CHAT MEMBERS : %v", err.Error())
			return
		}

		membersConv := helper.SliceAtoi(members)
		membersConv = append(membersConv, admin)

		dataToSend.Event = "event"
		dataToSend.Data = lastEvent

		err = queries.CreateNotification(data["userID"], members, "member_event_creation_member", data["groupID"])
		if err != nil {
			fmt.Printf("Error Calling Function CREATE NOTIFICATION: %v", err.Error())
			return
		}

		lastNotification, err := queries.GetLastNotification()
		if err != nil {
			fmt.Printf("Error Calling Function GET LAST NOTIFICATION : %v", err.Error())
			return
		}

		for client, c := range c.hub.connections {
			if c.Id == userID {
				fmt.Println("MATCHED CREATOR ID")
				sendWsMessage(dataToSend, client)
			}
			for _, member := range membersConv {
				if c.Id == member && member != userID {
					fmt.Println("MATCHED MEMBER ID")
					sendWsMessage(sqlite.ServToClientMessage{Event: "member_event_creation_member", Data: lastNotification}, client)
					time.Sleep(2 * time.Second)
					sendWsMessage(dataToSend, client)
				}
			}
		}

	case "member_event_creation_member_checked":
		groupID, _ := helper.MakeInt(data["GroupID"])
		id, _ := helper.MakeInt(data["ID"])

		dataToSend.Event = "member_event_creation_member_checked"
		dataToSend.Data = groupID

		for client, connec := range c.hub.connections {
			if connec.Id == id {
				sendWsMessage(dataToSend, client)
			}
		}

	case "vote":
		eventID := data["eventID"]
		userID := data["userID"]
		option := data["option"]
		fmt.Println("GROUP ID: ", data["groupID"])
		hasVoted := false

		count, err := queries.GetCountVotes(eventID, userID)
		if err != nil {
			fmt.Printf("Error Calling Function GET COUNT VOTES: %v", err.Error())
			return
		}

		if count > 0 {
			err = queries.UpdateVote(eventID, userID, option)
			if err != nil {
				fmt.Printf("Error Calling Function UPDATE VOTE: %v", err.Error())
				return
			}
			hasVoted = true
		} else {
			err = queries.CreateVote(eventID, userID, option)
			if err != nil {
				fmt.Printf("Error Calling Function CREATE VOTE: %v", err.Error())
				return
			}
		}

		admin, err := queries.GetAdmin(data["groupID"])
		if err != nil {
			fmt.Printf("Error Calling Function GET ADMIN: %v", err.Error())
			return
		}

		lastVote, err := queries.GetLastVote(eventID, userID)
		if err != nil {
			fmt.Printf("Error Calling Function GET LAST VOTE: %v", err.Error())
			return
		}

		groupID, _ := helper.MakeInt(data["groupID"])

		members, err := queries.GetGroupChatMembers(groupID)
		if err != nil {
			fmt.Printf("Error Calling Function GET GROUP CHAT MEMBERS : %v", err.Error())
			return
		}

		fmt.Println("MEMBERS: ", members)

		membersConv := helper.SliceAtoi(members)
		membersConv = append(membersConv, admin)

		fmt.Println("MEMBERS WITH ADMIN: ", membersConv)

		result := struct {
			Vote            queries.Vote
			HasAlreadyVoted bool
		}{
			Vote:            *lastVote,
			HasAlreadyVoted: hasVoted,
		}

		dataToSend.Event = "vote_successfull"
		dataToSend.Data = result

		for client, c := range c.hub.connections {
			for _, member := range membersConv {
				if member == c.Id {
					fmt.Println("MATCHED MEMBER ID: ", member)
					sendWsMessage(dataToSend, client)
				}
			}
		}

	case "chat":
		err := queries.CreatePrivateChat(data["ID"], data["recipient"], data["content"])
		if err != nil {
			fmt.Printf("Error detected: %v", err.Error())
		}
		lastPrivateChat, err := queries.GetLastPrivateChat()
		if err != nil {
			fmt.Printf("Error Calling Function GET LAST MESSAGE: %v", err.Error())
			return
		}
		fmt.Println("MESSAGE RETRIEVED: ", lastPrivateChat)
		dataToSend.Event = "MessageSuccessfull"
		dataToSend.Data = lastPrivateChat

		uId, _ := strconv.Atoi(data["ID"])
		fId, _ := strconv.Atoi(data["recipient"])
		for client, connec := range c.hub.connections {
			if connec.Id == uId || connec.Id == fId {
				fmt.Println("FOUND YOU")
				sendWsMessage(dataToSend, client)
			}
		}
		fmt.Println("MESSAGE CREATED ON PRIVATECHAT : ", lastPrivateChat.ID)

	default:
		c.refreshed = true
	}
	fmt.Println("connections: ", c.hub.connections)
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("NEW CLIENT: ", hub.connections)

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), refreshed: false}

	client.hub.connection <- client

	go client.readPump()
	go client.writePump()
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Write:", err.Error())
			return
		}
	}
}

func sendWsMessage(dataToSend sqlite.ServToClientMessage, client *Client) {
	message, err := json.Marshal(dataToSend)
	if err != nil {
		log.Fatal("ERROR WRITING MESSAGE: ", err.Error())
	}

	if client.conn != nil {
		client.send <- message
	}
}

func SessionUuid(name string) string {
	uuidGenerated, _ := uuid.NewV4()
	uuid := uuidGenerated.String()
	if err := queries.AddSession(name, uuid); err != nil {
		log.Println("Eroro creating session", err.Error())
	}
	return uuid
}
