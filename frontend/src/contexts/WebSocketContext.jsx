import { createContext, useEffect, useRef, useState } from 'react';

const WebSocketContext = createContext();

export function WebSocketProvider({ children }) {
    const [socket, setSocket] = useState(null);
    const connected = useRef(false);
    const [websocketData, setWebsocketData] = useState("");

    useEffect(() => {
        if (!connected.current) {
            const newSocket = new WebSocket(`ws://localhost:8080/ws`);
            console.log("WEBSOCKET CREATED");
            setSocket(newSocket);
            connected.current = true;

            const handleOpen = () => {
                console.log("WebSocket connection opened. Sending data...");
            }

            const handleClose = () => {
                console.log("WebSocket connection closed.");
            }

            const handleMessage = event => {
                if (event.data !== "PING") {
                    console.log(event);
                    const message = JSON.parse(event.data);
                    setWebsocketData(message);
                } else {
                    console.log(event.data);
                }
            }


            newSocket.addEventListener("open", handleOpen);
            newSocket.addEventListener("message", handleMessage);
            newSocket.addEventListener("close", handleClose);


            return () => {
                if (socket) {
                    newSocket.removeEventListener("open", handleOpen);
                    newSocket.removeEventListener("message", handleMessage);
                    newSocket.removeEventListener("close", handleClose);
                    socket.close();
                    console.log("CLOSED WS CONNECTION");
                }
            };
        }
    }, [socket]);
    
    window.onbeforeunload = (e) => {
        socket.send(JSON.stringify({ event: "refresh" }));
        e.preventDefault();
        e.returnValue = '';
    };

    const contextValue = {
        socket,
        websocketData,
    };


    return (
        <WebSocketContext.Provider value={contextValue}>
            {children}
        </WebSocketContext.Provider>
    );
}

export default WebSocketContext;
