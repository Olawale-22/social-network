import { useDispatch } from "react-redux";
import { addPostToUser } from '../redux/userSlice';
import { useState, useContext, useEffect } from "react";
import UserSearch from "./UserSearch";
import WebSocketContext from '../contexts/WebSocketContext';
import ImageUploader from "./ImageUploader";

const PostForm = ({ userId, privacy, setPrivacy }) => {
    const { socket } = useContext(WebSocketContext);
    const url = "http://localhost:8080/upload";
    const [textValue, setTextValue] = useState("");
    const [checkedUsers, setCheckedUsers] = useState([]);
    const [selectedImage, setSelectedImage] = useState(null);

    const handleTextAreaChange = event => {
        setTextValue(event.target.value);
    }

    const handlePrivacyChange = event => {
        setPrivacy(event.target.value);
    }

    const handleSum = async (event, imageName) => {
        try {
            // En premier, télécharger le fichier
            const formData = new FormData(event.target);
            const uploadResponse = await fetch(url, {
                method: "POST",
                body: formData,
            });

            console.log("Upload Response: ", uploadResponse);

            if (uploadResponse.status === 200) {
                const checkedIDs = checkedUsers.map(user => user.ID);
                // Ensuite, envoyer le message WebSocket
                const post = {
                    event: "posts",
                    post: textValue,
                    user_id: userId,
                    privacy: privacy,
                    mentions: JSON.stringify(checkedIDs),
                    image: imageName,
                };

                socket.send(JSON.stringify(post));

                setTextValue("");
                setCheckedUsers([]);
                setSelectedImage(null);
                event.target.reset();
                console.log("SEND DATA");
            } else {
                console.error("SERVER RESPONSE: UPLOAD FAILED");
            }
        } catch (error) {
            console.error("Error:", error);
        }
    };



    const handleSubmit = async () => {
        console.log("SOUMIS SANS IMAGE");
        if ((privacy === "public" || privacy === "private" || privacy === "mentions") && textValue.trim() !== "") {
            const checkedIDs = checkedUsers.map(user => user.ID);
            const post = {
                event: "posts",
                post: textValue,
                user_id: userId,
                privacy: privacy,
                mentions: JSON.stringify(checkedIDs),
                image: "",
            };

            socket.send(JSON.stringify(post));

            setTextValue("");
            setCheckedUsers([]);
            console.log("SEND DATA");
        } else {
            console.log("Veuillez remplir le texte et sélectionner une option de confidentialité valide.");
        }
    }


    return (
        <form
            onSubmit={async (event) => {
                event.preventDefault()
                if (selectedImage) {
                    console.log("TART", selectedImage);
                    handleSum(event, `${url}/${selectedImage.name}`);
                } else {
                    handleSubmit(event);
                }
            }}
            encType="multipart/form-data"
            className="form-post"
        >
            <div className="form-post-body">
                <div className="post-body-content">
                    <label htmlFor="content" className="form-label">Content</label>
                    <textarea
                        id="content"
                        name="content"
                        placeholder="Enter the content..."
                        onChange={handleTextAreaChange}
                        value={textValue}
                    />
                </div>

                <div className='post-body-privacy'>
                    <label className="form-label">Privacy</label>
                    <div className="privacies">
                        <label>
                            <input
                                type="radio"
                                name="privacy"
                                id="public"
                                value="public"
                                checked={privacy === "public"}
                                onChange={handlePrivacyChange}
                            />
                            Public
                        </label>
                        <label>
                            <input
                                type="radio"
                                name="privacy"
                                id="private"
                                value="private"
                                checked={privacy === "private"}
                                onChange={handlePrivacyChange}
                            />
                            Private
                        </label>
                        <label>
                            <input
                                type="radio"
                                name="privacy"
                                id="mentions"
                                value="mentions"
                                checked={privacy === "mentions"}
                                onChange={handlePrivacyChange}
                            />
                            Mentions
                        </label>
                    </div>


                    {privacy === "mentions" && (
                        <UserSearch userId={userId} checkedUsers={checkedUsers} setCheckedUsers={setCheckedUsers} />
                    )}

                    <div className="post-body-image">
                        <label htmlFor="image" className="form-label">Image</label>
                        <ImageUploader onImageSelected={setSelectedImage} />
                    </div>

                    <div className="btn-post">
                        <button type='submit'>Submit</button>
                    </div>
                </div>
            </div>
        </form>
    );
};

export default PostForm;
