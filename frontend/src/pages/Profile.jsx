import { useState } from 'react';
import UserInfos from '../components/UserInfos';
import '../styles/profile.css';

const Profile = () => {
  const [message, setMessage] = useState("");


  return (
    <>
      <div>
        {message && <div style={{ color: 'green' }}>{message}</div>}
        <UserInfos setMessage={setMessage} />
      </div>
    </>
  );
}

export default Profile