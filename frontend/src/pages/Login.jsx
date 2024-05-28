import { useState, useEffect, useContext } from 'react'
import { setCookie } from '../utils/utils'
import { NavLink, useNavigate } from 'react-router-dom'
import '../styles/login.css';
import WebSocketContext from '../contexts/WebSocketContext'

const Login = ({ setId }) => {
    const { socket, websocketData } = useContext(WebSocketContext);
    const navigate = useNavigate();
    const [error, setError] = useState("");


    useEffect(() => {
        if (websocketData) {
            if (websocketData.Event === 'LoginSuccess') {
                const userCode = websocketData.Data.Uuid;
                localStorage.setItem("user_id", JSON.stringify(websocketData.Data.ID));
                setId(websocketData.Data.ID);
                setCookie(userCode, 1);
                console.log("COOKIE SET");
                setError("")
                navigate("/")
            } else {
                navigate("/login")
                setError(websocketData.Data)
                setTimeout(() => {
                    setError("")
                }, 2000)
            }
        }
    }, [websocketData]);


    const handleSubmit = (event) => {
        event.preventDefault();
        const formData = event.target.elements;
        const data = {
            event: "login",
            username: formData.username.value,
            password: formData.password.value
        }
        socket.send(JSON.stringify(data));
    };

    return (
        <div className="container-login">
            <form onSubmit={handleSubmit} className="form-login">
                <div className="form-login-header">
                    <h2>LOGIN</h2>
                </div>

                <div className="form-login-body">
                    <div className="form-login-username">
                        <label className="form-login-label" htmlFor="username">
                            username / email
                        </label>
                        <input type="text" name="username" id="username" autoComplete="username" />
                    </div>

                    <div className="form-login-password">
                        <label className="form-login-label" htmlFor="password">
                            Password
                        </label>
                        <input type="password" name="password" id="password" required autoComplete="current-password" />
                    </div>

                    <div className="form-login-btn">
                        <button type="submit">LOGIN</button>
                    </div>

                    <hr />

                    <div className="not-registered">
                        <button onClick={() => navigate("/register")}>
                            CREATE NEW ACCOUNT
                        </button>
                    </div>
                </div>
                {error && <p style={{ color: "red" }}>Error: {error}</p>}
            </form>
        </div>
    )
}

export default Login