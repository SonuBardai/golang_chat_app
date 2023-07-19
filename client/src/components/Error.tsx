import React from "react";

const Error: React.FC<{ message: string }> = ({ message }) => {
  const handleReload = () => {
    window.location.reload();
  };

  return (
    <div className="fixed inset-0 z-10 flex items-center justify-center bg-black bg-opacity-50">
      <div className="p-4 bg-gray-900 text-white rounded-md shadow-lg">
        <p className="text-red-600 mb-4">{message}</p>
        <button
          onClick={handleReload}
          className="px-4 py-2 rounded-md bg-blue-500 text-white hover:bg-blue-600 focus:outline-none focus:ring"
        >
          Reload
        </button>
      </div>
    </div>
  );
};

export default Error;
