import { useEffect, useState } from "react";
import ChatInput from "./components/ChatInput";
import ChatWindow from "./components/ChatWindow";
import Error from "./components/Error";
import UsernamePrompt from "./components/UsernamePrompt";
import {
  BroadcastMessage,
  NewUserMessage,
  RenderableMessage,
  ServerMessage,
} from "./types/Message";

function App() {
  const [messages, setMessages] = useState<Array<RenderableMessage>>([]);
  const [error, setError] = useState<string | null>();
  const [username, setUsername] = useState("");
  const [showPrompt, setShowPrompt] = useState(true);
  const [socket, setSocket] = useState<WebSocket | null>(null);

  useEffect(() => {
    if (!socket) {
      setSocket(new WebSocket("ws://localhost:3001/ws"));
    }
    return () => {
      socket?.close();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleSendMessage = (message: string) => {
    const newMessage: BroadcastMessage = {
      messageType: "broadcastMessage",
      username,
      content: message,
    };
    socket?.send(JSON.stringify(newMessage));
  };

  useEffect(() => {
    if (socket) {
      socket.onmessage = (event) => {
        const message: ServerMessage = JSON.parse(event.data);
        switch (message.messageType) {
          case "newUser":
            setMessages((currentMessages) => [
              ...currentMessages,
              {
                message: `New user ${message.username} has joined the chat`,
                user: "server",
                level: "info",
              },
            ]);
            break;
          case "broadcastMessage":
            if ("content" in message) {
              setMessages((currentMessages) => [
                ...currentMessages,
                {
                  message: message.content,
                  user: message.username,
                  level: null,
                },
              ]);
            }
            break;
          default:
            console.error("Unknown message received");
            break;
        }
      };
      socket.onclose = (error) => {
        console.error(error);
        setError("Connection lost with server");
      };
    }
  }, [socket]);

  const handleSetUsername = (username: string) => {
    setUsername(username);
    setShowPrompt(false);
    const newUser: NewUserMessage = { messageType: "newUser", username };
    socket?.send(JSON.stringify(newUser));
  };

  return (
    <div className="h-full my-auto mx-auto bg-gray-900 text-white shadow rounded-md">
      {error && <Error message={error} />}
      {showPrompt && (
        <UsernamePrompt
          username={username}
          setUsername={setUsername}
          handleSetUsername={handleSetUsername}
        />
      )}
      <ChatWindow messages={messages} username={username} />
      <ChatInput onSendMessage={handleSendMessage} />
    </div>
  );
}

export default App;
