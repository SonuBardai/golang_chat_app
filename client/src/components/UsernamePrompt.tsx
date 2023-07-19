import React from "react";

const UsernamePrompt: React.FC<{
  username: string;
  setUsername: (username: string) => void;
  handleSetUsername: (username: string) => void;
}> = ({ username, setUsername, handleSetUsername }) => {
  return (
    <div className="fixed inset-0 z-10 flex items-center justify-center bg-black bg-opacity-50">
      <form className="p-8 bg-gray-900 text-white rounded-md shadow">
        <h2 className="text-lg font-bold mb-4">Enter your username</h2>
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          className="w-full bg-gray-800 px-4 py-2 rounded-md border focus:outline-none focus:ring mb-4"
        />
        <button
          onClick={() => handleSetUsername(username)}
          className="px-4 py-2 rounded-md bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring"
        >
          Submit
        </button>
      </form>
    </div>
  );
};

export default UsernamePrompt;
