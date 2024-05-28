import React, { useState } from 'react';
import '../styles/contact.css';
import { useSelector } from 'react-redux';

import UserCard from '../components/UserCard';

const Users = ({ setReloadComponent, buttonText, setButtonText }) => {
  const users = useSelector(state => state.data.users);
  const userId = localStorage.getItem("user_id");
  const filteredUsers = users.filter(user => user.ID !== parseInt(userId));


  return (
    <>
      {users && (
        <>
          <h3 className="title">USERS</h3>
          {filteredUsers.map(item => (
            <UserCard key={item.ID} user={item} userId={userId} setReloadComponent={setReloadComponent} buttonText={buttonText} setButtonText={setButtonText} />
          ))}
        </>
      )}
    </>
  );
};

export default Users;