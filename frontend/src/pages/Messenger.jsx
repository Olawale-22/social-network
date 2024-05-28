import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom'
import { checkCookie } from '../utils/utils';
import  MyGroups  from '../components/MyGroups';

const Messenger = () => {
  const [isLogged, setIsLogged] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    if (!checkCookie()) {
      navigate("/login");
    } else {
      setIsLogged(true);
    }
  })

  if (!isLogged) {
    return null;
  }
  
  return (
    <>
      <div className="container-groups">
        <MyGroups />
      </div>
    </>
  );
}

export default Messenger