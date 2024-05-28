import '../styles/navbar.css';
import { useContext, useState } from "react";
import { useSelector } from 'react-redux';
import { Link } from "react-router-dom";
import WebSocketContext from "../contexts/WebSocketContext";
import Notification from './Notification';

const NavBar = ({ handleAccept, handleDecline, countNotifications }) => {
	const userInfos = useSelector(state => state.user.userInfo);
	const { socket } = useContext(WebSocketContext);
	const [click, setClick] = useState(false);
	const [openNotifications, setOpenNotifications] = useState(false);
	const userId = localStorage.getItem("user_id");

	const handleClick = event => {
		event.preventDefault();
		const data = {
			event: "logout",
			id: localStorage.getItem("user_id"),
		};
		socket.send(JSON.stringify(data));
	};

	const handleNavClick = () => setClick(!click);
	const closeMenu = () => setClick(false);

	const handleNotifications = () => {
		if (countNotifications > 0) {
			setOpenNotifications(curr => !curr);
		}
	}


	return (
		<>
			{userInfos ? (
				<nav className="navbar">
					<div className='left-part'>
						<Link to='/' className='navbar-logo'>
							{userInfos.Nickname} <i className='underscores'>active</i>
						</Link>

						<div className='menu-icon' onClick={handleNavClick}>
							<i className={click ? 'fas fa-times' : 'fas fa-bars'} />
						</div>
					</div>

					<ul className={click ? 'right-part active' : 'right-part'}>
						<li className='nav-item'>
							<Link to='/' onClick={closeMenu}>
								Home
							</Link>
						</li>
						<li className='nav-item'>
							<Link to={`/profile/${userId}`} onClick={closeMenu}>
								Profile
							</Link>
						</li>
						<li className='nav-item'>
							<Link to="/messenger" onClick={closeMenu}>
								Messenger
							</Link>
						</li>
				
						<li>
							<div className="container-image-notif">
								<img
									src="/notification.png"
									alt="notif"
									className='notif-img'
									onClick={handleNotifications}
								/>

								<div className="notif-count">
									{countNotifications > 0 && <span>{countNotifications}</span>}
								</div>
							</div>


							{openNotifications && countNotifications > 0 && (
								<Notification
									onAccept={handleAccept}
									onDecline={handleDecline}
								/>
							)}
						</li>
						<li>
							<button type="submit" className="btn-logout" onClick={handleClick}>LOGOUT</button>
						</li>
					</ul>
				</nav>
			) : null}
		</>
	);
}

export default NavBar;
