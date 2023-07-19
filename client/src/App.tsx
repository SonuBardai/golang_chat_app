import { useEffect, useState } from "react";
import ChatInput from "./components/ChatInput";
import ChatWindow from "./components/ChatWindow";
import Error from "./components/Error";
import UsernamePrompt from "./components/UsernamePrompt";
import { RenderableMessage } from "./types/Message";

function App() {
  const [messages, setMessages] = useState<Array<RenderableMessage>>([]);
  const [error, setError] = useState<string | null>();
  const [username, setUsername] = useState("");
  const [showPrompt, setShowPrompt] = useState(true);
  const [socket, setSocket] = useState<WebSocket | null>(null);

  useEffect(() => {
    setSocket(new WebSocket("ws://localhost:3001/ws"));
    return () => {
      socket?.close();
    };
  }, []);

  const handleSendMessage = (message: string) => {
    setMessages([...messages, { user: username, message }]);
  };

  useEffect(() => {
    if (socket) {
      socket.onmessage = (event) => {
        // check message type and handle
      };
      socket.onclose = (event) => {
        setError("Connection lost with server");
      };
    }
  }, [socket]);

  const handleSetUsername = (username: string) => {
    setUsername(username);
    setShowPrompt(false);
    socket?.send(`New user: ${username}`);
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
