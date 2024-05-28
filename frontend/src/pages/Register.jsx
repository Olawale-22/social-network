import React, { useState, useEffect, useContext } from 'react';
import '../styles/register.css';
import { useNavigate } from 'react-router-dom';
import WebSocketContext from '../contexts/WebSocketContext';
import ImageUploader from '../components/ImageUploader';
import { formatDate } from '../utils/utils';

const Register = () => {
  const { socket, websocketData } = useContext(WebSocketContext);
  const navigate = useNavigate();
  const [error, setError] = useState('');
  const [selectedImage, setSelectedImage] = useState(null);

  useEffect(() => {
    if (websocketData) {
      if (websocketData.Event === 'RegisterSuccess') {
        navigate('/login');
      } else {
        setError(websocketData.Data);
        setTimeout(() => {
          setError('');
        }, 2000);
      }
    }
  }, [websocketData]);

  
  const handleSubmit = async (event) => {
    event.preventDefault();

    if (selectedImage) {
      try {
        const formData = new FormData();
        formData.append('image', selectedImage);

        for (const [k, v] of formData) {
          console.log(k, v);
        }

        const uploadResponse = await fetch('http://localhost:8080/upload', {
          method: 'POST',
          body: formData,
        });

        if (uploadResponse.status === 200) {
          const formElements = event.target.elements;

          const data = {
            event: "register",
            email: formElements.email.value,
            password: formElements.password.value,
            firstname: formElements.firstname.value,
            lastname: formElements.lastname.value,
            birthdate: formatDate(formElements.birthdate.value),
            avatar: `http://localhost:8080/upload/${selectedImage.name}`,
            nickname: formElements.nickname.value,
            aboutme: formElements.aboutme.value,
          };

          socket.send(JSON.stringify(data));
          setSelectedImage("");

        } else {
          console.error('SERVER RESPONSE: UPLOAD FAILED');
        }
      } catch (error) {
        console.error('Error:', error);
      }
    }

    // Le reste de votre logique de soumission de formulaire
  };

  return (
    <div className="container-register">
      <form onSubmit={handleSubmit} className="form-register">
        <div className="form-register-header">
          <h2>REGISTER</h2>
        </div>

        <div className="form-register-body">
          <div className="form-register-email">
            <label htmlFor="email" className="form-register-label">
              Email
            </label>
            <input type="email" name="email" id="email" required />
          </div>

          <div className="form-register-password">
            <label htmlFor="password" className="form-register-label">
              Password
            </label>
            <input type="password" name="password" id="password" required />
          </div>

          <div className="form-register-firstname">
            <label htmlFor="firstname" className="form-register-label">
              Firstname
            </label>
            <input type="text" name="firstname" id="firstname" required />
          </div>

          <div className="form-register-lastname">
            <label htmlFor="lastname" className="form-register-label">
              Lastname
            </label>
            <input type="text" name="lastname" id="lastname" required />
          </div>

          <div className="form-register-birthdate">
            <label htmlFor="birthdate" className="form-register-label">
              Birthdate
            </label>
            <input type="date" name="birthdate" id="birthdate" required />
          </div>

          <div className="form-register-avatar">
            <label htmlFor="image" className="form-register-label">
              Avatar
            </label>
            <ImageUploader onImageSelected={setSelectedImage} />
          </div>

          <div className="form-register-nickname">
            <label htmlFor="nickname" className="form-register-label">
              Nickname
            </label>
            <input type="text" name="nickname" id="nickname" />
          </div>

          <div className="form-register-aboutme">
            <label htmlFor="aboutme" className="form-register-label">
              Aboutme
            </label>
            <textarea name="aboutme" id="aboutme" />
          </div>

          <div className="form-register-btn">
            <button type="submit">REGISTER</button>
          </div>
        </div>

        {error && <p style={{ color: 'red' }}>{error}</p>}
      </form>
    </div>
  );
}

export default Register