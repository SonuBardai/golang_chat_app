import { RenderableMessage } from "../types/Message";

const ChatWindow: React.FC<{
  messages: RenderableMessage[];
  username: string;
}> = ({ messages, username }) => {
  return (
    <ul className="p-4 min-h-[90vh] overflow-y-auto bg-gray-800 rounded-md">
      {messages.map((message, index) => (
        <li
          key={index}
          className={
            message.user === username
              ? "mb-2 text-right"
              : message.user === "server"
              ? "mb-2 text-center"
              : "mb-2 text-left"
          }
        >
          {message.user && (
            <div
              className={
                message.user === username
                  ? "font-bold text-blue-500"
                  : message.user === "server"
                  ? "hidden"
                  : "font-bold text-green-500"
              }
            >
              {message.user}:
            </div>
          )}
          <div
            className={
              message.user === username
                ? "inline-block bg-blue-800 rounded px-2 py-1 text-white"
                : message.user === "server"
                ? "inline-block px-2 py-1 text-white text-sm"
                : "inline-block bg-green-800 rounded px-2 py-1 text-white"
            }
          >
            {message.message}
          </div>
        </li>
      ))}
    </ul>
  );
};

export default ChatWindow;
