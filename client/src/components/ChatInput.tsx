import React, { useState } from "react";

const ChatInput: React.FC<{ onSendMessage: (message: string) => void }> = ({
  onSendMessage,
}) => {
  const [message, setMessage] = useState("");

  const handleSendMessage = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    onSendMessage(message);
    setMessage("");
  };

  return (
    <form onSubmit={handleSendMessage} className="flex p-4">
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        className="flex-grow px-4 py-2 rounded-l-md border border-r-0 focus:outline-none focus:ring bg-gray-800 text-white"
      />
      <button
        type="submit"
        className="px-4 py-2 rounded-r-md bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring"
      >
        Send
      </button>
    </form>
  );
};

export default ChatInput;
