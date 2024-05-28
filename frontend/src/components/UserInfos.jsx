import React from 'react';
import { useSelector } from 'react-redux';
import PrivacyToggle from './PrivacyToggle';
import { useParams } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

const UserInfos = ({ setMessage }) => {
	const { id } = useParams()
	const users = useSelector(state => state.data.users);
	const follows = useSelector(state => state.user.follows);
	const followers = useSelector(state => state.user.followers);
	const userId = localStorage.getItem("user_id");
	const currentData = users.find(user => user.ID === parseInt(id));

	const checkFollowing = id => {
		return follows.some(item => item.User_ID === parseInt(userId) && item.Follower_ID === id)
	}

	const getUser = id => {
		return users.find(user => user.ID === id);
	}

	return (
		<div className="container-profil">
			{currentData && (
				(currentData.IsPublic || checkFollowing(currentData.ID) || id === userId) ? (
					<div className="profile-card">
						<div className="profile-card-header">
							<img src={currentData.Avatar} alt="avatar" />
						</div>

						<div className="profile-card-content">
							<div className="profile-card-identity">
								<label className="label-profile">IDENTITY</label>
								<div className="fullname">
									<FontAwesomeIcon icon="fa-solid fa-address-card" style={{ color: "#003ea8", width: "18px", height: "18px", marginRight: "10px" }} />
									<span>{currentData.FirstName} {currentData.LastName}</span>
								</div>

								<div className="email">
									<FontAwesomeIcon icon="fa-solid fa-envelope" style={{ color: "#003ea8", width: "18px", height: "18px", marginRight: "10px" }} />
									<span>{currentData.Email}</span>
								</div>

								<div className="nickname">
									<FontAwesomeIcon icon="fa-solid fa-user" style={{ color: "#003ea8", width: "18px", height: "18px", marginRight: "10px" }} />
									<span>{currentData.Nickname}</span>
								</div>

								<div className="birthdate">
									<FontAwesomeIcon icon="fa-solid fa-cake-candles" style={{ color: "#003ea8", width: "18px", height: "18px", marginRight: "10px" }} />
									<span>{currentData.BirthDate}</span>
								</div>
							</div>

							<div className="profile-separation">
								<hr />
							</div>

							<div className="profile-card-aboutme">
								<label className="label-profile">ABOUT ME</label>
								<p>{currentData.AboutMe}</p>
							</div>

							<div className="profile-separation">
								<hr />
							</div>

							<label className="label-profile">PRIVACY</label>
							{id === userId ? (
								<div className="profile-privacy-own">
									<PrivacyToggle setMessage={setMessage} />
								</div>
							) : (
								<div className="profile-privacy-other">
									<span>Status: {currentData.IsPublic ? "Public" : "Private"}</span>
								</div>
							)}

							<div className="profile-separation">
								<hr />
							</div>

							<label className="label-profile">FOLLOWS</label>
							{follows && follows.length > 0 && (
								<div className="profile-card-follows">
									{follows.map(follow => {
										const user = getUser(follow)
										return (
											<div className="profile-user" key={user.ID}>
												<img src={user.Avatar} />
												<span>{user.Nickname}</span>
											</div>
										)
									})}
								</div>
							)}

							<div className="profile-separation">
								<hr />
							</div>

							<label className="label-profile">FOLLOWERS</label>
							{followers && followers.length > 0 && (
								<div className="profile-card-followers">
									{followers.map(follower => {
										const user = getUser(follower)
										return (
											<div className="profile-user" key={user.ID}>
												<img src={user.Avatar} />
												<span>{user.Nickname}</span>
											</div>
										)
									})}
								</div>
							)}
						</div>
					</div>
				) : (
					<p>YOU MUST FOLLOW THIS USER BEFORE SEEING HIS PRIVATE PROFILE</p>
				)
			)}
		</div>
	);

}

export default UserInfos;