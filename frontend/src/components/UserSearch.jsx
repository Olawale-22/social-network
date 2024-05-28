import { useEffect, useState } from "react";
import { useDispatch, useSelector } from 'react-redux';

const UserSearch = ({ userId, checkedUsers, setCheckedUsers }) => {
    const users = useSelector(state => state.data.users);
    const [researchedValue, setResearchedValue] = useState("");
    const [isFocused, setIsFocused] = useState(false);
    const [suggestions, setSuggestions] = useState([]);
    

    useEffect(() => {
        if (users) {
            const newUsers = users.filter(user => user.ID !== parseInt(userId));
            setSuggestions(newUsers);
        }
    }, [users]);



    const handleResearch = event => {
        const inputValue = event.target.value;
        setResearchedValue(inputValue);
        const mentionedUsers = users.filter(user => user.Nickname.toLowerCase().includes(inputValue.toLowerCase()) && user.ID !== parseInt(userId));
        setSuggestions(mentionedUsers);
    };

    const handleToggle = user => {
        if (!checkedUsers.some(u => u.ID === user.ID)) {
            setCheckedUsers(current => [...current, user]);
        }
    };

    const handleInputFocus = () => {
        setIsFocused(true);
    };

    const handleInputBlur = () => {
        setIsFocused(false);
    };

    const removeSuggestion = id => {
        const copy = [...checkedUsers];
        const newCheckedValues = copy.filter(u => u.ID !== id);
        setCheckedUsers(newCheckedValues);
    }

    console.log(checkedUsers);

    return (
        <div className="container-mentions">
            {checkedUsers.length > 0 && (
                <ul className="results">
                    {checkedUsers.map(user => {
                        return (
                            < li key={user.ID} >
                                <img className="avatar-mentioned" src={user.Avatar} alt="avatar" />
                                <span>{user.Nickname}</span>
                                <span onClick={() => removeSuggestion(user.ID)}>x</span>
                            </li>
                        )
                    })}
                </ul>
            )}

            <input className="search-mentions"
                name="researchedValue"
                placeholder="Search a user..."
                onChange={handleResearch}
                value={researchedValue}
                onFocus={handleInputFocus}
                onBlur={handleInputBlur}
            />

            {
                isFocused && (
                    <ul className="users-mentioned">
                        {suggestions.length > 0 ? (
                            suggestions.map(user => (
                                <li key={user.ID} onMouseDown={() => handleToggle(user)}>
                                    <img className="avatar-mentioned" src={user.Avatar} alt="avatar" />
                                    <span>{user.Nickname}</span>
                                </li>
                            ))
                        ) : (
                            <li>No user found</li>
                        )}
                    </ul>
                )
            }
        </div >
    );
};

export default UserSearch;